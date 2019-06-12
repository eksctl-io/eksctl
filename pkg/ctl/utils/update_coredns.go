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

func updateCoreDNSCmd(rc *cmdutils.ResourceCmd) {
	cfg := api.NewClusterConfig()
	rc.ClusterConfig = cfg

	rc.SetDescription("update-coredns", "Update coredns add-on to ensure image the standard Amazon EKS version", "")

	rc.SetRunFuncWithNameArg(func() error {
		return doUpdateCoreDNS(rc)
	})

	rc.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddNameFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, rc.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &rc.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, rc)
	})

	cmdutils.AddCommonFlagsForAWS(rc.FlagSetGroup, rc.ProviderConfig, false)
}

func doUpdateCoreDNS(rc *cmdutils.ResourceCmd) error {
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

	if err := ctl.GetCredentials(cfg); err != nil {
		return errors.Wrapf(err, "getting credentials for cluster %q", meta.Name)
	}

	clientSet, err := ctl.NewStdClientSet(cfg)
	if err != nil {
		return err
	}

	updateRequired, err := defaultaddons.UpdateCoreDNSImageTag(clientSet, rc.Plan)
	if err != nil {
		return err
	}

	cmdutils.LogPlanModeWarning(rc.Plan && updateRequired)

	return nil
}
