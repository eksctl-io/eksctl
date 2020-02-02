package unset

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/managed"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

func unsetLabelsCmd(cmd *cmdutils.Cmd) {
	unsetLabelsWithRunFunc(cmd, func(cmd *cmdutils.Cmd, nodeGroupName string, removeLabels []string) error {
		return unsetLabels(cmd, nodeGroupName, removeLabels)
	})
}

func unsetLabelsWithRunFunc(cmd *cmdutils.Cmd, runFunc func(cmd *cmdutils.Cmd, nodeGroupName string, removeLabels []string) error) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("labels", "Remove Labels", "")

	var (
		nodeGroupName string
		removeLabels  []string
	)
	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return runFunc(cmd, nodeGroupName, removeLabels)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name")
		fs.StringVarP(&nodeGroupName, "nodegroup", "n", "", "Nodegroup name")
		fs.StringSliceVarP(&removeLabels, "labels", "l", nil, "List of labels to remove")

		_ = cobra.MarkFlagRequired(fs, "labels")

		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)

}

func unsetLabels(cmd *cmdutils.Cmd, nodeGroupName string, removeLabels []string) error {
	cfg := cmd.ClusterConfig
	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet(cmdutils.ClusterNameFlag(cmd))
	}

	if cmd.NameArg != "" {
		return cmdutils.ErrUnsupportedNameArg()
	}

	ctl := eks.New(cmd.ProviderConfig, cmd.ClusterConfig)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	stackCollection := manager.NewStackCollection(ctl.Provider, cfg)
	managedService := managed.NewService(ctl.Provider, stackCollection, cfg.Metadata.Name)
	return managedService.UpdateLabels(nodeGroupName, nil, removeLabels)
}
