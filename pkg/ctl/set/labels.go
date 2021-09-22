package set

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

type labelOptions struct {
	nodeGroupName string
	labels        map[string]string
}

func setLabelsCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("labels", "Create or overwrite labels for managed nodegroups", "")

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

		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)

}

func setLabels(cmd *cmdutils.Cmd, options labelOptions) error {
	cfg := cmd.ClusterConfig
	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet(cmdutils.ClusterNameFlag(cmd))
	}
	if options.nodeGroupName == "" {
		return cmdutils.ErrMustBeSet("--nodegroup")
	}

	if cmd.NameArg != "" {
		return cmdutils.ErrUnsupportedNameArg()
	}

	ctl, err := cmd.NewProviderForExistingCluster()
	if err != nil {
		return err
	}

	cmdutils.LogRegionAndVersionInfo(cmd.ClusterConfig.Metadata)
	logger.Info("setting label(s) on nodegroup %s in cluster %s", options.nodeGroupName, cmd.ClusterConfig.Metadata)

	service := managed.NewService(ctl.Provider.EKS(), ctl.Provider.SSM(), ctl.Provider.EC2(), manager.NewStackCollection(ctl.Provider, cfg), cfg.Metadata.Name)
	manager := label.New(cfg.Metadata.Name, service, ctl.Provider.EKS())
	if err := manager.Set(options.nodeGroupName, options.labels); err != nil {
		return err
	}

	logger.Info("done")
	return nil
}
