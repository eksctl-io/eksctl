package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/kubicorn/kubicorn/pkg/logger"

	"github.com/weaveworks/eksctl/pkg/eks"
)

func getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get resource(s)",
		Run: func(c *cobra.Command, _ []string) {
			c.Help()
		},
	}

	cmd.AddCommand(getClusterCmd())

	return cmd
}

func getClusterCmd() *cobra.Command {
	cfg := &eks.ClusterConfig{}

	cmd := &cobra.Command{
		Use:     "cluster",
		Short:   "Get custer(s)",
		Aliases: []string{"clusters"},
		Run: func(_ *cobra.Command, args []string) {
			if err := doGetCluster(cfg, getNameArg(args)); err != nil {
				logger.Critical(err.Error())
				os.Exit(1)
			}
		},
	}

	fs := cmd.Flags()

	fs.StringVarP(&cfg.ClusterName, "name", "n", "", "EKS cluster name")

	fs.StringVarP(&cfg.Region, "region", "r", DEFAULT_EKS_REGION, "AWS region")
	fs.StringVarP(&cfg.Profile, "profile", "p", "", "AWS creditials profile to use (overrides the AWS_PROFILE environment variable)")

	return cmd
}

func doGetCluster(cfg *eks.ClusterConfig, name string) error {
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

	if err := ctl.ListClusters(); err != nil {
		return err
	}

	return nil
}
