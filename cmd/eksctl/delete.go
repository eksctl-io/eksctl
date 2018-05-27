package main

import (
	"fmt"
	"os"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/spf13/cobra"

	cloudformation "github.com/weaveworks/eksctl/pkg/cfn"
)

func deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "delete",
		Run: func(c *cobra.Command, _ []string) {
			c.Help()
		},
	}

	cmd.AddCommand(deleteClusterCmd())

	return cmd
}

func deleteClusterCmd() *cobra.Command {
	config := &cloudformation.Config{}

	cmd := &cobra.Command{
		Use: "cluster",
		Run: func(_ *cobra.Command, _ []string) {
			if err := doDeleteCluster(config); err != nil {
				logger.Critical(err.Error())
				os.Exit(1)
			}
		},
	}

	fs := cmd.Flags()

	fs.StringVarP(&config.ClusterName, "cluster-name", "n", "", "EKS cluster name (required)")
	fs.StringVarP(&config.Region, "region", "r", DEFAULT_EKS_REGION, "AWS region")

	return cmd
}

func doDeleteCluster(config *cloudformation.Config) error {
	cfn := cloudformation.New(config)

	if err := cfn.CheckAuth(); err != nil {
		return err
	}

	if config.ClusterName == "" {
		return fmt.Errorf("--cluster-name must be set")
	}

	logger.Info("deleting EKS cluster %q", config.ClusterName)

	if err := cfn.DeleteStackServiceRole(); err != nil {
		return err
	}

	if err := cfn.DeleteStackVPC(); err != nil {
		return err
	}

	logger.Success("all EKS cluster %q resource will be deleted (if in doubt, check CloudForamtion console)", config.ClusterName)

	return nil
}
