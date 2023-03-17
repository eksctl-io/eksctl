package delete

import (
	"context"
	"fmt"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/fargate"
)

func deleteFargateProfileWithRunFunc(cmd *cmdutils.Cmd, runFunc func(cmd *cmdutils.Cmd, opts *fargate.Options) error) {
	cmd.ClusterConfig = api.NewClusterConfig()
	cmd.SetDescription(
		"fargateprofile",
		"Delete Fargate profile",
		"",
	)
	opts := configureDeleteFargateProfileCmd(cmd)
	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		if err := cmdutils.NewDeleteFargateProfileLoader(cmd, opts).Load(); err != nil {
			return err
		}
		return runFunc(cmd, opts)
	}
}

func deleteFargateProfile(cmd *cmdutils.Cmd) {
	deleteFargateProfileWithRunFunc(cmd, func(cmd *cmdutils.Cmd, opts *fargate.Options) error {
		return doDeleteFargateProfile(cmd, opts)
	})
}

func configureDeleteFargateProfileCmd(cmd *cmdutils.Cmd) *fargate.Options {
	var opts fargate.Options
	cmd.FlagSetGroup.InFlagSet("Fargate", func(fs *pflag.FlagSet) {
		cmdutils.AddFlagsForFargate(fs, &opts)
	})
	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cmd.ClusterConfig.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddWaitFlag(fs, &cmd.Wait, "wait for the deletion of the Fargate profile, which may take from a couple seconds to a couple minutes.")
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})
	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, false)
	return &opts
}

func doDeleteFargateProfile(cmd *cmdutils.Cmd, opts *fargate.Options) error {
	ctx := context.Background()
	ctl, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}

	clusterName := cmd.ClusterConfig.Metadata.Name
	manager := fargate.NewFromProvider(clusterName, ctl.AWSProvider, ctl.NewStackManager(cmd.ClusterConfig))
	if cmd.Wait {
		logger.Info(deletingFargateProfileMsg(clusterName, opts.ProfileName))
	} else {
		logger.Debug(deletingFargateProfileMsg(clusterName, opts.ProfileName))
	}
	if err := manager.DeleteProfile(ctx, opts.ProfileName, cmd.Wait); err != nil {
		return err
	}
	logger.Info("deleted Fargate profile %q on EKS cluster %q", opts.ProfileName, clusterName)
	return nil
}

func deletingFargateProfileMsg(clusterName, profileName string) string {
	return fmt.Sprintf("deleting Fargate profile %q on EKS cluster %q", profileName, clusterName)
}
