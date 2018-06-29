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
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

var (
	utilsKubeconfigInputPath  string
	utilsKubeconfigOutputPath string
	utilsSetContext           bool
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
		Run: func(_ *cobra.Command, args []string) {
			if err := doWriteKubeconfigCmd(cfg, getNameArg(args)); err != nil {
				logger.Critical(err.Error())
				os.Exit(1)
			}
		},
	}

	fs := cmd.Flags()

	fs.StringVarP(&cfg.ClusterName, "name", "n", "", "EKS cluster name (required)")

	fs.StringVarP(&cfg.Region, "region", "r", DEFAULT_EKS_REGION, "AWS region")
	fs.StringVarP(&cfg.Profile, "profile", "p", "", "AWS profile to use. If provided, this overrides the AWS_PROFILE environment variable")

	fs.StringVar(&utilsKubeconfigOutputPath, "kubeconfig", "", "path to write kubeconfig")
	fs.BoolVar(&utilsSetContext, "set-kubeconfig-context", true, "if true then current-context will be set in kubeconfig; if a context is already set then it will be overwritten")

	return cmd
}

func doWriteKubeconfigCmd(cfg *eks.ClusterConfig, name string) error {
	ctl := eks.New(cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.ClusterName != "" && name != "" {
		return fmt.Errorf("--name=%s and argument %s cannot be used at the same time", cfg.ClusterName, name)
	}

	if name != "" {
		cfg.ClusterName = name
	}

	if cfg.ClusterName == "" {
		return fmt.Errorf("--name must be set")
	}

	if utilsKubeconfigOutputPath == "" {
		utilsKubeconfigOutputPath = kubeconfig.AutoPath(cfg.ClusterName)
	}

	cluster, err := ctl.DescribeControlPlane()
	if err != nil {
		return err
	}

	logger.Debug("cluster = %#v", cluster)

	if err := ctl.GetCredentials(*cluster); err != nil {
		return err
	}

	clientConfigBase, err := ctl.NewClientConfig(setContext)
	if err != nil {
		return err
	}

	config := clientConfigBase.WithExecHeptioAuthenticator()
	if err := kubeconfig.WriteKubeCfg(utilsKubeconfigOutputPath, config.Client, setContext); err != nil {
		return errors.Wrap(err, "writing kubeconfig")
	}

	logger.Success("wrote kubeconfig file %q", utilsKubeconfigOutputPath)

	return nil
}
