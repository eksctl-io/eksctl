package cluster

import (
	"time"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

type OwnedCluster struct {
	cfg                 *api.ClusterConfig
	ctl                 *eks.ClusterProvider
	clusterStack        *manager.Stack
	stackManager        manager.StackManager
	newClientSet        func() (kubernetes.Interface, error)
	newNodeGroupManager func(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, clientSet kubernetes.Interface) NodeGroupDrainer
}

func NewOwnedCluster(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, clusterStack *manager.Stack, stackManager manager.StackManager) *OwnedCluster {
	return &OwnedCluster{
		cfg:          cfg,
		ctl:          ctl,
		clusterStack: clusterStack,
		stackManager: stackManager,
		newClientSet: func() (kubernetes.Interface, error) {
			return ctl.NewStdClientSet(cfg)
		},
		newNodeGroupManager: func(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, clientSet kubernetes.Interface) NodeGroupDrainer {
			return nodegroup.New(cfg, ctl, clientSet)
		},
	}
}

func (c *OwnedCluster) Upgrade(dryRun bool) error {
	if err := vpc.UseFromClusterStack(c.ctl.Provider, c.clusterStack, c.cfg); err != nil {
		return errors.Wrapf(err, "getting VPC configuration for cluster %q", c.cfg.Metadata.Name)
	}

	versionUpdateRequired, err := upgrade(c.cfg, c.ctl, dryRun)
	if err != nil {
		return err
	}

	if err := c.ctl.RefreshClusterStatus(c.cfg); err != nil {
		return err
	}

	supportsManagedNodes, err := c.ctl.SupportsManagedNodes(c.cfg)
	if err != nil {
		return err
	}

	stackUpdateRequired, err := c.stackManager.AppendNewClusterStackResource(dryRun, supportsManagedNodes)
	if err != nil {
		return err
	}

	nodeGroupService := eks.NodeGroupService{Provider: c.ctl.Provider}
	if err := nodeGroupService.ValidateExistingNodeGroupsForCompatibility(c.cfg, c.stackManager); err != nil {
		logger.Critical("failed checking nodegroups", err.Error())
	}

	cmdutils.LogPlanModeWarning(dryRun && (stackUpdateRequired || versionUpdateRequired))
	return nil
}

func (c *OwnedCluster) Delete(_ time.Duration, wait, force, disableNodegroupEviction bool) error {
	var (
		clientSet kubernetes.Interface
		oidc      *iamoidc.OpenIDConnectManager
	)

	clusterOperable, err := c.ctl.CanOperate(c.cfg)
	if err != nil {
		logger.Debug("failed to check if cluster is operable: %v", err)
	}

	oidcSupported := true
	if clusterOperable {
		var err error
		clientSet, err = c.newClientSet()
		if err != nil {
			if force {
				logger.Warning("error occurred during deletion: %v", err)
			} else {
				return err
			}
		}

		oidc, err = c.ctl.NewOpenIDConnectManager(c.cfg)
		if err != nil {
			if _, ok := err.(*eks.UnsupportedOIDCError); !ok {
				if force {
					logger.Warning("error occurred during deletion: %v", err)
				} else {
					return err
				}
			}
			oidcSupported = false
		}
		allStacks, err := c.stackManager.ListNodeGroupStacks()
		if err != nil {
			return err
		}

		nodeGroupManager := c.newNodeGroupManager(c.cfg, c.ctl, clientSet)
		if err := drainAllNodeGroups(c.cfg, c.ctl, clientSet, allStacks, disableNodegroupEviction, nodeGroupManager, attemptVpcCniDeletion); err != nil {
			if !force {
				return err
			}

			logger.Warning("an error occurred during nodegroups draining, force=true so proceeding with deletion: %q", err.Error())
		}
	}

	if err := deleteSharedResources(c.cfg, c.ctl, c.stackManager, clusterOperable, clientSet); err != nil {
		if err != nil {
			if force {
				logger.Warning("error occurred during deletion: %v", err)
			} else {
				return err
			}
		}
	}

	deleteOIDCProvider := clusterOperable && oidcSupported
	tasks, err := c.stackManager.NewTasksToDeleteClusterWithNodeGroups(deleteOIDCProvider, oidc, kubernetes.NewCachedClientSet(clientSet), wait, func(errs chan error, _ string) error {
		logger.Info("trying to cleanup dangling network interfaces")
		if err := c.ctl.LoadClusterVPC(c.cfg, c.stackManager); err != nil {
			return errors.Wrapf(err, "getting VPC configuration for cluster %q", c.cfg.Metadata.Name)
		}

		go func() {
			errs <- vpc.CleanupNetworkInterfaces(c.ctl.Provider.EC2(), c.cfg)
			close(errs)
		}()
		return nil
	})

	if err != nil {
		return err
	}

	if tasks.Len() == 0 {
		logger.Warning("no cluster resources were found for %q", c.cfg.Metadata.Name)
		return nil
	}

	logger.Info(tasks.Describe())
	if errs := tasks.DoAllSync(); len(errs) > 0 {
		return handleErrors(errs, "cluster with nodegroup(s)")
	}

	if err := c.deleteKarpenterStackIfExists(); err != nil {
		return err
	}

	if err := checkForUndeletedStacks(c.stackManager); err != nil {
		return err
	}

	logger.Success("all cluster resources were deleted")

	return nil
}

func (c *OwnedCluster) deleteKarpenterStackIfExists() error {
	stack, err := c.stackManager.GetKarpenterStack()
	if err != nil {
		return err
	}

	if stack != nil {
		logger.Info("deleting karpenter stack")
		return c.stackManager.DeleteStackByNameSync(*stack.StackName)
	}

	return nil
}
