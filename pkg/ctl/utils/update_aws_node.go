package utils

import (
	"github.com/spf13/pflag"

	defaultaddons "github.com/weaveworks/eksctl/pkg/addons/default"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func updateAWSNodeCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("update-aws-node", "Update aws-node add-on to latest released version", "")

	cmd.SetRunFuncWithNameArg(func() error {
		return doUpdateAWSNode(cmd)
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

func doUpdateAWSNode(cmd *cmdutils.Cmd) error {
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

	updateRequired, err := defaultaddons.UpdateAWSNode(rawClient, meta.Region, cmd.Plan)
	if err != nil {
		return err
	}

	cmdutils.LogPlanModeWarning(cmd.Plan && updateRequired)

	return nil
}
