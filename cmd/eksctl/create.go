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
		Use:   "create",
		Short: "Create resource(s)",
		Run: func(c *cobra.Command, _ []string) {
			c.Help()
		},
	}

	cmd.AddCommand(createClusterCmd())

	return cmd
}

var (
	writeKubeconfig bool
	kubeconfigPath  string
)

func getClusterName() string {
	return fmt.Sprintf("%s-%d", namer.RandomName(), time.Now().Unix())
}

func createClusterCmd() *cobra.Command {
	cfg := &eks.ClusterConfig{Interactive: true}

	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Create a custer (all flags are optional)",
		Run: func(_ *cobra.Command, _ []string) {
			if err := doCreateCluster(cfg); err != nil {
				logger.Critical(err.Error())
				os.Exit(1)
			}
		},
	}

	fs := cmd.Flags()

	fs.StringVarP(&cfg.ClusterName, "cluster-name", "n", "", fmt.Sprintf("EKS cluster name (generated if unspecified, e.g. %q)", getClusterName()))
	fs.StringVarP(&cfg.Region, "region", "r", eks.DEFAULT_REGION, "AWS region")

	fs.StringVarP(&cfg.NodeType, "node-type", "t", eks.DEFAULT_NODE_TYPE, "node instance type")
	fs.IntVarP(&cfg.Nodes, "nodes", "N", eks.DEFAULT_NODE_COUNT, "total number of nodes (for a static ASG)")
	fs.StringVarP(&cfg.NodeAMI, "node-ami", "", "", "custom node AMI")

	fs.IntVarP(&cfg.MinNodes, "nodes-min", "m", 0, "minimum nodes in ASG")
	fs.IntVarP(&cfg.MaxNodes, "nodes-max", "M", 0, "maximum nodes in ASG")

	fs.StringVar(&cfg.SSHPublicKeyPath, "ssh-public-key", eks.DEFAULT_SSH_PUBLIC_KEY, "SSH public key to use for nodes (import from local path, or use existing EC2 key pair)")

	fs.BoolVar(&writeKubeconfig, "write-kubeconfig", true, "toggle writing of kubeconfig")
	fs.StringVar(&kubeconfigPath, "kubeconfig", "kubeconfig", "path to write kubeconfig")

	return cmd
}

func doCreateCluster(cfg *eks.ClusterConfig) error {
	ctl := eks.New(cfg)

	if cfg.ClusterName == "" {
		cfg.ClusterName = getClusterName()
	}

	if err := ctl.CheckConfig(); err != nil {
		return err
	}

	if err := ctl.CheckNodeCountConfig(); err != nil {
		return err
	}

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if err := ctl.LoadSSHPublicKey(); err != nil {
		return err
	}

	logger.Debug("cfg = %#v", cfg)

	logger.Info("creating EKS cluster %q", cfg.ClusterName)

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
		clientConfigBase, err := ctl.NewClientConfig()
		if err != nil {
			return err
		}

		// TODO: https://github.com/weaveworks/eksctl/issues/29
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

		// check kubectl version, and offer install instructions if missing or old
		// also check heptio-authenticator
		// TODO: https://github.com/weaveworks/eksctl/issues/30
		if err := utils.CheckAllCommands(kubeconfigPath); err != nil {
			logger.Critical(err.Error())
			logger.Info("cluster should be functions despite missing client binaries that need to be installed in the PATH")
		}
	}

	return nil
}
