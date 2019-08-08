package utils

import (
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	defaultaddons "github.com/weaveworks/eksctl/pkg/addons/default"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func updateKubeProxyCmd(rc *cmdutils.ResourceCmd) {
	cfg := api.NewClusterConfig()
	rc.ClusterConfig = cfg

	rc.SetDescription("update-kube-proxy", "Update kube-proxy add-on to ensure image matches Kubernetes control plane version", "")

	rc.SetRunFuncWithNameArg(func() error {
		return doUpdateKubeProxy(rc)
	})

	rc.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddNameFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, rc.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &rc.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, rc)
		cmdutils.AddTimeoutFlag(fs, &rc.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(rc.FlagSetGroup, rc.ProviderConfig, false)
}

func doUpdateKubeProxy(rc *cmdutils.ResourceCmd) error {
	if err := cmdutils.NewMetadataLoader(rc).Load(); err != nil {
		return err
	}

	cfg := rc.ClusterConfig
	meta := rc.ClusterConfig.Metadata

	ctl, err := rc.NewCtl()
	if err != nil {
		return err
	}
	logger.Info("using region %s", meta.Region)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if err := ctl.RefreshClusterConfig(cfg); err != nil {
		return errors.Wrapf(err, "getting credentials for cluster %q", meta.Name)
	}

	rawClient, err := ctl.NewRawClient(cfg)
	if err != nil {
		return err
	}

	kubernetesVersion, err := rawClient.ServerVersion()
	if err != nil {
		return err
	}

	updateRequired, err := defaultaddons.UpdateKubeProxyImageTag(rawClient.ClientSet(), kubernetesVersion, rc.Plan)
	if err != nil {
		return err
	}

	cmdutils.LogPlanModeWarning(rc.Plan && updateRequired)

	return nil
}
