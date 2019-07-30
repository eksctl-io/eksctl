package utils

import (
	"fmt"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"k8s.io/apimachinery/pkg/util/sets"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func enableLoggingCmd(rc *cmdutils.ResourceCmd) {
	cfg := api.NewClusterConfig()
	rc.ClusterConfig = cfg

	rc.SetDescription("enable-logging", "Update cluster logging configuration", "")

	rc.SetRunFuncWithNameArg(func() error {
		return doEnableLogging(rc)
	})

	rc.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddNameFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, rc.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &rc.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, rc)
	})

	rc.FlagSetGroup.InFlagSet("Enable/disable log types", func(fs *pflag.FlagSet) {
		enableExplicitly := sets.NewString()
		allSupportedTypes := api.SupportedCloudWatchClusterLogTypes()
		inherit := sets.NewString()

		enableAll := fs.Bool("all", true, fmt.Sprintf("Enable all supported log types (%s)", strings.Join(allSupportedTypes, ", ")))

		for _, logType := range allSupportedTypes {
			_ = fs.Bool(logType, false, fmt.Sprintf("Enable %q log type", logType))
		}

		cmdutils.AddPreRun(rc.Command, func(cmd *cobra.Command, args []string) {
			if *enableAll {
				inherit.Insert(allSupportedTypes...)
			}

			for _, logType := range allSupportedTypes {
				f := cmd.Flag(logType)

				shouldEnable := f.Value.String() == "true"
				shouldDisable := f.Changed && !shouldEnable

				if shouldEnable {
					enableExplicitly.Insert(logType)
					if !cmd.Flag("all").Changed {
						// do not include all log types if --all wasn't specified explicitly
						inherit = sets.NewString()
					}
				}

				if shouldDisable {
					inherit.Delete(logType)
				}
			}

			cfg.AppendClusterCloudWatchLogTypes(enableExplicitly.List()...)
			cfg.AppendClusterCloudWatchLogTypes(inherit.List()...)
		})
	})

	cmdutils.AddCommonFlagsForAWS(rc.FlagSetGroup, rc.ProviderConfig, false)
}

func doEnableLogging(rc *cmdutils.ResourceCmd) error {
	if err := cmdutils.NewUtilsEnableLoggingLoader(rc).Load(); err != nil {
		return err
	}

	cfg := rc.ClusterConfig
	meta := rc.ClusterConfig.Metadata

	if err := api.SetClusterConfigDefaults(cfg); err != nil {
		return err
	}

	printer := printers.NewJSONPrinter()
	ctl := eks.New(rc.ProviderConfig, cfg)

	if !ctl.IsSupportedRegion() {
		return cmdutils.ErrUnsupportedRegion(rc.ProviderConfig)
	}
	logger.Info("using region %s", meta.Region)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	currentlyEnabled, _, err := ctl.GetCurrentClusterConfigForLogging(meta)
	if err != nil {
		return err
	}

	shouldEnable := sets.NewString()

	if cfg.HasClusterCloudWatchLogging() {
		shouldEnable.Insert(cfg.CloudWatch.ClusterLogging.EnableTypes...)
	}

	shouldDisable := sets.NewString(api.SupportedCloudWatchClusterLogTypes()...).Difference(shouldEnable)

	updateRequired := !currentlyEnabled.Equal(shouldEnable)

	if err := printer.LogObj(logger.Debug, "cfg.json = \\\n%s\n", cfg); err != nil {
		return err
	}

	if updateRequired {
		describeTypesToEnable := "no types to enable"
		if len(shouldEnable.List()) > 0 {
			describeTypesToEnable = fmt.Sprintf("enable types: %s", strings.Join(shouldEnable.List(), ", "))
		}

		describeTypesToDisable := "no types to disable"
		if len(shouldDisable.List()) > 0 {
			describeTypesToDisable = fmt.Sprintf("disable types: %s", strings.Join(shouldDisable.List(), ", "))
		}

		cmdutils.LogIntendedAction(rc.Plan, "update CloudWatch logging for cluster %q in %q (%s & %s)",
			meta.Name, meta.Region, describeTypesToEnable, describeTypesToDisable,
		)
		if !rc.Plan {
			if err := ctl.UpdateClusterConfigForLogging(cfg); err != nil {
				return err
			}
		}
	} else {
		logger.Success("CloudWatch logging for cluster %q in %q is already up-to-date", meta.Name, meta.Region)
	}

	cmdutils.LogPlanModeWarning(rc.Plan && updateRequired)

	return nil
}
