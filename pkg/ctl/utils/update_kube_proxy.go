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

func updateKubeProxyCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("update-kube-proxy", "Update kube-proxy add-on to ensure image matches Kubernetes control plane version", "")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doUpdateKubeProxy(cmd)
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

func doUpdateKubeProxy(cmd *cmdutils.Cmd) error {
	ctx := context.TODO()
	return updateAddon(ctx, cmd, api.KubeProxyAddon, func(rawClient *kubernetes.RawClient, addonDescriber defaultaddons.AddonVersionDescriber) (bool, error) {
		kubernetesVersion, err := rawClient.ServerVersion()
		if err != nil {
			return false, err
		}
		return defaultaddons.UpdateKubeProxy(ctx, defaultaddons.AddonInput{
			RawClient:             rawClient,
			ControlPlaneVersion:   kubernetesVersion,
			Region:                cmd.ClusterConfig.Metadata.Region,
			AddonVersionDescriber: addonDescriber,
		}, cmd.Plan)
	})
}
