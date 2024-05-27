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

func updateCoreDNSCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("update-coredns", "Update coredns add-on to ensure image matches the standard Amazon EKS version", "")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doUpdateCoreDNS(cmd)
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

func doUpdateCoreDNS(cmd *cmdutils.Cmd) error {
	ctx := context.TODO()
	return updateAddon(ctx, cmd, api.CoreDNSAddon, func(rawClient *kubernetes.RawClient, _ defaultaddons.AddonVersionDescriber) (bool, error) {
		kubernetesVersion, err := rawClient.ServerVersion()
		if err != nil {
			return false, err
		}
		return defaultaddons.UpdateCoreDNS(ctx, defaultaddons.AddonInput{
			RawClient:           rawClient,
			ControlPlaneVersion: kubernetesVersion,
			Region:              cmd.ClusterConfig.Metadata.Region,
		}, cmd.Plan)
	})
}
