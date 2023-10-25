package utils

import (
	"context"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func publicAccessCIDRsCmdWithHandler(cmd *cmdutils.Cmd, handler func(cmd *cmdutils.Cmd) error) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.CobraCommand.Deprecated = "this command is deprecated and will be removed soon. Use `eksctl utils update-cluster-vpc-config --public-access-cidrs=<> instead."
	cmd.SetDescription("set-public-access-cidrs", "Update public access CIDRs", "CIDR blocks that EKS uses to create a security group on the public endpoint")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		if err := cmdutils.NewUtilsPublicAccessCIDRsLoader(cmd).Load(); err != nil {
			return err
		}
		return handler(cmd)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cmd)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, false)
}

func publicAccessCIDRsCmd(cmd *cmdutils.Cmd) {
	publicAccessCIDRsCmdWithHandler(cmd, doUpdatePublicAccessCIDRs)
}

func doUpdatePublicAccessCIDRs(cmd *cmdutils.Cmd) error {
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

	cfg.VPC.ClusterEndpoints = nil
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
