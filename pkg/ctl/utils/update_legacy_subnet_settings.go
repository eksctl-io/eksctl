package utils

import (
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func updateLegacySubnetSettings(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("update-legacy-subnet-settings", "Update the configuration of the cluster's public subnets with MapPublicIpOnLaunch enabled",
		"MapPublicIpOnLaunch is a new property for subnets that is required for creating new nodegroups in them")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doUpdateLegacySubnetSettings(cmd)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlagWithDeprecated(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
}

func doUpdateLegacySubnetSettings(cmd *cmdutils.Cmd) error {
	cfg := cmd.ClusterConfig
	meta := cmd.ClusterConfig.Metadata

	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}
	cmdutils.LogRegionAndVersionInfo(cfg.Metadata)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if meta.Name != "" && cmd.NameArg != "" {
		return cmdutils.ErrFlagAndArg(cmdutils.ClusterNameFlag(cmd), meta.Name, cmd.NameArg)
	}
	if cmd.NameArg != "" {
		meta.Name = cmd.NameArg
	}

	if meta.Name == "" {
		return cmdutils.ErrMustBeSet(cmdutils.ClusterNameFlag(cmd))
	}

	if err := ctl.LoadClusterVPC(cfg); err != nil {
		return errors.Wrapf(err, "getting VPC configuration for cluster %q", cfg.Metadata.Name)
	}

	stackManager := ctl.NewStackManager(cfg)

	logger.Info("updating settings { MapPublicIpOnLaunch: enabled } for public subnets %q", cfg.PublicSubnetIDs())
	err = stackManager.EnsureMapPublicIPOnLaunchEnabled()
	if err != nil {
		logger.Warning(err.Error())
		return err
	}

	logger.Success("public subnets up to date")
	return nil
}
