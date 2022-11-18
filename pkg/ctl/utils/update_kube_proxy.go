package utils

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	defaultaddons "github.com/weaveworks/eksctl/pkg/addons/default"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
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
	if err := cmdutils.NewMetadataLoader(cmd).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig
	meta := cmd.ClusterConfig.Metadata

	ctx := context.TODO()
	ctl, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}

	if ok, err := ctl.CanUpdate(cfg); !ok {
		return err
	}

	rawClient, err := ctl.NewRawClient(cfg)
	if err != nil {
		return err
	}

	kubernetesVersion, err := rawClient.ServerVersion()
	if err != nil {
		return err
	}

	updateRequired, err := defaultaddons.UpdateKubeProxy(ctx, defaultaddons.AddonInput{
		RawClient:           rawClient,
		ControlPlaneVersion: kubernetesVersion,
		Region:              meta.Region,
		EKSAPI:              ctl.AWSProvider.EKS(),
	}, cmd.Plan)
	if err != nil {
		return err
	}

	cmdutils.LogPlanModeWarning(cmd.Plan && updateRequired)

	return nil
}
