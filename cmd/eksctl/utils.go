package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/utils"
)

var (
	utilsKubeconfigInputPath  string
	utilsKubeconfigOutputPath string
)

func utilsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "utils",
		Short: "Various utils",
		Run: func(c *cobra.Command, _ []string) {
			c.Help()
		},
	}

	cmd.AddCommand(waitNodesCmd())
	cmd.AddCommand(writeKubeconfigCmd())

	return cmd
}

func waitNodesCmd() *cobra.Command {
	cfg := &eks.ClusterConfig{}

	cmd := &cobra.Command{
		Use:   "wait-nodes",
		Short: "Wait for nodes",
		Run: func(_ *cobra.Command, _ []string) {
			if err := doWaitNodes(cfg); err != nil {
				logger.Critical(err.Error())
				os.Exit(1)
			}
		},
	}

	fs := cmd.Flags()

	fs.StringVar(&utilsKubeconfigInputPath, "kubeconfig", "kubeconfig", "path to read kubeconfig")
	fs.IntVarP(&cfg.MinNodes, "nodes-min", "m", DEFAULT_NODE_COUNT, "minimum nodes to wait for")

	return cmd
}

func doWaitNodes(cfg *eks.ClusterConfig) error {
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

	if err := cfg.WaitForNodes(clientset); err != nil {
		return err
	}

	return nil
}

func writeKubeconfigCmd() *cobra.Command {
	cfg := &eks.ClusterConfig{}

	cmd := &cobra.Command{
		Use:   "write-kubeconfig",
		Short: "Write kubeconfig file for a given cluster",
		Run: func(_ *cobra.Command, _ []string) {
			if err := doWriteKubeconfigCmd(cfg); err != nil {
				logger.Critical(err.Error())
				os.Exit(1)
			}
		},
	}

	fs := cmd.Flags()

	fs.StringVarP(&cfg.ClusterName, "cluster-name", "n", "", fmt.Sprintf("EKS cluster name (generated if unspecified, e.g. %q)", utils.ClusterName()))
	fs.StringVarP(&cfg.Region, "region", "r", DEFAULT_EKS_REGION, "AWS region")

	fs.StringVar(&utilsKubeconfigOutputPath, "kubeconfig", "", "path to write kubeconfig")

	return cmd
}

func doWriteKubeconfigCmd(cfg *eks.ClusterConfig) error {
	ctl := eks.New(cfg)

	if cfg.ClusterName == "" {
		return fmt.Errorf("--cluster-name must be set")
	}

	if utilsKubeconfigOutputPath == "" {
		utilsKubeconfigOutputPath = utils.ConfigPath(cfg.ClusterName)
	}

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	cluster, err := ctl.DescribeControlPlane()
	if err != nil {
		return err
	}

	logger.Debug("cluster = %#v", cluster)

	if err := ctl.GetCredentials(*cluster); err != nil {
		return err
	}

	clientConfigBase, err := ctl.NewClientConfig()
	if err != nil {
		return err
	}

	if err := clientConfigBase.WithExecHeptioAuthenticator().WriteToFile(utilsKubeconfigOutputPath); err != nil {
		return errors.Wrap(err, "writing kubeconfig")
	}

	return nil
}
