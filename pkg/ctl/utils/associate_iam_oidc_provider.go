package utils

import (
	"github.com/kris-nova/logger"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func associateIAMOIDCProviderCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("associate-iam-oidc-provider", "Setup IAM OIDC provider for a cluster to enable IAM roles for pods", "")

	cmd.SetRunFuncWithNameArg(func() error {
		return doAssociateIAMOIDCProvider(cmd)
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlagWithDeprecated(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cmd)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
}

func doAssociateIAMOIDCProvider(cmd *cmdutils.Cmd) error {
	if err := cmdutils.NewUtilsAssociateIAMOIDCProviderLoader(cmd).Load(); err != nil {
		return err
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

	if !providerExists {
		cmdutils.LogIntendedAction(cmd.Plan, "create IAM Open ID Connect provider for cluster %q in %q", meta.Name, meta.Region)
		if !cmd.Plan {
			if err := oidc.CreateProvider(); err != nil {
				return err
			}
			logger.Success("created IAM Open ID Connect provider for cluster %q in %q", meta.Name, meta.Region)
		}
	} else {
		logger.Info("IAM Open ID Connect provider is already associated with cluster %q in %q", meta.Name, meta.Region)
	}

	cmdutils.LogPlanModeWarning(cmd.Plan && !providerExists)

	return nil
}
