package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/aws/aws-sdk-go/service/cloudformation"

	"github.com/kubicorn/kubicorn/pkg/logger"

	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/api"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

var (
	utilsKubeconfigInputPath  string
	utilsKubeconfigOutputPath string
	utilsSetContext           bool
	utilsAutoKubeconfigPath   bool
	utilsDescribeStackAll     bool
	utilsDescribeStackEvents  bool
)

func utilsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "utils",
		Short: "Various utils",
		Run: func(c *cobra.Command, _ []string) {
			if err := c.Help(); err != nil {
				logger.Debug("ignoring error %q", err.Error())
			}
		},
	}

	cmd.AddCommand(waitNodesCmd())
	cmd.AddCommand(writeKubeconfigCmd())
	cmd.AddCommand(describeStacksCmd())

	return cmd
}

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
	fs.IntVarP(&cfg.MinNodes, "nodes-min", "m", DefaultNodeCount, "minimum number of nodes to wait for")
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

func writeKubeconfigCmd() *cobra.Command {
	cfg := &api.ClusterConfig{}

	cmd := &cobra.Command{
		Use:   "write-kubeconfig",
		Short: "Write kubeconfig file for a given cluster",
		Run: func(_ *cobra.Command, args []string) {
			if err := doWriteKubeconfigCmd(cfg, getNameArg(args)); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	fs := cmd.Flags()

	fs.StringVarP(&cfg.ClusterName, "name", "n", "", "EKS cluster name (required)")

	fs.StringVarP(&cfg.Region, "region", "r", "", "AWS region")
	fs.StringVarP(&cfg.Profile, "profile", "p", "", "AWS credentials profile to use (overrides the AWS_PROFILE environment variable)")

	fs.BoolVar(&utilsAutoKubeconfigPath, "auto-kubeconfig", false, fmt.Sprintf("save kubconfig file by cluster name – %q", kubeconfig.AutoPath("<name>")))
	fs.StringVar(&utilsKubeconfigOutputPath, "kubeconfig", kubeconfig.DefaultPath, "path to write kubeconfig")
	fs.BoolVar(&utilsSetContext, "set-kubeconfig-context", true, "if true then current-context will be set in kubeconfig; if a context is already set then it will be overwritten")

	return cmd
}

func doWriteKubeconfigCmd(cfg *api.ClusterConfig, name string) error {
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

	if utilsAutoKubeconfigPath {
		if utilsKubeconfigOutputPath != kubeconfig.DefaultPath {
			return fmt.Errorf("--kubeconfig and --auto-kubeconfig cannot be used at the same time")
		}
		utilsKubeconfigOutputPath = kubeconfig.AutoPath(cfg.ClusterName)
	}

	cluster, err := ctl.DescribeControlPlane()
	if err != nil {
		return err
	}

	logger.Debug("cluster = %#v", cluster)

	if err = ctl.GetCredentials(*cluster); err != nil {
		return err
	}

	clientConfigBase, err := ctl.NewClientConfig()
	if err != nil {
		return err
	}

	config := clientConfigBase.WithExecAuthenticator()
	filename, err := kubeconfig.Write(utilsKubeconfigOutputPath, config.Client, utilsSetContext)
	if err != nil {
		return errors.Wrap(err, "writing kubeconfig")
	}

	logger.Success("saved kubeconfig as %q", filename)

	return nil
}

func describeStacksCmd() *cobra.Command {
	cfg := &api.ClusterConfig{}

	cmd := &cobra.Command{
		Use:   "describe-stacks",
		Short: "Describe CloudFormation stack for a given cluster",
		Run: func(_ *cobra.Command, args []string) {
			if err := doDescribeStacksCmd(cfg, getNameArg(args)); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	fs := cmd.Flags()

	fs.StringVarP(&cfg.ClusterName, "name", "n", "", "EKS cluster name (required)")

	fs.StringVarP(&cfg.Region, "region", "r", "", "AWS region")
	fs.StringVarP(&cfg.Profile, "profile", "p", "", "AWS credentials profile to use (overrides the AWS_PROFILE environment variable)")

	fs.BoolVar(&utilsDescribeStackAll, "all", false, "include deleted stacks")
	fs.BoolVar(&utilsDescribeStackEvents, "events", true, "include stack events")

	return cmd
}

func doDescribeStacksCmd(cfg *api.ClusterConfig, name string) error {
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

	stackManager := ctl.NewStackManager()

	stacks, err := stackManager.DescribeStacks(cfg.ClusterName)
	if err != nil {
		return err
	}

	if len(stacks) < 2 {
		logger.Warning("only %d stacks found, for a ready-to-use cluster there should be at least 2", len(stacks))
	}

	for _, s := range stacks {
		if !utilsDescribeStackAll && *s.StackStatus == cloudformation.StackStatusDeleteComplete {
			continue
		}
		logger.Info("stack/%s = %#v", *s.StackName, s)
		if utilsDescribeStackEvents {
			events, err := stackManager.DescribeStackEvents(s)
			if err != nil {
				logger.Critical(err.Error())
			}
			for i, e := range events {
				logger.Info("events/%s[%d] = %#v", *s.StackName, i, e)
			}
		}
	}

	return nil
}
