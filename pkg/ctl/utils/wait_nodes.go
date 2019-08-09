package utils

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func waitNodesCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()
	cmd.ClusterConfig = cfg

	var kubeconfigPath string

	cmd.SetDescription("wait-nodes", "Wait for nodes", "")

	cmd.SetRunFunc(func() error {
		return doWaitNodes(cmd, ng, kubeconfigPath)
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&kubeconfigPath, "kubeconfig", "kubeconfig", "path to read kubeconfig")
		minSize := fs.IntP("nodes-min", "m", api.DefaultNodeCount, "minimum number of nodes to wait for")
		cmdutils.AddPreRun(cmd.CobraCommand, func(cobraCmd *cobra.Command, args []string) {
			if f := cobraCmd.Flag("nodes-min"); f.Changed {
				ng.MinSize = minSize
			}
		})
		fs.DurationVar(&cmd.ProviderConfig.WaitTimeout, "timeout", api.DefaultWaitTimeout, "how long to wait")
	})
}

func doWaitNodes(cmd *cmdutils.Cmd, ng *api.NodeGroup, kubeconfigPath string) error {
	cfg := cmd.ClusterConfig

	ctl := eks.New(cmd.ProviderConfig, cfg)

	if kubeconfigPath == "" {
		return cmdutils.ErrMustBeSet("--kubeconfig")
	}

	clientConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return err
	}

	clientSet, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return err
	}

	if err := ctl.WaitForNodes(clientSet, ng); err != nil {
		return err
	}

	return nil
}
