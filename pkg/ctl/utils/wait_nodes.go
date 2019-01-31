package utils

import (
	"fmt"
	"os"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var waitNodesKubeconfigPath string

func waitNodesCmd(g *cmdutils.Grouping) *cobra.Command {
	p := &api.ProviderConfig{}
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()

	cmd := &cobra.Command{
		Use:   "wait-nodes",
		Short: "Wait for nodes",
		Run: func(_ *cobra.Command, _ []string) {
			if err := doWaitNodes(p, cfg, ng); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := g.New(cmd)

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&waitNodesKubeconfigPath, "kubeconfig", "kubeconfig", "path to read kubeconfig")
		minSize := fs.IntP("nodes-min", "m", api.DefaultNodeCount, "minimum number of nodes to wait for")
		cmd.PreRun = func(cmd *cobra.Command, args []string) {
			if f := cmd.Flag("nodes-min"); f.Changed {
				ng.MinSize = minSize
			}
		}
		fs.DurationVar(&p.WaitTimeout, "timeout", api.DefaultWaitTimeout, "how long to wait")
	})

	group.AddTo(cmd)

	return cmd
}

func doWaitNodes(p *api.ProviderConfig, cfg *api.ClusterConfig, ng *api.NodeGroup) error {
	ctl := eks.New(p, cfg)

	if waitNodesKubeconfigPath == "" {
		return fmt.Errorf("--kubeconfig must be set")
	}

	clientConfig, err := clientcmd.BuildConfigFromFlags("", waitNodesKubeconfigPath)
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return err
	}

	if err := ctl.WaitForNodes(clientset, ng); err != nil {
		return err
	}

	return nil
}
