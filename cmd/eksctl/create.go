package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/kubicorn/kubicorn/pkg/logger"

	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/utils"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconf"
)

func createCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create resource(s)",
		Run: func(c *cobra.Command, _ []string) {
			c.Help()
		},
	}

	cmd.AddCommand(createClusterCmd())

	return cmd
}

const (
	EKS_REGION_US_WEST_2   = "us-west-2"
	EKS_REGION_US_EAST_1   = "us-east-1"
	DEFAULT_EKS_REGION     = EKS_REGION_US_WEST_2
	DEFAULT_NODE_COUNT     = 2
	DEFAULT_NODE_TYPE      = "m5.large"
	DEFAULT_SSH_PUBLIC_KEY = "~/.ssh/id_rsa.pub"

	DEFAULT_KUBECONFIG_PATH = "kubeconfig"
)

var (
	writeKubeconfig    bool
	kubeconfigPath     string
	autoKubeconfigPath bool
	setContext         bool
)

func createClusterCmd() *cobra.Command {
	cfg := &eks.ClusterConfig{}

	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Create a custer",
		Run: func(_ *cobra.Command, _ []string) {
			if err := doCreateCluster(cfg); err != nil {
				logger.Critical(err.Error())
				os.Exit(1)
			}
		},
	}

	fs := cmd.Flags()

	exampleClusterName := utils.ClusterName()

	fs.StringVarP(&cfg.ClusterName, "cluster-name", "n", "", fmt.Sprintf("EKS cluster name (generated if unspecified, e.g. %q)", exampleClusterName))
	fs.StringVarP(&cfg.Region, "region", "r", DEFAULT_EKS_REGION, "AWS region")
	fs.StringVarP(&cfg.Profile, "profile", "p", "", "AWS profile to use. If provided, this overrides the AWS_PROFILE environment variable")

	fs.StringVarP(&cfg.NodeType, "node-type", "t", DEFAULT_NODE_TYPE, "node instance type")
	fs.IntVarP(&cfg.Nodes, "nodes", "N", DEFAULT_NODE_COUNT, "total number of nodes (for a static ASG)")

	// TODO: https://github.com/weaveworks/eksctl/issues/28
	fs.IntVarP(&cfg.MinNodes, "nodes-min", "m", 0, "minimum nodes in ASG")
	fs.IntVarP(&cfg.MaxNodes, "nodes-max", "M", 0, "maximum nodes in ASG")

	fs.StringVar(&cfg.SSHPublicKeyPath, "ssh-public-key", DEFAULT_SSH_PUBLIC_KEY, "SSH public key to use for nodes (import from local path, or use existing EC2 key pair)")

	fs.BoolVar(&writeKubeconfig, "write-kubeconfig", true, "toggle writing of kubeconfig")
	fs.BoolVar(&autoKubeconfigPath, "auto-kubeconfig", false, fmt.Sprintf("save kubconfig file by cluster name, e.g. %q", utils.ConfigPath(exampleClusterName)))
	fs.StringVar(&kubeconfigPath, "kubeconfig", DEFAULT_KUBECONFIG_PATH, "path to write kubeconfig (incompatible with --auto-kubeconfig)")
	fs.BoolVar(&setContext, "set-context", true, "If true then current-context will be set in kubeconfig. If a context is already set then it will be overwritten.")

	return cmd
}

func doCreateCluster(cfg *eks.ClusterConfig) error {
	ctl := eks.New(cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.ClusterName == "" {
		cfg.ClusterName = utils.ClusterName()
	}

	if autoKubeconfigPath {
		if kubeconfigPath != DEFAULT_KUBECONFIG_PATH {
			return fmt.Errorf("--kubeconfig and --auto-kubeconfig cannot be used at the same time")
		}
		kubeconfigPath = utils.ConfigPath(cfg.ClusterName)
	}
	if kubeconfigPath == DEFAULT_KUBECONFIG_PATH {
		kubeconfigPath = kubeconf.GetRecommendedPath()
		logger.Debug("no kubeconfig specified so using recommended client path: %s", kubeconfigPath)
	}

	if cfg.SSHPublicKeyPath == "" {
		return fmt.Errorf("--ssh-public-key must be non-empty string")
	}

	if cfg.Region != EKS_REGION_US_WEST_2 && cfg.Region != EKS_REGION_US_EAST_1 {
		return fmt.Errorf("--region=%s is not supported only %s and %s are supported", cfg.Region, EKS_REGION_US_WEST_2, EKS_REGION_US_EAST_1)
	}

	if err := ctl.LoadSSHPublicKey(); err != nil {
		return err
	}

	logger.Debug("cfg = %#v", cfg)

	logger.Info("creating EKS cluster %q in %q region", cfg.ClusterName, cfg.Region)

	{ // core action
		taskErr := make(chan error)
		// create each of the core cloudformation stacks
		ctl.CreateCluster(taskErr)
		// read any errors (it only gets non-nil errors)
		for err := range taskErr {
			return err
		}
	}

	logger.Success("all EKS cluster %q resources has been created", cfg.ClusterName)

	// obtain cluster credentials, write kubeconfig

	{ // post-creation action
		clientConfigBase, err := ctl.NewClientConfig(setContext)
		if err != nil {
			return err
		}

		if writeKubeconfig {
			config := clientConfigBase.WithExecHeptioAuthenticator()
			if err := kubeconf.WriteToFile(kubeconfigPath, config.Client, setContext); err != nil {
				return errors.Wrap(err, "writing kubeconfig")
			}

			logger.Success("wrote %q", kubeconfigPath)
		} else {
			kubeconfigPath = ""
		}

		// create Kubernetes client
		clientSet, err := clientConfigBase.NewClientSetWithEmbeddedToken()
		if err != nil {
			return err
		}

		// authorise nodes to join
		if err := cfg.CreateDefaultNodeGroupAuthConfigMap(clientSet); err != nil {
			return err
		}

		// wait for nodes to join
		if err := cfg.WaitForNodes(clientSet); err != nil {
			return err
		}

		// check kubectl version, and offer install instructions if missing or old
		// also check heptio-authenticator
		// TODO: https://github.com/weaveworks/eksctl/issues/30
		if err := utils.CheckAllCommands(kubeconfigPath); err != nil {
			logger.Critical(err.Error())
			logger.Info("cluster should be functional despite missing client binaries that need to be installed in the PATH")
		}
	}

	logger.Info("EKS cluster %q in %q region is ready", cfg.ClusterName, cfg.Region)

	return nil
}
