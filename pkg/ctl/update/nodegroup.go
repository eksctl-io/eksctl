package update

import (
	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/managed"
)

func updateNodeGroupCmd(cmd *cmdutils.Cmd) {
	cmd.ClusterConfig = api.NewClusterConfig()

	cmd.SetDescription(
		"nodegroup",
		"Update nodegroup",
		dedent.Dedent(`Update nodegroup and its configuration.

		To upgrade a nodegroup, please use 'eksctl upgrade nodegroup' instead.
		Please consult the eksctl documentation for more info on which fields can be updated through 'eksctl update nodegroup'.
		Note that this is only available for managed nodegroups. 
	`),
	)

	var options managed.UpdateOptions
	cmd.FlagSetGroup.InFlagSet("Nodegroup", func(fs *pflag.FlagSet) {
		fs.StringVar(&options.NodegroupName, "name", "", "Nodegroup name")
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cmd.ClusterConfig.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return updateNodegroup(cmd, options)
	}
}

func updateNodegroup(cmd *cmdutils.Cmd, options managed.UpdateOptions) error {
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
		return cmdutils.ErrMustBeSet("--name")
	}

	if cmd.ClusterConfigFile == "" {
		return cmdutils.ErrMustBeSet("--config-file")
	}

	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}

	if ok, err := ctl.CanOperate(cfg); !ok {
		return err
	}

	return nodegroup.New(cfg, ctl, nil).Update(options)
}
