package utils

import (
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/sets"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func publicAccessCIDRsCmdWithHandler(cmd *cmdutils.Cmd, handler func(cmd *cmdutils.Cmd) error) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

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
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cmd)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
}

func publicAccessCIDRsCmd(cmd *cmdutils.Cmd) {
	publicAccessCIDRsCmdWithHandler(cmd, doUpdatePublicAccessCIDRs)
}

func doUpdatePublicAccessCIDRs(cmd *cmdutils.Cmd) error {
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

	logger.Info("current public access CIDRs: %v", clusterVPCConfig.PublicAccessCIDRs)

	if cidrsEqual(clusterVPCConfig.PublicAccessCIDRs, cfg.VPC.PublicAccessCIDRs) {
		logger.Success("Public Endpoint Restrictions for cluster %q in %q is already up to date",
			meta.Name, meta.Region)
		return nil
	}

	cmdutils.LogIntendedAction(
		cmd.Plan, "update Public Endpoint Restrictions for cluster %q in %q to: %v",
		meta.Name, meta.Region, cfg.VPC.PublicAccessCIDRs)

	if !cmd.Plan {
		if err := ctl.UpdatePublicAccessCIDRs(cfg); err != nil {
			return errors.Wrap(err, "error updating CIDRs for public access")
		}
		cmdutils.LogCompletedAction(
			false,
			"Public Endpoint Restrictions for cluster %q in %q have been updated to: %v",
			meta.Name, meta.Region, cfg.VPC.PublicAccessCIDRs)
	}
	cmdutils.LogPlanModeWarning(cmd.Plan)
	return nil
}

func cidrsEqual(currentValues, newValues []string) bool {
	return sets.NewString(currentValues...).Equal(sets.NewString(newValues...))
}
