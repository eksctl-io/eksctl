package delete

import (
	"github.com/kris-nova/logger"
	"github.com/spf13/pflag"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/fargate"
)

func deleteFargateProfile(cmd *cmdutils.Cmd) {
	cmd.ClusterConfig = api.NewClusterConfig()
	cmd.SetDescription(
		"fargateprofile",
		"Delete Fargate profile",
		"",
	)
	opts := configureDeleteFargateProfileCmd(cmd)
	cmd.SetRunFuncWithNameArg(func() error {
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
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddWaitFlag(fs, &cmd.Wait, "wait for the deletion of the Fargate profile, which may take from a couple seconds to a couple minutes.")
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})
	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
	return &opts
}

func doDeleteFargateProfile(cmd *cmdutils.Cmd, opts *fargate.Options) error {
	if err := opts.Validate(); err != nil {
		return err
	}
	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}
	if err := ctl.CheckAuth(); err != nil {
		return err
	}
	clusterName := cmd.ClusterConfig.Metadata.Name
	awsClient := fargate.NewClientWithWaitTimeout(clusterName, ctl.Provider.EKS(), cmd.ProviderConfig.WaitTimeout)
	logger.Info("deleting Fargate profile %q", opts.ProfileName)
	return awsClient.DeleteProfile(opts.ProfileName, cmd.Wait)
}
