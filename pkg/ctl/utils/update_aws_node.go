package utils

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	defaultaddons "github.com/weaveworks/eksctl/pkg/addons/default"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

func updateAWSNodeCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("update-aws-node", "Update aws-node add-on to latest released version", "")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doUpdateAWSNode(cmd)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlagWithDeprecated(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cmd)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, false)
}

func doUpdateAWSNode(cmd *cmdutils.Cmd) error {
	ctx := context.TODO()
	return updateAddon(ctx, cmd, api.VPCCNIAddon, func(rawClient *kubernetes.RawClient, _ defaultaddons.AddonVersionDescriber) (bool, error) {
		return defaultaddons.UpdateAWSNode(ctx, defaultaddons.AddonInput{
			RawClient: rawClient,
			Region:    cmd.ClusterConfig.Metadata.Region,
		}, cmd.Plan)
	})
}
