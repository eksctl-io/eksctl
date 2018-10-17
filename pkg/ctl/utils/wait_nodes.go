package utils

import (
	"fmt"
	"os"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/api"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func waitNodesCmd() *cobra.Command {
	cfg := &api.ClusterConfig{}

	cmd := &cobra.Command{
		Use:   "wait-nodes",
		Short: "Wait for nodes",
		Run: func(_ *cobra.Command, _ []string) {
			if err := doWaitNodes(cfg); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	fs := cmd.Flags()

	fs.StringVar(&utilsKubeconfigInputPath, "kubeconfig", "kubeconfig", "path to read kubeconfig")
	fs.IntVarP(&cfg.MinNodes, "nodes-min", "m", api.DefaultNodeCount, "minimum number of nodes to wait for")
	fs.DurationVar(&cfg.WaitTimeout, "timeout", api.DefaultWaitTimeout, "how long to wait")

	return cmd
}

func doWaitNodes(cfg *api.ClusterConfig) error {
	ctl := eks.New(cfg)

	if utilsKubeconfigInputPath == "" {
		return fmt.Errorf("--kubeconfig must be set")
	}

	clientConfig, err := clientcmd.BuildConfigFromFlags("", utilsKubeconfigInputPath)
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return err
	}

	if err := ctl.WaitForNodes(clientset); err != nil {
		return err
	}

	return nil
}
