package set

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/managed"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

type labelOptions struct {
	nodeGroupName string
	labels        map[string]string
}

func setLabelsCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("labels", "Create or overwrite labels", "")

	var options labelOptions
	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return setLabels(cmd, options)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name")
		fs.StringVarP(&options.nodeGroupName, "nodegroup", "n", "", "Nodegroup name")
		cmdutils.AddStringToStringVarPFlag(fs, &options.labels, "labels", "l", nil, "Labels")
		_ = cobra.MarkFlagRequired(fs, "labels")

		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)

}

func setLabels(cmd *cmdutils.Cmd, options labelOptions) error {
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
	return managedService.UpdateLabels(options.nodeGroupName, options.labels, nil)
}
