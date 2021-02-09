package utils

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func waitNodesCmd(cmd *cmdutils.Cmd) {
	cmd.CobraCommand.Deprecated = "This command will be removed in version 0.26.0, if you use it please let us know at https://github.com/weaveworks/eksctl/issues/1185."
	cmd.CobraCommand.Hidden = true
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()
	cmd.ClusterConfig = cfg

	var kubeconfigPath string

	cmd.SetDescription("wait-nodes", "Wait for nodes", "")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doWaitNodes(cmd, ng, kubeconfigPath)
	}

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
	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}

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
