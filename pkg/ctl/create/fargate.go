package create

import (
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/fargate"
)

func createFargateProfile(cmd *cmdutils.Cmd) {
	cmd.ClusterConfig = api.NewClusterConfig()
	cmd.SetDescription(
		"fargateprofile",
		"Create a Fargate profile",
		"",
	)
	options := configureCreateFargateProfileCmd(cmd)
	cmd.SetRunFuncWithNameArg(func() error {
		return doCreateFargateProfile(cmd, options)
	})
}

func configureCreateFargateProfileCmd(cmd *cmdutils.Cmd) *fargate.CreateOptions {
	var options fargate.CreateOptions
	cmd.FlagSetGroup.InFlagSet("Fargate", func(fs *pflag.FlagSet) {
		cmdutils.AddFlagsForFargateProfileCreation(fs, &options)
	})
	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cmd.ClusterConfig.Metadata)
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})
	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
	return &options
}

func doCreateFargateProfile(cmd *cmdutils.Cmd, options *fargate.CreateOptions) error {
	if err := cmdutils.NewCreateFargateProfileLoader(cmd, options).Load(); err != nil {
		return err
	}
	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}
	if err := ctl.CheckAuth(); err != nil {
		return err
	}
	cfg := cmd.ClusterConfig
	roleARN, err := getClusterRoleARN(ctl, cfg.Metadata)
	if err != nil {
		return err
	}
	return doCreateFargateProfiles(cmd, ctl, roleARN)
}

func getClusterRoleARN(ctl *eks.ClusterProvider, meta *api.ClusterMeta) (string, error) {
	eksCluster, err := ctl.DescribeControlPlane(meta)
	if err != nil {
		return "", errors.Wrapf(err, "failed to retrieve EKS cluster role ARN for \"%v\"", meta.Name)
	}
	roleARN := *eksCluster.RoleArn
	logger.Debug("default Fargate profile pod execution role ARN: %v", roleARN)
	return roleARN, nil
}

func doCreateFargateProfiles(cmd *cmdutils.Cmd, ctl *eks.ClusterProvider, defaultPodExecRoleARN string) error {
	clusterName := cmd.ClusterConfig.Metadata.Name
	awsClient := fargate.NewClient(clusterName, ctl.Provider.EKS())
	for _, profile := range cmd.ClusterConfig.FargateProfiles {
		logger.Info("creating Fargate profile \"%s\" on EKS cluster \"%s\"", profile.Name, clusterName)
		// Default the pod execution role ARN to be the same as the cluster
		// role defined in CloudFormation:
		if profile.PodExecutionRoleARN == "" {
			profile.PodExecutionRoleARN = defaultPodExecRoleARN
		}
		if err := awsClient.CreateProfile(profile); err != nil {
			return errors.Wrapf(err, "failed to create Fargate profile \"%s\" on EKS cluster \"%s\"", profile.Name, clusterName)
		}
		logger.Info("created Fargate profile \"%s\" on EKS cluster \"%s\"", profile.Name, clusterName)
	}
	return nil
}
