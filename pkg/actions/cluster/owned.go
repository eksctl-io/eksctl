package cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/amazon-ec2-instance-selector/v2/pkg/selector"

	"github.com/kris-nova/logger"

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

func NewOwnedCluster(ctx context.Context, cfg *api.ClusterConfig, ctl *eks.ClusterProvider, clusterStack *manager.Stack, stackManager manager.StackManager) (*OwnedCluster, error) {
	instanceSelector, err := selector.New(context.Background(), ctl.AWSProvider.AWSConfig())
	if err != nil {
		return nil, err
	}
	return &OwnedCluster{
		cfg:          cfg,
		ctl:          ctl,
		clusterStack: clusterStack,
		stackManager: stackManager,
		newClientSet: func() (kubernetes.Interface, error) {
			return ctl.NewStdClientSet(cfg)
		},
		newNodeGroupManager: func(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, clientSet kubernetes.Interface) NodeGroupDrainer {
			return nodegroup.New(cfg, ctl, clientSet, instanceSelector)
		},
	}, nil
}

func (c *OwnedCluster) Upgrade(ctx context.Context, dryRun bool) error {
	if err := vpc.UseFromClusterStack(ctx, c.ctl.AWSProvider, c.clusterStack, c.cfg); err != nil {
		return fmt.Errorf("getting VPC configuration for cluster %q: %w", c.cfg.Metadata.Name, err)
	}

	versionUpdateRequired, err := upgrade(ctx, c.cfg, c.ctl, dryRun)
	if err != nil {
		return err
	}

	stackUpdateRequired, err := c.stackManager.AppendNewClusterStackResource(ctx, false, dryRun)
	if err != nil {
		return err
	}

	if err := eks.ValidateExistingNodeGroupsForCompatibility(ctx, c.cfg, c.stackManager); err != nil {
		logger.Critical("failed checking nodegroups", err.Error())
	}

	cmdutils.LogPlanModeWarning(dryRun && (stackUpdateRequired || versionUpdateRequired))
	return nil
}

func (c *OwnedCluster) Delete(ctx context.Context, _, podEvictionWaitPeriod time.Duration, wait, force, disableNodegroupEviction bool, parallel int) error {
	clusterOperable, err := c.ctl.CanOperate(c.cfg)
	if err != nil {
		logger.Debug("failed to check if cluster is operable: %v", err)
	}

	// moving this here was fine because inside `NewTasksToDeleteClusterWithNodeGroups` we did it anyway.
	allStacks, err := c.stackManager.ListNodeGroupStacksWithStatuses(ctx)
	if err != nil {
		return err
	}

	var clientSet kubernetes.Interface
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

		nodeGroupManager := c.newNodeGroupManager(c.cfg, c.ctl, clientSet)
		if err := drainAllNodeGroups(ctx, c.cfg, c.ctl, clientSet, allStacks, disableNodegroupEviction, parallel, nodeGroupManager, func(clusterConfig *api.ClusterConfig, ctl *eks.ClusterProvider, clientSet kubernetes.Interface) {
			attemptVpcCniDeletion(ctx, clusterConfig, ctl, clientSet)
		}, podEvictionWaitPeriod); err != nil {
			if !force {
				return err
			}

			logger.Warning("an error occurred during nodegroups draining, force=true so proceeding with deletion: %q", err.Error())
		}
	}

	if err := deleteSharedResources(ctx, c.cfg, c.ctl, c.stackManager, clusterOperable, clientSet); err != nil {
		if err != nil {
			if force {
				logger.Warning("error occurred during deletion: %v", err)
			} else {
				return err
			}
		}
	}

	newOIDCManager := func() (*iamoidc.OpenIDConnectManager, error) {
		return c.ctl.NewOpenIDConnectManager(ctx, c.cfg)
	}
	tasks, err := c.stackManager.NewTasksToDeleteClusterWithNodeGroups(ctx, c.clusterStack, allStacks, clusterOperable, newOIDCManager, c.ctl.Status.ClusterInfo.Cluster, kubernetes.NewCachedClientSet(clientSet), wait, force, func(errs chan error, _ string) error {
		logger.Info("trying to cleanup dangling network interfaces")
		stack, err := c.stackManager.DescribeClusterStack(ctx)
		if err != nil {
			return fmt.Errorf("error describing cluster stack: %w", err)
		}
		if err := c.ctl.LoadClusterVPC(ctx, c.cfg, stack); err != nil {
			return fmt.Errorf("getting VPC configuration for cluster %q: %w", c.cfg.Metadata.Name, err)
		}

		go func() {
			errs <- vpc.CleanupNetworkInterfaces(ctx, c.ctl.AWSProvider.EC2(), c.cfg)
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

	if err := c.deleteKarpenterStackIfExists(ctx); err != nil {
		return err
	}

	if err := checkForUndeletedStacks(ctx, c.stackManager); err != nil {
		return err
	}

	logger.Success("all cluster resources were deleted")

	return nil
}

func (c *OwnedCluster) deleteKarpenterStackIfExists(ctx context.Context) error {
	stack, err := c.stackManager.GetKarpenterStack(ctx)
	if err != nil {
		return err
	}

	if stack != nil {
		logger.Info("deleting karpenter stack")
		return c.stackManager.DeleteStackSync(ctx, stack)
	}

	return nil
}
