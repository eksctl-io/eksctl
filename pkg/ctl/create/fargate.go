package create

import (
	"fmt"
	"time"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/fargate"
	"github.com/weaveworks/eksctl/pkg/fargate/coredns"
	"github.com/weaveworks/eksctl/pkg/utils"
	"github.com/weaveworks/eksctl/pkg/utils/retry"
	"github.com/weaveworks/eksctl/pkg/utils/strings"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func createFargateProfileWithRunFunc(cmd *cmdutils.Cmd, runFunc func(cmd *cmdutils.Cmd) error) {
	cmd.ClusterConfig = api.NewClusterConfig()
	cmd.SetDescription(
		"fargateprofile",
		"Create a Fargate profile",
		"",
	)
	options := configureCreateFargateProfileCmd(cmd)
	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		if err := cmdutils.NewCreateFargateProfileLoader(cmd, options).Load(); err != nil {
			return err
		}
		return runFunc(cmd)
	}
}

func createFargateProfile(cmd *cmdutils.Cmd) {
	createFargateProfileWithRunFunc(cmd, doCreateFargateProfile)
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

func doCreateFargateProfile(cmd *cmdutils.Cmd) error {
	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}
	if err := ctl.CheckAuth(); err != nil {
		return err
	}
	cfg := cmd.ClusterConfig
	if ok, err := ctl.CanOperate(cfg); !ok {
		return err
	}

	supportsFargate, err := ctl.SupportsFargate(cfg)
	if err != nil {
		return err
	}
	if !supportsFargate {
		return fmt.Errorf("Fargate is not supported for this cluster version. Please update the cluster to be at least eks.%d", fargate.MinPlatformVersion)
	}

	if err := ctl.LoadClusterVPC(cfg); err != nil {
		return err
	}

	// Read back the default Fargate pod execution role ARN from CloudFormation:
	if err := ctl.NewStackManager(cfg).RefreshFargatePodExecutionRoleARN(); err != nil {
		return err
	}

	if err := doCreateFargateProfiles(cmd, ctl); err != nil {
		return err
	}
	clientSet, err := clientSet(cfg, ctl)
	if err != nil {
		return err
	}
	return scheduleCoreDNSOnFargateIfRelevant(cmd, clientSet)
}

func clientSet(cfg *api.ClusterConfig, ctl *eks.ClusterProvider) (kubernetes.Interface, error) {
	kubernetesClientConfigs, err := ctl.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	k8sConfig := kubernetesClientConfigs.Config
	k8sRestConfig, err := clientcmd.NewDefaultClientConfig(*k8sConfig, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, err
	}
	k8sClientSet, err := kubernetes.NewForConfig(k8sRestConfig)
	if err != nil {
		return nil, err
	}
	return k8sClientSet, nil
}

func doCreateFargateProfiles(cmd *cmdutils.Cmd, ctl *eks.ClusterProvider) error {
	clusterName := cmd.ClusterConfig.Metadata.Name
	awsClient := fargate.NewClientWithWaitTimeout(clusterName, ctl.Provider.EKS(), cmd.ProviderConfig.WaitTimeout)
	for _, profile := range cmd.ClusterConfig.FargateProfiles {
		logger.Info("creating Fargate profile %q on EKS cluster %q", profile.Name, clusterName)

		// Default the pod execution role ARN to be the same as the cluster
		// role defined in CloudFormation:
		if profile.PodExecutionRoleARN == "" {
			profile.PodExecutionRoleARN = strings.EmptyIfNil(cmd.ClusterConfig.IAM.FargatePodExecutionRoleARN)
		}
		// Linearise the creation of Fargate profiles by passing
		// wait = true, as the API otherwise errors out with:
		//   ResourceInUseException: Cannot create Fargate Profile
		//   ${name2} because cluster ${clusterName} currently has
		//   Fargate profile ${name1} in status CREATING
		if err := awsClient.CreateProfile(profile, true); err != nil {
			return errors.Wrapf(err, "failed to create Fargate profile %q on EKS cluster %q", profile.Name, clusterName)
		}
		logger.Info("created Fargate profile %q on EKS cluster %q", profile.Name, clusterName)
	}
	return nil
}

func scheduleCoreDNSOnFargateIfRelevant(cmd *cmdutils.Cmd, clientSet kubernetes.Interface) error {
	if coredns.IsSchedulableOnFargate(cmd.ClusterConfig.FargateProfiles) {
		betaAPIDeprecated, err := utils.IsMinVersion(api.Version1_16, cmd.ClusterConfig.Metadata.Version)
		if err != nil {
			return err
		}
		useBetaAPIGroup := !betaAPIDeprecated
		scheduled, err := coredns.IsScheduledOnFargate(clientSet, useBetaAPIGroup)
		if err != nil {
			return err
		}
		if !scheduled {
			if err := coredns.ScheduleOnFargate(clientSet, useBetaAPIGroup); err != nil {
				return err
			}
			retryPolicy := &retry.TimingOutExponentialBackoff{
				Timeout:  cmd.ProviderConfig.WaitTimeout,
				TimeUnit: time.Second,
			}
			if err := coredns.WaitForScheduleOnFargate(clientSet, retryPolicy, useBetaAPIGroup); err != nil {
				return err
			}
		}
	}
	return nil
}
