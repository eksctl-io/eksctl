package utils

import (
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	defaultaddons "github.com/weaveworks/eksctl/pkg/addons/default"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

func updateAWSNodeCmd(rc *cmdutils.ResourceCmd) {
	cfg := api.NewClusterConfig()
	rc.ClusterConfig = cfg

	rc.SetDescription("update-aws-node", "Update aws-node add-on to latest released version", "")

	rc.SetRunFuncWithNameArg(func() error {
		return doUpdateAWSNode(rc)
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

func doUpdateAWSNode(rc *cmdutils.ResourceCmd) error {
	if err := cmdutils.NewMetadataLoader(rc).Load(); err != nil {
		return err
	}

	cfg := rc.ClusterConfig
	meta := rc.ClusterConfig.Metadata

	ctl := eks.New(rc.ProviderConfig, cfg)

	if !ctl.IsSupportedRegion() {
		return cmdutils.ErrUnsupportedRegion(rc.ProviderConfig)
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

	updateRequired, err := defaultaddons.UpdateAWSNode(rawClient, meta.Region, kubernetesVersion, rc.Plan)
	if err != nil {
		return err
	}

	cmdutils.LogPlanModeWarning(rc.Plan && updateRequired)

	return nil
}
