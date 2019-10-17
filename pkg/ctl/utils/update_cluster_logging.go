package utils

import (
	"fmt"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/spf13/pflag"

	"k8s.io/apimachinery/pkg/util/sets"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func enableLoggingCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("update-cluster-logging", "Update cluster logging configuration", "")

	var typesEnabled []string
	var typesDisabled []string
	cmd.SetRunFuncWithNameArg(func() error {
		return doEnableLogging(cmd, typesEnabled, typesDisabled)
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlagWithDeprecated(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cmd)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmd.FlagSetGroup.InFlagSet("Enable/disable log types", func(fs *pflag.FlagSet) {
		allSupportedTypes := api.SupportedCloudWatchClusterLogTypes()

		fs.StringSliceVar(&typesEnabled, "enable-types", []string{}, fmt.Sprintf("Log types to be enabled. Supported log types: (all, none, %s)", strings.Join(allSupportedTypes, ", ")))
		fs.StringSliceVar(&typesDisabled, "disable-types", []string{}, fmt.Sprintf("Log types to be disabled, the rest will be disabled. Supported log types: (all, none, %s)", strings.Join(allSupportedTypes, ", ")))

	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
}

func doEnableLogging(cmd *cmdutils.Cmd, logTypesToEnable []string, logTypesToDisable []string) error {
	if err := cmdutils.NewUtilsEnableLoggingLoader(cmd).Load(); err != nil {
		return err
	}

	if !cmd.ClusterConfig.HasClusterCloudWatchLogging() {
		if err := validateLoggingFlags(logTypesToEnable, logTypesToDisable); err != nil {
			return err
		}
	}

	cfg := cmd.ClusterConfig
	meta := cmd.ClusterConfig.Metadata

	printer := printers.NewJSONPrinter()

	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}
	cmdutils.LogRegionAndVersionInfo(meta)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	currentlyEnabled, _, err := ctl.GetCurrentClusterConfigForLogging(cfg)
	if err != nil {
		return err
	}

	var willBeEnabled sets.String
	if cfg.HasClusterCloudWatchLogging() {
		willBeEnabled = sets.NewString(cfg.CloudWatch.ClusterLogging.EnableTypes...)
	} else {
		baselineEnabled := currentlyEnabled.List()
		willBeEnabled = processTypesToEnable(baselineEnabled, logTypesToEnable, logTypesToDisable)
	}

	cfg.CloudWatch.ClusterLogging.EnableTypes = willBeEnabled.List()
	willBeDisabled := sets.NewString(api.SupportedCloudWatchClusterLogTypes()...).Difference(willBeEnabled)
	updateRequired := !currentlyEnabled.Equal(willBeEnabled)

	if err = printer.LogObj(logger.Debug, "cfg.json = \\\n%s\n", cfg); err != nil {
		return err
	}

	if updateRequired {
		describeTypesToEnable := "no types to enable"
		if len(willBeEnabled.List()) > 0 {
			describeTypesToEnable = fmt.Sprintf("enable types: %s", strings.Join(willBeEnabled.List(), ", "))
		}

		describeTypesToDisable := "no types to disable"
		if len(willBeDisabled.List()) > 0 {
			describeTypesToDisable = fmt.Sprintf("disable types: %s", strings.Join(willBeDisabled.List(), ", "))
		}

		cmdutils.LogIntendedAction(cmd.Plan, "update CloudWatch logging for cluster %q in %q (%s & %s)",
			meta.Name, meta.Region, describeTypesToEnable, describeTypesToDisable,
		)
		if !cmd.Plan {
			if err := ctl.UpdateClusterConfigForLogging(cfg); err != nil {
				return err
			}
		}
	} else {
		logger.Success("CloudWatch logging for cluster %q in %q is already up-to-date", meta.Name, meta.Region)
	}

	cmdutils.LogPlanModeWarning(cmd.Plan && updateRequired)

	return nil
}

func validateLoggingFlags(toEnable []string, toDisable []string) error {
	// At least enable-types or disable-types should be provided
	if len(toEnable) == 0 && len(toDisable) == 0 {
		return fmt.Errorf("at least one flag has to be provided: --enable-types, --disable-types")
	}

	isEnableAll := len(toEnable) == 1 && toEnable[0] == "all"
	isDisableAll := len(toDisable) == 1 && toDisable[0] == "all"

	// Can't enable all and disable all
	if isDisableAll && isEnableAll {
		return fmt.Errorf("cannot use `all` for both --enable-types and --disable-types at the same time")
	}

	// Check all are valid values
	// TODO if this is too restrictive we can drop it
	if err := checkAllTypesAreSupported(toEnable); err != nil {
		return err
	}
	if err := checkAllTypesAreSupported(toDisable); err != nil {
		return err
	}
	// both options are provided but without "all"
	toEnableSet := sets.NewString(toEnable...)
	toDisableSet := sets.NewString(toDisable...)

	appearInBoth := toEnableSet.Intersection(toDisableSet)

	if appearInBoth.Len() != 0 {
		return fmt.Errorf("log types cannot be part of --enable-types and --disable-types simultaneously")
	}
	return nil
}

func processTypesToEnable(existingEnabled, toEnable, toDisable []string) sets.String {
	emptyToEnable := toEnable == nil || len(toEnable) == 0
	emptyToDisable := toDisable == nil || len(toDisable) == 0

	isEnableAll := !emptyToEnable && toEnable[0] == "all"
	isDisableAll := !emptyToDisable && toDisable[0] == "all"

	// When all is provided in one of the options
	if isDisableAll {
		return sets.NewString(toEnable...)
	}
	if isEnableAll {
		toDisableSet := sets.NewString(toDisable...)
		toEnableSet := sets.NewString(api.SupportedCloudWatchClusterLogTypes()...).Difference(toDisableSet)
		return toEnableSet
	}

	// willEnable = existing - toDisable + toEnable
	willEnable := sets.NewString(existingEnabled...)
	willEnable.Insert(toEnable...)
	willEnable.Delete(toDisable...)

	return willEnable
}

func checkAllTypesAreSupported(logTypes []string) error {
	if len(logTypes) == 1 && logTypes[0] == "all" {
		return nil
	}
	allSupportedTypesSet := sets.NewString(api.SupportedCloudWatchClusterLogTypes()...)
	for _, logType := range logTypes {
		if !allSupportedTypesSet.Has(logType) {
			return fmt.Errorf("unknown log type %s. Supported log types: all, %s", logType, strings.Join(api.SupportedCloudWatchClusterLogTypes(), ", "))
		}
	}
	return nil
}
