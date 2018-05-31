package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/kubicorn/kubicorn/pkg/logger"

	"github.com/weaveworks/eksctl/pkg/eks"
)

func getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "get",
		Run: func(c *cobra.Command, _ []string) {
			c.Help()
		},
	}

	cmd.AddCommand(getClusterCmd())

	return cmd
}

func getClusterCmd() *cobra.Command {
	cfg := &eks.Config{}

	cmd := &cobra.Command{
		Use:     "cluster",
		Aliases: []string{"clusters"},
		Run: func(_ *cobra.Command, _ []string) {
			if err := doGetCluster(cfg); err != nil {
				logger.Critical(err.Error())
				os.Exit(1)
			}
		},
	}

	fs := cmd.Flags()

	fs.StringVarP(&cfg.ClusterName, "cluster-name", "n", "", "EKS cluster name")
	fs.StringVarP(&cfg.Region, "region", "r", DEFAULT_EKS_REGION, "AWS region")

	return cmd
}

func doGetCluster(cfg *eks.Config) error {
	ctl := eks.New(cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if err := ctl.ListClusters(); err != nil {
		return err
	}

	return nil
}
