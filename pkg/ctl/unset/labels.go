package unset

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/actions/label"
	"github.com/weaveworks/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
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

	ctl := eks.New(&cmd.ProviderConfig, cmd.ClusterConfig)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	cmdutils.LogRegionAndVersionInfo(cmd.ClusterConfig.Metadata)
	logger.Info("removing label(s) from nodegroup %s in cluster %s", nodeGroupName, cmd.ClusterConfig.Metadata)

	manager := label.New(cfg, ctl)
	if err := manager.Unset(nodeGroupName, removeLabels); err != nil {
		return err
	}

	logger.Info("done")
	return nil
}
