package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/kubicorn/kubicorn/pkg/logger"

	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

func deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete resource(s)",
		Run: func(c *cobra.Command, _ []string) {
			c.Help()
		},
	}

	cmd.AddCommand(deleteClusterCmd())

	return cmd
}

func deleteClusterCmd() *cobra.Command {
	cfg := &eks.ClusterConfig{}

	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Delete a cluster",
		Run: func(_ *cobra.Command, args []string) {
			if err := doDeleteCluster(cfg, getNameArg(args)); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	fs := cmd.Flags()

	fs.StringVarP(&cfg.ClusterName, "name", "n", "", "EKS cluster name (required)")

	fs.StringVarP(&cfg.Region, "region", "r", DEFAULT_EKS_REGION, "AWS region")
	fs.StringVarP(&cfg.Profile, "profile", "p", "", "AWS creditials profile to use (overrides the AWS_PROFILE environment variable)")

	return cmd
}

func doDeleteCluster(cfg *eks.ClusterConfig, name string) error {
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

	logger.Info("deleting EKS cluster %q", cfg.ClusterName)

	debugError := func(err error) {
		logger.Debug("continue despite error: %v", err)
	}

	if err := ctl.DeleteControlPlane(); err != nil {
		debugError(err)
	}

	if err := ctl.DeleteStackServiceRole(); err != nil {
		debugError(err)
	}

	if err := ctl.DeleteStackVPC(); err != nil {
		debugError(err)
	}

	if err := ctl.DeleteStackDefaultNodeGroup(); err != nil {
		debugError(err)
	}

	ctl.MaybeDeletePublicSSHKey()

	kubeconfig.MaybeDeleteConfig(cfg.ClusterName)

	logger.Success("all EKS cluster %q resource will be deleted (if in doubt, check CloudFormation console)", cfg.ClusterName)

	return nil
}
