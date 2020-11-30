package create

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/fargate"
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
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})
	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)
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

	client := fargate.NewClientWithWaitTimeout(cfg.Metadata.Name, ctl.Provider.EKS(), ctl.Provider.WaitTimeout())
	if err := eks.DoCreateFargateProfiles(cmd.ClusterConfig, client); err != nil {
		return err
	}
	clientSet, err := clientSet(cfg, ctl)
	if err != nil {
		return err
	}
	return eks.ScheduleCoreDNSOnFargateIfRelevant(cmd.ClusterConfig, ctl, clientSet)
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
