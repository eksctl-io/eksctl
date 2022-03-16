package set

import (
	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/eksctl/pkg/actions/label"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/managed"
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

		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)
}

func setLabels(cmd *cmdutils.Cmd, options labelOptions) error {
	// if nodeGroupName is defined, the loader will filter managed nodegroups
	// for that nodegroup only.
	if err := cmdutils.NewSetLabelLoader(cmd, options.nodeGroupName, options.labels).Load(); err != nil {
		return err
	}
	cfg := cmd.ClusterConfig
	ctl, err := cmd.NewProviderForExistingCluster()
	if err != nil {
		return err
	}

	service := managed.NewService(ctl.Provider.EKS(), ctl.Provider.EC2(), manager.NewStackCollection(ctl.Provider, cfg), cfg.Metadata.Name)

	if options.nodeGroupName == "" && cmd.ClusterConfigFile != "" {
		logger.Info("setting label(s) on %d nodegroup(s) in cluster %s", len(cfg.ManagedNodeGroups), cmd.ClusterConfig.Metadata)
	} else if options.nodeGroupName != "" {
		logger.Info("setting label(s) on nodegroup %s in cluster %s", options.nodeGroupName, cmd.ClusterConfig.Metadata)
	}

	manager := label.New(cfg.Metadata.Name, service, ctl.Provider.EKS())
	// when there is no config file provided
	if cmd.ClusterConfigFile == "" {
		if err := manager.Set(options.nodeGroupName, options.labels); err != nil {
			return err
		}
		logger.Info("done")
		return nil
	}
	// when there is a config file, we call GetLabels first.
	for _, mng := range cfg.ManagedNodeGroups {
		existingLabels, err := service.GetLabels(mng.Name)
		if err != nil {
			return err
		}
		for k := range existingLabels {
			delete(mng.Labels, k)
		}
		if len(mng.Labels) == 0 {
			logger.Info("no new labels to add for nodegroup %s", mng.Name)
			continue
		}
		if err := manager.Set(mng.Name, mng.Labels); err != nil {
			return err
		}
	}

	logger.Info("done")
	return nil
}
