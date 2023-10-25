package utils

import (
	"context"

	"github.com/kris-nova/logger"
	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func updateClusterVPCConfigWithHandler(cmd *cmdutils.Cmd, handler func(cmd *cmdutils.Cmd) error) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("update-cluster-vpc-config", "Update Kubernetes API endpoint access configuration and public access CIDRs",
		dedent.Dedent(`Updates the Kubernetes API endpoint access configuration and public access CIDRs.

			When a config file is passed, only changes to vpc.clusterEndpoints and vpc.publicAccessCIDRs are updated in the EKS API.
			Changes to any other fields are ignored.
		`),
	)
	var options cmdutils.UpdateClusterVPCOptions
	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		if err := cmdutils.NewUpdateClusterVPCLoader(cmd, options).Load(); err != nil {
			return err
		}
		return handler(cmd)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cmd)
	})

	cmd.FlagSetGroup.InFlagSet("Endpoint Access", func(fs *pflag.FlagSet) {
		fs.BoolVar(&options.PrivateAccess, "private-access", false, "access for private (VPC) clients")
		fs.BoolVar(&options.PublicAccess, "public-access", false, "access for public clients")
	})
	cmd.FlagSetGroup.InFlagSet("Public Access CIDRs", func(fs *pflag.FlagSet) {
		fs.StringSliceVar(&options.PublicAccessCIDRs, "public-access-cidrs", nil, "CIDR blocks that EKS uses to create a security group on the public endpoint")
	})
	cmd.FlagSetGroup.InFlagSet("Control plane subnets and security groups", func(fs *pflag.FlagSet) {
		fs.StringSliceVar(&options.ControlPlaneSubnetIDs, "control-plane-subnet-ids", nil, "Subnet IDs for the control plane")
		fs.StringSliceVar(&options.ControlPlaneSecurityGroupIDs, "control-plane-security-group-ids", nil, "Security group IDs for the control plane")
	})

	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, false)
}

func updateClusterVPCConfigCmd(cmd *cmdutils.Cmd) {
	updateClusterVPCConfigWithHandler(cmd, doUpdateClusterVPCConfig)
}

func doUpdateClusterVPCConfig(cmd *cmdutils.Cmd) error {
	ctx := context.Background()
	ctl, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}
	cfg := cmd.ClusterConfig
	logger.Info("using region %s", cfg.Metadata.Region)

	if ok, err := ctl.CanUpdate(cfg); !ok {
		return err
	}

	vpcHelper := &VPCHelper{
		VPCUpdater:  ctl,
		ClusterMeta: cfg.Metadata,
		Cluster:     ctl.Status.ClusterInfo.Cluster,
		PlanMode:    cmd.Plan,
	}

	return vpcHelper.UpdateClusterVPCConfig(ctx, cfg.VPC)
}
