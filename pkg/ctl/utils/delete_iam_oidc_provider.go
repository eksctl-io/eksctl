package utils

import (
	"errors"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func deleteIAMOIDCProviderCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cfg.IAM.WithOIDC = api.Disabled()

	cmd.SetDescription("delete-iam-oidc-provider", "Delete IAM OIDC provider for a cluster", "")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doDeleteIAMOIDCProvider(cmd)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlagWithDeprecated(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cmd)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
}

func doDeleteIAMOIDCProvider(cmd *cmdutils.Cmd) error {
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

	if ok, err := ctl.CanOperate(cfg); !ok {
		return err
	}

	oidc, err := ctl.NewOpenIDConnectManager(cfg)
	if err != nil {
		return err
	}

	if err := printer.LogObj(logger.Debug, "cfg.json = \\\n%s\n", cfg); err != nil {
		return err
	}

	providerExists, err := oidc.CheckProviderExists()
	if err != nil {
		return err
	}

	if providerExists {

		stackManager := ctl.NewStackManager(cfg)

		existing, err := stackManager.ListIAMServiceAccountStacks()
		if err != nil {
			return err
		}

		if len(existing) > 0 {
			logger.Warning("found existing iamserviceaccount(s); can't delete IAM OIDC provider associated with cluster %q in %q", meta.Name, meta.Region)
			return errors.New("unable to delete IAM OIDC provider with existing iamserviceaccount(s)")
		}

		cmdutils.LogIntendedAction(cmd.Plan, "delete IAM Open ID Connect provider for cluster %q in %q", meta.Name, meta.Region)
		if !cmd.Plan {
			if err := oidc.DeleteProvider(); err != nil {
				return err
			}
			logger.Success("deleted IAM Open ID Connect provider for cluster %q in %q", meta.Name, meta.Region)
		}
	} else {
		logger.Info("no IAM OIDC provider associated with cluster %q in %q", meta.Name, meta.Region)
	}

	cmdutils.LogPlanModeWarning(cmd.Plan && !providerExists)

	return nil
}
