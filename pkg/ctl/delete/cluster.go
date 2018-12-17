package delete

import (
	"fmt"
	"os"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/api"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

func deleteClusterCmd() *cobra.Command {
	p := &api.ProviderConfig{}
	cfg := api.NewClusterConfig()

	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Delete a cluster",
		Run: func(_ *cobra.Command, args []string) {
			if err := doDeleteCluster(p, cfg, cmdutils.GetNameArg(args)); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := &cmdutils.NamedFlagSetGroup{}

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVarP(&cfg.Metadata.Name, "name", "n", "", "EKS cluster name (required)")
		cmdutils.AddRegionFlag(fs, p)
		fs.BoolVarP(&waitDelete, "wait", "w", false, "Wait for deletion of all resources before exiting")
	})

	cmdutils.AddCommonFlagsForAWS(group, p)

	group.AddTo(cmd)
	return cmd
}

func doDeleteCluster(p *api.ProviderConfig, cfg *api.ClusterConfig, nameArg string) error {
	ctl := eks.New(p, cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.Metadata.Name != "" && nameArg != "" {
		return cmdutils.ErrNameFlagAndArg(cfg.Metadata.Name, nameArg)
	}

	if nameArg != "" {
		cfg.Metadata.Name = nameArg
	}

	if cfg.Metadata.Name == "" {
		return fmt.Errorf("--name must be set")
	}

	logger.Info("deleting EKS cluster %q", cfg.Metadata.Name)

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

	stackManager := ctl.NewStackManager(cfg)

	{
		errs := stackManager.WaitDeleteAllNodeGroups()
		if len(errs) > 0 {
			logger.Info("%d error(s) occurred while deleting nodegroup(s)", len(errs))
			for _, err := range errs {
				logger.Critical("%s\n", err.Error())
			}
			handleIfError(fmt.Errorf("failed to delete nodegroup(s)"), "nodegroup(s)")
		}
		logger.Debug("all nodegroups were deleted")
	}

	var clusterErr bool
	if waitDelete {
		clusterErr = handleIfError(stackManager.WaitDeleteCluster(), "cluster")
	} else {
		clusterErr = handleIfError(stackManager.DeleteCluster(), "cluster")
	}

	if clusterErr {
		if handleIfError(ctl.DeprecatedDeleteControlPlane(cfg.Metadata), "control plane") {
			handleIfError(stackManager.DeprecatedDeleteStackControlPlane(waitDelete), "stack control plane (deprecated)")
		}
	}

	handleIfError(stackManager.DeprecatedDeleteStackServiceRole(waitDelete), "service group (deprecated)")
	handleIfError(stackManager.DeprecatedDeleteStackVPC(waitDelete), "stack VPC (deprecated)")
	handleIfError(stackManager.DeprecatedDeleteStackDefaultNodeGroup(waitDelete), "default nodegroup (deprecated)")

	ctl.MaybeDeletePublicSSHKey(cfg.Metadata.Name)

	kubeconfig.MaybeDeleteConfig(cfg.Metadata)

	if len(deletedResources) == 0 {
		logger.Warning("no EKS cluster resources were found for %q", cfg.Metadata.Name)
	} else {
		logger.Success("the following EKS cluster resource(s) for %q will be deleted: %s. If in doubt, check CloudFormation console", cfg.Metadata.Name, strings.Join(deletedResources, ", "))
	}

	return nil
}
