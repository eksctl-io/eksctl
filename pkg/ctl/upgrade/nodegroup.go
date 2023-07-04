package upgrade

import (
	"context"
	"time"

	"github.com/aws/amazon-ec2-instance-selector/v2/pkg/selector"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

const upgradeNodegroupTimeout = 45 * time.Minute

func upgradeNodeGroupCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("nodegroup", "Upgrade nodegroup", "")

	var options nodegroup.UpgradeOptions
	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return upgradeNodeGroup(cmd, options)
	}

	cmd.FlagSetGroup.InFlagSet("Nodegroup", func(fs *pflag.FlagSet) {
		fs.StringVar(&options.NodegroupName, "name", "", "Nodegroup name")
		fs.StringVar(&options.LaunchTemplateVersion, "launch-template-version", "", "Launch template version")
		fs.StringVar(&options.KubernetesVersion, "kubernetes-version", "", "Kubernetes version")
		fs.BoolVar(&options.ForceUpgrade, "force-upgrade", false, "Force the update if the existing node group's pods are unable to be drained due to a pod disruption budget issue")
		fs.StringVar(&options.ReleaseVersion, "release-version", "", "AMI version of the EKS optimized AMI to use")
		fs.BoolVar(&options.Wait, "wait", true, "nodegroup upgrade to complete")
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cmd.ClusterConfig.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		// found with experimentation
		cmdutils.AddTimeoutFlagWithValue(fs, &cmd.ProviderConfig.WaitTimeout, upgradeNodegroupTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, false)

}

func upgradeNodeGroup(cmd *cmdutils.Cmd, options nodegroup.UpgradeOptions) error {
	cfg := cmd.ClusterConfig
	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet(cmdutils.ClusterNameFlag(cmd))
	}

	if options.NodegroupName != "" && cmd.NameArg != "" {
		return cmdutils.ErrFlagAndArg("--name", options.NodegroupName, cmd.NameArg)
	}

	if cmd.NameArg != "" {
		options.NodegroupName = cmd.NameArg
	}

	if options.NodegroupName == "" {
		return cmdutils.ErrMustBeSet("name")
	}

	ctx := context.TODO()
	ctl, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}

	if ok, err := ctl.CanOperate(cfg); !ok {
		return err
	}

	clientSet, err := ctl.NewStdClientSet(cfg)
	if err != nil {
		return err
	}

	instanceSelector, err := selector.New(ctx, ctl.AWSProvider.AWSConfig())
	if err != nil {
		return err
	}
	return nodegroup.New(cfg, ctl, clientSet, instanceSelector).Upgrade(ctx, options)
}
