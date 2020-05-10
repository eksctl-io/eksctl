package utils

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/addons"
	client "github.com/weaveworks/eksctl/pkg/kubernetes"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func installCloudwatchAgentCmd(cmd *cmdutils.Cmd) {
	installCloudwatchAgentWithRunFunc(cmd, func(cmd *cmdutils.Cmd) error {
		return doInstallCloudwatchAgent(cmd)
	})
}

func installCloudwatchAgentWithRunFunc(cmd *cmdutils.Cmd, runFunc func(cmd *cmdutils.Cmd) error) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("install-cloudwatch-agent", "Install cloudwatch agent with default prometheus setting", "")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		return runFunc(cmd)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cmd)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
}

func doInstallCloudwatchAgent(cmd *cmdutils.Cmd) error {
	if err := cmdutils.NewMetadataLoader(cmd).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig
	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet(cmdutils.ClusterNameFlag(cmd))
	}

	clientSet, restConfig, err := kubernetesClientAndConfigFrom(cmd)
	if err != nil {
		return err
	}

	c, err := client.NewRawClient(clientSet, restConfig)
	if err != nil {
		return err
	}

	cwAgent := addons.NewCloudwatchAgent(c, cfg.Metadata.Name, cfg.Metadata.Region, cmd.Plan)
	if err := cwAgent.Deploy(); err != nil {
		return errors.Wrap(err, "error installing cloudwatch agent")
	}
	cmdutils.LogPlanModeWarning(cmd.Plan)
	return nil
}

func kubernetesClientAndConfigFrom(cmd *cmdutils.Cmd) (*kubernetes.Clientset, *rest.Config, error) {
	ctl, err := cmd.NewCtl()
	if err != nil {
		return nil, nil, err
	}
	if err := ctl.CheckAuth(); err != nil {
		return nil, nil, err
	}
	cfg := cmd.ClusterConfig
	if ok, err := ctl.CanOperate(cfg); !ok {
		return nil, nil, err
	}
	kubernetesClientConfigs, err := ctl.NewClient(cfg)
	if err != nil {
		return nil, nil, err
	}
	k8sConfig := kubernetesClientConfigs.Config
	k8sRestConfig, err := clientcmd.NewDefaultClientConfig(*k8sConfig, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, nil, err
	}
	k8sClientSet, err := kubernetes.NewForConfig(k8sRestConfig)
	if err != nil {
		return nil, nil, err
	}
	return k8sClientSet, k8sRestConfig, nil
}
