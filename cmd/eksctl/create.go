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
	DEFAULT_NODE_TYPE      = "m5.large"          // seems like good value for money
	DEFAULT_SSH_PUBLIC_KEY = "~/.ssh/id_rsa.pub" // TODO kops does this, let's just upload one to make it easy
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
	fs.IntVarP(&config.MinNodes, "nodes-min", "m", 0, "maximum nodes in ASG")
	fs.IntVarP(&config.MaxNodes, "nodes-max", "M", 0, "minimum nodes in ASG")

	// fs.StringVar(&config.TODO, "ssh-public-key", DEFAULT_SSH_PUBLIC_KEY, "SSH public key to use for nodes")

	// TODO:
	// --nodes
	// --kubeconfig <path>
	// --write-kuhbeconfig <booL>
	return cmd
}

func doCreateCluster(config *cloudformation.Config) error {
	cfn := cloudformation.New(config)

	if err := cfn.CheckAuth(); err != nil {
		return err
	}

	logger.Info("creating EKS cluster %q", config.ClusterName)

	// create each of the cloudformation stacks

	// TODO waitgroups
	{
		done := make(chan error)
		if err := cfn.CreateStackServiceRole(done); err != nil {
			return err
		}
		if err := <-done; err != nil {
			return err
		}
	}
	{
		done := make(chan error)
		if err := cfn.CreateStackVPC(done); err != nil {
			return err
		}
		if err := <-done; err != nil {
			return err
		}
	}

	logger.Success("all EKS cluster %q resources has been created", config.ClusterName)

	// obtain cluster credentials

	// login to the cluster and authorise nodes to join

	// watch nodes joining

	// validate (like in kops)

	return nil
}
