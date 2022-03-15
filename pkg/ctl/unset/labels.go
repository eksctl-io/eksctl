package unset

import (
	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/eksctl/pkg/actions/label"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/managed"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func unsetLabelsCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("labels", "Remove labels from managed nodegroups", "")

	var (
		nodeGroupName string
		removeLabels  []string
	)
	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return unsetLabels(cmd, nodeGroupName, removeLabels)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name")
		fs.StringVarP(&nodeGroupName, "nodegroup", "n", "", "Nodegroup name")
		fs.StringSliceVarP(&removeLabels, "labels", "l", nil, "List of labels to remove")

		_ = cobra.MarkFlagRequired(fs, "labels")

		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)

}

func unsetLabels(cmd *cmdutils.Cmd, nodeGroupName string, removeLabels []string) error {
	cfg := cmd.ClusterConfig
	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet(cmdutils.ClusterNameFlag(cmd))
	}
	if nodeGroupName == "" {
		return cmdutils.ErrMustBeSet("--nodegroup")
	}

	if cmd.NameArg != "" {
		return cmdutils.ErrUnsupportedNameArg()
	}

	ctl, err := cmd.NewProviderForExistingCluster()
	if err != nil {
		return err
	}

	logger.Info("removing label(s) from nodegroup %s in cluster %s", nodeGroupName, cmd.ClusterConfig.Metadata)
	service := managed.NewService(ctl.Provider.EKS(), ctl.Provider.EC2(), manager.NewStackCollection(ctl.Provider, cfg), cfg.Metadata.Name)
	manager := label.New(cfg.Metadata.Name, service, ctl.Provider.EKS())
	if err := manager.Unset(nodeGroupName, removeLabels); err != nil {
		return err
	}

	logger.Info("done")
	return nil
}
