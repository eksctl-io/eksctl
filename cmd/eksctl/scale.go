package main

import (
	"fmt"
	"os"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/api"
)

func scaleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scale",
		Short: "Scale resources(s)",
		Run: func(c *cobra.Command, _ []string) {
			c.Help()
		},
	}

	cmd.AddCommand(scaleNodeGroupCmd())

	return cmd
}

func scaleNodeGroupCmd() *cobra.Command {
	cfg := &api.ClusterConfig{}

	cmd := &cobra.Command{
		Use:   "nodegroup",
		Short: "Scale a nodegroup",
		Run: func(_ *cobra.Command, args []string) {
			if err := doScaleNodeGroup(cfg); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	fs := cmd.Flags()

	fs.StringVarP(&cfg.ClusterName, "name", "n", "", "EKS cluster name")
	fs.IntVarP(&cfg.Nodes, "nodes", "N", 0, "total number of nodes (scale to this number)")

	fs.StringVarP(&cfg.Region, "region", "r", api.DEFAULT_EKS_REGION, "AWS region")
	fs.StringVarP(&cfg.Profile, "profile", "p", "", "AWS creditials profile to use (overrides the AWS_PROFILE environment variable)")

	fs.DurationVar(&cfg.WaitTimeout, "timeout", api.DefaultWaitTimeout, "max wait time in any polling operations")

	return cmd
}

func doScaleNodeGroup(cfg *api.ClusterConfig) error {
	ctl := eks.New(cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.ClusterName == "" {
		return fmt.Errorf("no cluster name supplied. Use the --name= flag")
	}

	if cfg.Nodes < 1 {
		return fmt.Errorf("number of nodes must be greater than 0. Use the --nodes/-N flag")
	}

	stackManager := ctl.NewStackManager()
	errs := stackManager.ScaleInitialNodeGroup()
	if len(errs) > 0 {
		logger.Info("%d error(s) occurred and nodegroup can't be scaled", len(errs))
		for _, err := range errs {
			logger.Critical("%s\n", err.Error())
		}
		return fmt.Errorf("failed to scale nodegroup for cluster %q", cfg.ClusterName)
	}

	return nil
}
