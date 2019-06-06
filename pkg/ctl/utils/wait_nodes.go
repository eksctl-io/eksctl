package utils

import (
	"os"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var waitNodesKubeconfigPath string

func waitNodesCmd(g *cmdutils.Grouping) *cobra.Command {
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()
	cp := cmdutils.NewCommonParams(cfg)

	cp.Command = &cobra.Command{
		Use:   "wait-nodes",
		Short: "Wait for nodes",
		Run: func(_ *cobra.Command, _ []string) {
			if err := doWaitNodes(cp, ng); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := g.New(cp.Command)

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&waitNodesKubeconfigPath, "kubeconfig", "kubeconfig", "path to read kubeconfig")
		minSize := fs.IntP("nodes-min", "m", api.DefaultNodeCount, "minimum number of nodes to wait for")
		cmdutils.AddPreRun(cp.Command, func(cmd *cobra.Command, args []string) {
			if f := cmd.Flag("nodes-min"); f.Changed {
				ng.MinSize = minSize
			}
		})
		fs.DurationVar(&cp.ProviderConfig.WaitTimeout, "timeout", api.DefaultWaitTimeout, "how long to wait")
	})

	group.AddTo(cp.Command)
	return cp.Command
}

func doWaitNodes(cp *cmdutils.CommonParams, ng *api.NodeGroup) error {
	cfg := cp.ClusterConfig

	ctl := eks.New(cp.ProviderConfig, cfg)

	if waitNodesKubeconfigPath == "" {
		return cmdutils.ErrMustBeSet("--kubeconfig")
	}

	clientConfig, err := clientcmd.BuildConfigFromFlags("", waitNodesKubeconfigPath)
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
