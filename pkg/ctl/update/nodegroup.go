package update

import (
	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func updateNodeGroupCmd(cmd *cmdutils.Cmd) {
	cmd.ClusterConfig = api.NewClusterConfig()

	cmd.SetDescription(
		"nodegroup",
		"Update nodegroup",
		dedent.Dedent(`Update nodegroup and its configuration.

		Please consult the eksctl documentation for more info on which config fields can be updated with this command.
		To upgrade a nodegroup, please use 'eksctl upgrade nodegroup' instead.
		Note that this is only available for managed nodegroups. 
	`),
	)

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return updateNodegroup(cmd)
	}
}

func updateNodegroup(cmd *cmdutils.Cmd) error {
	if err := cmdutils.NewUpdateNodegroupLoader(cmd).Load(); err != nil {
		return err
	}

	ctl, err := cmd.NewProviderForExistingCluster()
	if err != nil {
		return err
	}

	if ok, err := ctl.CanOperate(cmd.ClusterConfig); !ok {
		return err
	}

	return nodegroup.New(cmd.ClusterConfig, ctl, nil).Update()
}
