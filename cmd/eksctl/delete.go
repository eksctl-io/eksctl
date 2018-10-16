package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kubicorn/kubicorn/pkg/logger"

	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/api"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

var (
	waitDelete bool
)

func deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete resource(s)",
		Run: func(c *cobra.Command, _ []string) {
			if err := c.Help(); err != nil {
				logger.Debug("ignoring error %q", err.Error())
			}
		},
	}

	cmd.AddCommand(deleteClusterCmd())

	return cmd
}

func deleteClusterCmd() *cobra.Command {
	cfg := &api.ClusterConfig{}

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

	fs.StringVarP(&cfg.Region, "region", "r", "", "AWS region")
	fs.StringVarP(&cfg.Profile, "profile", "p", "", "AWS credentials profile to use (overrides the AWS_PROFILE environment variable)")

	fs.BoolVarP(&waitDelete, "wait", "w", false, "Wait for deletion of all resources before exiting")

	fs.DurationVar(&cfg.WaitTimeout, "timeout", api.DefaultWaitTimeout, "max wait time in any polling operations")

	return cmd
}

func doDeleteCluster(cfg *api.ClusterConfig, name string) error {
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

	var deletedResources []string

	handleIfError := func(err error, name string) bool {
		if err != nil {
			logger.Debug("continue despite error: %v", err)
			return true
		}
		logger.Debug("deleted %q", name)
		deletedResources = append(deletedResources, name)
		return false
	}

	// We can remove all 'DeprecatedDelete*' calls in 0.2.0

	stackManager := ctl.NewStackManager()

	handleIfError(stackManager.WaitDeleteNodeGroup(), "node group")

	var clusterErr bool
	if waitDelete {
		clusterErr = handleIfError(stackManager.WaitDeleteCluster(), "cluster")
	} else {
		clusterErr = handleIfError(stackManager.DeleteCluster(), "cluster")
	}

	if clusterErr {
		if handleIfError(ctl.DeprecatedDeleteControlPlane(), "control plane") {
			handleIfError(stackManager.DeprecatedDeleteStackControlPlane(waitDelete), "stack control plane")
		}
	}

	handleIfError(stackManager.DeprecatedDeleteStackServiceRole(waitDelete), "node group")
	handleIfError(stackManager.DeprecatedDeleteStackVPC(waitDelete), "stack VPC")
	handleIfError(stackManager.DeprecatedDeleteStackDefaultNodeGroup(waitDelete), "default node group")

	ctl.MaybeDeletePublicSSHKey()

	kubeconfig.MaybeDeleteConfig(cfg)

	if len(deletedResources) == 0 {
		logger.Warning("no EKS cluster resources were found for %q", ctl.Spec.ClusterName)
	} else {
		logger.Success("the following EKS cluster resource(s) for %q will be deleted: %s. If in doubt, check CloudFormation console", ctl.Spec.ClusterName, strings.Join(deletedResources, ", "))
	}

	return nil
}
