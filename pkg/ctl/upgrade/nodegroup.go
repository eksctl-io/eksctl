package upgrade

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/managed"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

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
		fs.StringVarP(&cfg.Metadata.Name, "cluster", "", "", "EKS cluster name")
		fs.StringVarP(&options.NodegroupName, "name", "", "", "Nodegroup name")
		fs.StringVarP(&options.KubernetesVersion, "kubernetes-version", "", "", "Kubernetes version")
		fs.StringVarP(&options.LaunchTemplateVersion, "launch-template-version", "", "", "Launch template version")

		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmd.Wait = true
		cmdutils.AddWaitFlag(fs, &cmd.Wait, "nodegroup upgrade to complete")

		// during testing, a nodegroup update took 20-25 minutes
		cmdutils.AddTimeoutFlagWithValue(fs, &cmd.ProviderConfig.WaitTimeout, 35*time.Minute)
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

	ctl := eks.New(&cmd.ProviderConfig, cmd.ClusterConfig)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	stackCollection := manager.NewStackCollection(ctl.Provider, cfg)
	managedService := managed.NewService(ctl.Provider, stackCollection, cfg.Metadata.Name)
	return managedService.UpgradeNodeGroup(options)
}
