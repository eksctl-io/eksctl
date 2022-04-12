package utils

import (
	"context"
	"time"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func installWindowsVPCController(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("install-vpc-controllers", "Install Windows VPC controller to support running Windows workloads", "")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doInstallWindowsVPCController(cmd)
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

func doInstallWindowsVPCController(cmd *cmdutils.Cmd) error {
	if err := cmdutils.NewMetadataLoader(cmd).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig
	meta := cmd.ClusterConfig.Metadata

	ctl, err := cmd.NewProviderForExistingCluster()
	if err != nil {
		return err
	}
	logger.Info("using region %s", meta.Region)

	if ok, err := ctl.CanUpdate(cfg); !ok {
		return err
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(cmd.ProviderConfig.WaitTimeout))
	defer cancel()

	vpcControllerTask := &eks.VPCControllerTask{
		Info:            "install Windows VPC controller",
		Context:         ctx,
		ClusterConfig:   cfg,
		ClusterProvider: ctl,
		PlanMode:        cmd.Plan,
	}

	taskTree := &tasks.TaskTree{
		Tasks: []tasks.Task{vpcControllerTask},
	}

	if errs := taskTree.DoAllSync(); len(errs) > 0 {
		return errs[0]
	}

	cmdutils.LogPlanModeWarning(cmd.Plan)
	return nil
}
