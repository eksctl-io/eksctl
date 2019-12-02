package get

import (
	"fmt"
	"os"

	"github.com/kris-nova/logger"
	"github.com/spf13/pflag"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/fargate"
)

type options struct {
	fargate.Options
	getCmdParams
}

func getFargateProfile(cmd *cmdutils.Cmd) {
	cmd.ClusterConfig = api.NewClusterConfig()
	cmd.SetDescription(
		"fargateprofile",
		"Get Fargate profile(s)",
		"",
	)
	options := configureGetFargateProfileCmd(cmd)
	cmd.SetRunFuncWithNameArg(func() error {
		return doGetFargateProfile(cmd, options)
	})
}

func configureGetFargateProfileCmd(cmd *cmdutils.Cmd) *options {
	var options options
	cmd.FlagSetGroup.InFlagSet("Fargate", func(fs *pflag.FlagSet) {
		cmdutils.AddFlagsForFargate(fs, &options.Options)
	})
	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cmd.ClusterConfig.Metadata)
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
		cmdutils.AddCommonFlagsForGetCmd(fs, &options.chunkSize, &options.output)
	})
	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
	return &options
}

func doGetFargateProfile(cmd *cmdutils.Cmd, options *options) error {
	if err := cmdutils.NewMetadataLoader(cmd).Load(); err != nil {
		cmd.CobraCommand.Help()
		return err
	}

	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}
	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	supportsFargate, err := ctl.SupportsFargate(cmd.ClusterConfig)
	if err != nil {
		return err
	}
	if !supportsFargate {
		return fmt.Errorf("Fargate is not supported for this cluster version. Please update the cluster to be at least eks.%d", fargate.MinPlatformVersion)
	}

	clusterName := cmd.ClusterConfig.Metadata.Name
	awsClient := fargate.NewClient(clusterName, ctl.Provider.EKS())

	logger.Debug("getting EKS cluster %q's Fargate profile(s)", clusterName)
	profiles, err := getProfiles(awsClient, options.ProfileName)
	if err != nil {
		return err
	}
	return fargate.PrintProfiles(profiles, os.Stdout, options.output)
}

func getProfiles(awsClient *fargate.Client, name string) ([]*api.FargateProfile, error) {
	if name == "" {
		return awsClient.ReadProfiles()
	}
	profile, err := awsClient.ReadProfile(name)
	if err != nil {
		return nil, err
	}
	return []*api.FargateProfile{profile}, nil
}
