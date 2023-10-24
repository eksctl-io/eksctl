package utils

import (
	"context"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func updateClusterEndpointsCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.CobraCommand.Deprecated = "this command is deprecated and will be removed soon. Use `eksctl utils update-cluster-vpc-config --public-access=<> --private-access=<> instead."
	cmd.SetDescription("update-cluster-endpoints", "Update Kubernetes API endpoint access configuration", "")

	var (
		private bool
		public  bool
	)
	cmd.CobraCommand.RunE = func(_ *cobra.Command, _ []string) error {
		return doUpdateClusterEndpoints(cmd, private, public)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlagWithDeprecated(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cmd)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmd.FlagSetGroup.InFlagSet("Endpoint Access",
		func(fs *pflag.FlagSet) {
			fs.BoolVar(&private, "private-access", false, "access for private (VPC) clients")
			fs.BoolVar(&public, "public-access", false, "access for public clients")
		})
	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, false)
}

func doUpdateClusterEndpoints(cmd *cmdutils.Cmd, newPrivate bool, newPublic bool) error {
	if err := cmdutils.NewUtilsEnableEndpointAccessLoader(cmd, newPrivate, newPublic).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig
	meta := cmd.ClusterConfig.Metadata

	ctx := context.TODO()
	ctl, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}
	logger.Info("using region %s", meta.Region)

	if ok, err := ctl.CanUpdate(cfg); !ok {
		return err
	}

	cfg.VPC.PublicAccessCIDRs = nil
	cfg.VPC.ControlPlaneSubnetIDs = nil
	cfg.VPC.ControlPlaneSecurityGroupIDs = nil
	vpcHelper := &VPCHelper{
		VPCUpdater:  ctl,
		ClusterMeta: cfg.Metadata,
		Cluster:     ctl.Status.ClusterInfo.Cluster,
		PlanMode:    cmd.Plan,
	}

	return vpcHelper.UpdateClusterVPCConfig(ctx, cfg.VPC)
}
