package upgrade

import (
	"time"

	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/managed"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

const upgradeNodegroupTimeout = 35 * time.Minute

func upgradeNodeGroupCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("nodegroup", "Upgrade nodegroup", "")

	var options managed.UpgradeOptions
	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return upgradeNodeGroup(cmd, options)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name")
		fs.StringVar(&options.NodegroupName, "name", "", "Nodegroup name")
		fs.StringVar(&options.KubernetesVersion, "kubernetes-version", "", "Kubernetes version")
		fs.StringVar(&options.LaunchTemplateVersion, "launch-template-version", "", "Launch template version")
		fs.BoolVar(&options.ForceUpgrade, "force-upgrade", false, "Force the update if the existing node group's pods are unable to be drained due to a pod disruption budget issue")

		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmd.Wait = true
		cmdutils.AddWaitFlag(fs, &cmd.Wait, "nodegroup upgrade to complete")

		// found with experimentation
		cmdutils.AddTimeoutFlagWithValue(fs, &cmd.ProviderConfig.WaitTimeout, upgradeNodegroupTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)

}

func upgradeNodeGroup(cmd *cmdutils.Cmd, options managed.UpgradeOptions) error {
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

	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if ok, err := ctl.CanOperate(cfg); !ok {
		return err
	}

	clientSet, err := ctl.NewStdClientSet(cfg)
	if err != nil {
		return err
	}

	return nodegroup.New(cfg, ctl, clientSet).Upgrade(options.NodegroupName, options.KubernetesVersion, options.LaunchTemplateVersion, options.ForceUpgrade)

}
