package main

import (
	"fmt"
	"os"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/kubicorn/kubicorn/pkg/namer"
	"github.com/spf13/cobra"

	cloudformation "github.com/weaveworks/eksctl/pkg/cfn"
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

func createClusterCmd() *cobra.Command {
	config := &cloudformation.Config{}

	cmd := &cobra.Command{
		Use: "cluster",
		Run: func(_ *cobra.Command, _ []string) {
			if err := doCreateCluster(config); err != nil {
				logger.Critical(err.Error())
				os.Exit(1)
			}
		},
	}

	fs := cmd.Flags()

	fs.StringVarP(&config.ClusterName, "cluster-name", "n", "", fmt.Sprintf("EKS cluster name (generated, e.g. %q)", namer.RandomName()))
	fs.StringVarP(&config.Region, "region", "r", DEFAULT_EKS_REGION, "AWS region")

	fs.StringVarP(&config.NodeType, "node-type", "t", DEFAULT_NODE_TYPE, "node instance type")
	fs.IntVarP(&config.Nodes, "nodes", "N", DEFAULT_NODE_COUNT, "total number of nodes for a fixed ASG")

	// TODO(p2): review parameter validation, this shouldn't be needed initially
	// fs.IntVarP(&config.MinNodes, "nodes-min", "m", 0, "maximum nodes in ASG")
	// fs.IntVarP(&config.MaxNodes, "nodes-max", "M", 0, "minimum nodes in ASG")

	// TODO(p1): upload SSH key
	// fs.StringVar(&config.TODO, "ssh-public-key", DEFAULT_SSH_PUBLIC_KEY, "SSH public key to use for nodes")

	// TODO(p0):
	// --kubeconfig <path>
	// --write-kuhbeconfig <booL>
	return cmd
}

func doCreateCluster(config *cloudformation.Config) error {
	cfn := cloudformation.New(config)

	if err := cfn.CheckAuth(); err != nil {
		return err
	}

	if config.ClusterName == "" {
		config.ClusterName = namer.RandomName()
	}

	logger.Info("creating EKS cluster %q", config.ClusterName)

	taskErr := make(chan error)
	// create each of the core cloudformation stacks
	cfn.CreateCoreStacks(taskErr)
	// read any errors (it only gets non-nil errors)
	for err := range taskErr {
		return err
	}

	logger.Success("all EKS cluster %q resources has been created", config.ClusterName)

	// TODO(p0): obtain cluster credentials, write kubeconfig

	// TODO(p0): login to the cluster and authorise nodes to join

	// TODO(p1): watch nodes joining

	// TODO(p2): validate (like in kops)

	// TODO(p2): addons

	// TODO(p0): check kubectl version, and offer install instructions if missing or old
	// TODO(p0): check heptio-authenticator, and offer install instructions if missing

	return nil
}
