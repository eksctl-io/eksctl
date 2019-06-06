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

func waitNodesCmd(rc *cmdutils.ResourceCmd) {
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()
	rc.ClusterConfig = cfg

	var kubeconfigPath string

	rc.SetDescription("wait-nodes", "Wait for nodes", "")

	rc.SetRunFunc(func() error {
		return doWaitNodes(rc, ng, kubeconfigPath)
	})

	rc.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&kubeconfigPath, "kubeconfig", "kubeconfig", "path to read kubeconfig")
		minSize := fs.IntP("nodes-min", "m", api.DefaultNodeCount, "minimum number of nodes to wait for")
		cmdutils.AddPreRun(rc.Command, func(cmd *cobra.Command, args []string) {
			if f := cmd.Flag("nodes-min"); f.Changed {
				ng.MinSize = minSize
			}
		})
		fs.DurationVar(&rc.ProviderConfig.WaitTimeout, "timeout", api.DefaultWaitTimeout, "how long to wait")
	})
}

func doWaitNodes(rc *cmdutils.ResourceCmd, ng *api.NodeGroup, kubeconfigPath string) error {
	cfg := rc.ClusterConfig

	ctl := eks.New(rc.ProviderConfig, cfg)

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
