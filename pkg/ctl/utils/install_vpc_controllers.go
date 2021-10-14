package utils

import (
	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/addons"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func installWindowsVPCController(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("install-vpc-controllers", "Install Windows VPC controller to support running Windows workloads", "")
	cmd.CobraCommand.Deprecated = addons.VPCControllerInfoMessage

	var deleteController bool

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doInstallWindowsVPCController(cmd, deleteController)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.BoolVar(&deleteController, "delete", false, "Deletes VPC resource controller from worker nodes")
		cmdutils.AddClusterFlagWithDeprecated(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cmd)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)
}

func doInstallWindowsVPCController(cmd *cmdutils.Cmd, deleteController bool) error {
	if !deleteController {
		logger.Warning(addons.VPCControllerInfoMessage)
		return nil
	}

	if err := cmdutils.NewMetadataLoader(cmd).Load(); err != nil {
		return err
	}

	ctl, err := cmd.NewProviderForExistingCluster()
	if err != nil {
		return err
	}
	logger.Info("using region %s", cmd.ClusterConfig.Metadata.Region)
	rawClient, err := ctl.NewRawClient(cmd.ClusterConfig)
	if err != nil {
		return err
	}
	deleteControllerTask := &eks.DeleteVPCControllerTask{
		Info:      "delete Windows VPC controller",
		PlanMode:  cmd.Plan,
		RawClient: rawClient,
	}

	taskTree := &tasks.TaskTree{
		Tasks: []tasks.Task{deleteControllerTask},
	}

	if errs := taskTree.DoAllSync(); len(errs) > 0 {
		return errs[0]
	}

	cmdutils.LogPlanModeWarning(cmd.Plan)
	return nil
}
