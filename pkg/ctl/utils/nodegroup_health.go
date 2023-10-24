package utils

import (
	"context"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/managed"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func nodeGroupHealthCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	var nodeGroupName string

	cmd.SetDescription("nodegroup-health", "Get nodegroup health for a managed node", "")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return getNodeGroupHealth(cmd, nodeGroupName)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cfg.Metadata)
		fs.StringVarP(&nodeGroupName, "name", "n", "", "Name of the nodegroup")

		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, false)
}

func getNodeGroupHealth(cmd *cmdutils.Cmd, nodeGroupName string) error {
	cfg := cmd.ClusterConfig

	ctx := context.TODO()
	ctl, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}
	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet(cmdutils.ClusterNameFlag(cmd))
	}

	if nodeGroupName != "" && cmd.NameArg != "" {
		return cmdutils.ErrFlagAndArg("--name", nodeGroupName, cmd.NameArg)
	}

	if cmd.NameArg != "" {
		nodeGroupName = cmd.NameArg
	}

	if nodeGroupName == "" {
		return cmdutils.ErrMustBeSet("name")
	}

	if cfg.IsControlPlaneOnOutposts() {
		return api.ErrUnsupportedLocalCluster
	}
	stackCollection := manager.NewStackCollection(ctl.AWSProvider, cfg)
	managedService := managed.NewService(ctl.AWSProvider.EKS(), ctl.AWSProvider.EC2(), stackCollection, cfg.Metadata.Name)
	healthIssues, err := managedService.GetHealth(ctx, nodeGroupName)
	if err != nil {
		return err
	}

	if len(healthIssues) == 0 {
		logger.Info("No health issues found. Node group %q is active", nodeGroupName)
		return nil
	}

	for _, issue := range healthIssues {
		logger.Warning(issue.Message)
	}

	return nil
}
