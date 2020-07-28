package utils

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/addons"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func installSpotOceanController(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("install-spot-ocean-controller", "Install Spot Ocean controller", "")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doInstallSpotOceanController(cmd)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlagWithDeprecated(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cmd)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)
}

func doInstallSpotOceanController(cmd *cmdutils.Cmd) error {
	if err := cmdutils.NewMetadataLoader(cmd).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig

	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}

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

	controller := addons.NewSpotOceanController(
		rawClient,
		cfg,
		cmd.Plan,
		"")

	if err := controller.Deploy(); err != nil {
		return fmt.Errorf("ocean: error installing controller: %w", err)
	}

	cmdutils.LogPlanModeWarning(cmd.Plan)
	return nil
}
