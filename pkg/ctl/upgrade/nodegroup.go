package upgrade

import (
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/managed"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

type upgradeOptions struct {
	nodeGroupName     string
	kubernetesVersion string
}

func upgradeNodeGroupCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("nodegroup", "Upgrade nodegroup", "")

	var options upgradeOptions
	cmd.SetRunFuncWithNameArg(func() error {
		return upgradeNodeGroup(cmd, options)
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVarP(&cfg.Metadata.Name, "cluster", "", "", "EKS cluster name")
		fs.StringVarP(&options.nodeGroupName, "name", "", "", "Nodegroup name")
		fs.StringVarP(&options.kubernetesVersion, "kubernetes-version", "", "", "Kubernetes version")
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)

		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)

}

func upgradeNodeGroup(cmd *cmdutils.Cmd, options upgradeOptions) error {
	cfg := cmd.ClusterConfig
	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet(cmdutils.ClusterNameFlag(cmd))
	}

	if options.nodeGroupName != "" && cmd.NameArg != "" {
		return cmdutils.ErrFlagAndArg("--name", options.nodeGroupName, cmd.NameArg)
	}

	if cmd.NameArg != "" {
		options.nodeGroupName = cmd.NameArg
	}

	if options.nodeGroupName == "" {
		return cmdutils.ErrMustBeSet("name")
	}

	ctl := eks.New(cmd.ProviderConfig, cmd.ClusterConfig)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	stackCollection := manager.NewStackCollection(ctl.Provider, cfg)
	managedService := managed.NewService(ctl.Provider, stackCollection, cfg.Metadata.Name)
	return managedService.UpgradeNodeGroup(options.nodeGroupName, options.kubernetesVersion)
}
