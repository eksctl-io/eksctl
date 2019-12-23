package utils

import (
	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var (
	private bool
	public  bool
)

func updateClusterEndpointsCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("update-cluster-endpoints", "Update Kubernetes API endpoint access configuration", "")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, _ []string) error {
		return doUpdateClusterEndpoints(cmd, private, public)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlagWithDeprecated(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cmd)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmd.FlagSetGroup.InFlagSet("Update private/public Kubernetes API endpoint access configuration",
		func(fs *pflag.FlagSet) {
			fs.BoolVar(&private, "private-access", false, "access for private (VPC) clients")
			fs.BoolVar(&public, "public-access", false, "access for public clients")
		})
	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
}

func accessFlagsSet(cmd *cmdutils.Cmd) (privateSet, publicSet bool) {
	cmd.FlagSetGroup.InFlagSet("Update private/public Kubernetes API endpoint access configuration",
		func(fs *pflag.FlagSet) {
			if priv := fs.Lookup("private-access"); priv != nil {
				privateSet = priv.Changed
			}
			if pub := fs.Lookup("public-access"); pub != nil {
				publicSet = pub.Changed
			}
		})
	return
}

func doUpdateClusterEndpoints(cmd *cmdutils.Cmd, newPrivate bool, newPublic bool) error {
	if err := cmdutils.NewUtilsEnableEndpointAccessLoader(cmd).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig
	meta := cmd.ClusterConfig.Metadata

	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}
	logger.Info("using region %s", meta.Region)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if ok, err := ctl.CanUpdate(cfg); !ok {
		return err
	}

	clusterVPCConfig, err := ctl.GetCurrentClusterVPCConfig(cfg)
	if err != nil {
		return err
	}

	curPrivate, curPublic := *clusterVPCConfig.ClusterEndpoints.PrivateAccess, *clusterVPCConfig.ClusterEndpoints.PublicAccess

	logger.Info("current Kubernetes API endpoint access: privateAccess=%v, publicAccess=%v",
		curPrivate, curPublic)

	if cfg.VPC.ClusterEndpoints.PublicAccess != nil && cfg.VPC.ClusterEndpoints.PrivateAccess != nil {
		newPublic = *cfg.VPC.ClusterEndpoints.PublicAccess
		newPrivate = *cfg.VPC.ClusterEndpoints.PrivateAccess
	}else{
		privateSet, publicSet := accessFlagsSet(cmd)
		if !privateSet {
			newPrivate = curPrivate
		}
		if !publicSet {
			newPublic = curPublic
		}
	}

	// Nothing changed?
	if newPrivate == curPrivate && newPublic == curPublic {
		logger.Success("Kubernetes API endpoint access for cluster %q in %q is already up to date",
			meta.Name, meta.Region)
		return nil
	}

	cfg.VPC.ClusterEndpoints.PrivateAccess = &newPrivate
	cfg.VPC.ClusterEndpoints.PublicAccess = &newPublic

	cmdutils.LogIntendedAction(
		cmd.Plan, "update Kubernetes API endpoint access for cluster %q in %q to: privateAccess=%v, publicAccess=%v",
		meta.Name, meta.Region, newPrivate, newPublic)

	if err := cfg.ValidateClusterEndpointConfig(); err != nil {
		// Error for everything except private-only (which leaves the cluster accessible)
		if err != api.ErrClusterEndpointPrivateOnly {
			return err
		}
		logger.Warning(err.Error())
	}

	if !cmd.Plan {
		if err := ctl.UpdateClusterConfigForEndpoints(cfg); err != nil {
			return err
		}
		cmdutils.LogCompletedAction(
			false,
			"the Kubernetes API endpoint access for cluster %q in %q has been updated to: "+
				"privateAccess=%v, publicAccess=%v",
			meta.Name, meta.Region, newPrivate, newPublic)
	}
	cmdutils.LogPlanModeWarning(cmd.Plan)

	return nil
}
