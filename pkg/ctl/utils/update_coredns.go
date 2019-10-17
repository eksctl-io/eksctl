package utils

import (
	"github.com/spf13/pflag"

	defaultaddons "github.com/weaveworks/eksctl/pkg/addons/default"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func updateCoreDNSCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("update-coredns", "Update coredns add-on to ensure image matches the standard Amazon EKS version", "")

	cmd.SetRunFuncWithNameArg(func() error {
		return doUpdateCoreDNS(cmd)
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlagWithDeprecated(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cmd)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
}

func doUpdateCoreDNS(cmd *cmdutils.Cmd) error {
	if err := cmdutils.NewMetadataLoader(cmd).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig
	meta := cmd.ClusterConfig.Metadata

	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}
	cmdutils.LogRegionAndVersionInfo(meta)

	if err := ctl.CheckAuth(); err != nil {
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

	updateRequired, err := defaultaddons.UpdateCoreDNS(rawClient, meta.Region, kubernetesVersion, cmd.Plan)
	if err != nil {
		return err
	}

	cmdutils.LogPlanModeWarning(cmd.Plan && updateRequired)

	return nil
}
