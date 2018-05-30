package main

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/kubicorn/kubicorn/pkg/namer"

	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/utils"
)

func createCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "create",
		Run: func(c *cobra.Command, _ []string) {
			c.Help()
		},
	}

	cmd.AddCommand(createClusterCmd())

	return cmd
}

const (
	DEFAULT_EKS_REGION     = "us-west-2"
	DEFAULT_NODE_COUNT     = 2
	DEFAULT_NODE_TYPE      = "m5.large"
	DEFAULT_SSH_PUBLIC_KEY = "~/.ssh/id_rsa.pub"
)

var (
	writeKubeconfig bool
	kubeconfigPath  string
)

func createClusterCmd() *cobra.Command {
	cfg := &eks.Config{}

	cmd := &cobra.Command{
		Use: "cluster",
		Run: func(_ *cobra.Command, _ []string) {
			if err := doCreateCluster(cfg); err != nil {
				logger.Critical(err.Error())
				os.Exit(1)
			}
		},
	}

	fs := cmd.Flags()

	fs.StringVarP(&cfg.ClusterName, "cluster-name", "n", "", fmt.Sprintf("EKS cluster name (generated, e.g. %q)", namer.RandomName()))
	fs.StringVarP(&cfg.Region, "region", "r", DEFAULT_EKS_REGION, "AWS region")

	fs.StringVarP(&cfg.NodeType, "node-type", "t", DEFAULT_NODE_TYPE, "node instance type")
	fs.IntVarP(&cfg.Nodes, "nodes", "N", DEFAULT_NODE_COUNT, "total number of nodes for a fixed ASG")

	// TODO(p2): review parameter validation, this shouldn't be needed initially
	fs.IntVarP(&cfg.MinNodes, "nodes-min", "m", 0, "maximum nodes in ASG")
	fs.IntVarP(&cfg.MaxNodes, "nodes-max", "M", 0, "minimum nodes in ASG")

	fs.StringVar(&cfg.SSHPublicKeyPath, "ssh-public-key", DEFAULT_SSH_PUBLIC_KEY, "SSH public key to use for nodes (import from local path, or use existing EC2 key pair)")

	fs.BoolVar(&writeKubeconfig, "write-kubeconfig", true, "toggle writing of kubeconfig")
	fs.StringVar(&kubeconfigPath, "kubeconfig", "kubeconfig", "path to write kubeconfig")

	return cmd
}

func doCreateCluster(cfg *eks.Config) error {
	ctl := eks.New(cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.ClusterName == "" {
		cfg.ClusterName = fmt.Sprintf("%s-%d", namer.RandomName(), time.Now().Unix())
	}

	if cfg.SSHPublicKeyPath == "" {
		return fmt.Errorf("--ssh-public-key must be non-empty string")
	}

	if err := ctl.LoadSSHPublicKey(); err != nil {
		return err
	}

	logger.Debug("cfg = %#v", cfg)

	logger.Info("creating EKS cluster %q", cfg.ClusterName)

	{ // core action
		taskErr := make(chan error)
		// create each of the core cloudformation stacks
		ctl.CreateAllStacks(taskErr)
		// read any errors (it only gets non-nil errors)
		for err := range taskErr {
			return err
		}
	}

	logger.Success("all EKS cluster %q resources has been created", cfg.ClusterName)

	// obtain cluster credentials, write kubeconfig

	{ // post-creation action
		clientConfigBase, err := ctl.NewClientConfig()
		if err != nil {
			return err
		}

		// TODO(p2): make kubeconfig writter merge with current default kubeconfig and respect KUBECONFIG env var for writing
		if writeKubeconfig {
			if err := clientConfigBase.WithExecHeptioAuthenticator().WriteToFile(kubeconfigPath); err != nil {
				return errors.Wrap(err, "writing kubeconfig")
			}
			logger.Info("wrote %q", kubeconfigPath)
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

		// TODO(p2): addons

		// check kubectl version, and offer install instructions if missing or old
		// also check heptio-authenticator
		// TODO(p2): and offer install instructions if missing
		// TODO(p2): add sub-command for these checks
		// TODO(p3): few more extensive checks, i.e. some basic validation
		if err := utils.CheckAllCommands(kubeconfigPath); err != nil {
			return err
		}
	}

	return nil
}
