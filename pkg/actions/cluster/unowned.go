package cluster

import (
	"context"
	"time"

	"github.com/aws/amazon-ec2-instance-selector/v2/pkg/selector"

	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	"github.com/weaveworks/eksctl/pkg/actions/addon"
	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

type UnownedCluster struct {
	cfg                 *api.ClusterConfig
	ctl                 *eks.ClusterProvider
	stackManager        manager.StackManager
	newClientSet        func() (kubernetes.Interface, error)
	newNodeGroupManager func(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, clientSet kubernetes.Interface) NodeGroupDrainer
}

func NewUnownedCluster(ctx context.Context, cfg *api.ClusterConfig, ctl *eks.ClusterProvider, stackManager manager.StackManager) (*UnownedCluster, error) {
	instanceSelector, err := selector.New(context.Background(), ctl.AWSProvider.AWSConfig())
	if err != nil {
		return nil, err
	}
	return &UnownedCluster{
		cfg:          cfg,
		ctl:          ctl,
		stackManager: stackManager,
		newClientSet: func() (kubernetes.Interface, error) {
			return ctl.NewStdClientSet(cfg)
		},
		newNodeGroupManager: func(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, clientSet kubernetes.Interface) NodeGroupDrainer {
			return nodegroup.New(cfg, ctl, clientSet, instanceSelector)
		},
	}, nil
}

func (c *UnownedCluster) Upgrade(ctx context.Context, dryRun bool) error {
	versionUpdateRequired, err := upgrade(ctx, c.cfg, c.ctl, dryRun)
	if err != nil {
		return err
	}

	// if no version update is required, don't log asking them to rerun with --approve
	cmdutils.LogPlanModeWarning(dryRun && versionUpdateRequired)
	return nil
}

func (c *UnownedCluster) Delete(ctx context.Context, waitInterval, podEvictionWaitPeriod time.Duration, wait, force, disableNodegroupEviction bool, parallel int) error {
	clusterName := c.cfg.Metadata.Name

	if err := c.checkClusterExists(ctx, clusterName); err != nil {
		return err
	}

	clusterOperable, err := c.ctl.CanOperate(c.cfg)
	if err != nil {
		logger.Debug("failed to check if cluster is operable: %v", err)
	}

	allStacks, err := c.stackManager.ListNodeGroupStacksWithStatuses(ctx)
	if err != nil {
		return err
	}

	var clientSet kubernetes.Interface
	if clusterOperable {
		clientSet, err = c.newClientSet()
		if err != nil {
			return err
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

	if err := c.deleteFargateRoleIfExists(ctx); err != nil {
		return err
	}

	// we have to wait for nodegroups to delete before deleting the cluster
	// so the `wait` value is ignored here
	if err := c.deleteAndWaitForNodegroupsDeletion(ctx, waitInterval, allStacks); err != nil {
		return err
	}

	if err := c.deleteIAMAndOIDC(ctx, wait, clusterOperable, clientSet, force); err != nil {
		if err != nil {
			if force {
				logger.Warning("error occurred during deletion: %v", err)
			} else {
				return err
			}
		}
	}

	if err := c.deleteCluster(ctx, wait); err != nil {
		return err
	}

	if err := checkForUndeletedStacks(ctx, c.stackManager); err != nil {
		return err
	}

	logger.Success("all cluster resources were deleted")
	return nil
}

func (c *UnownedCluster) deleteFargateRoleIfExists(ctx context.Context) error {
	stack, err := c.stackManager.GetFargateStack(ctx)
	if err != nil {
		return err
	}

	if stack != nil {
		logger.Info("deleting fargate role")
		_, err = c.stackManager.DeleteStackBySpec(ctx, stack)
		return err
	}

	logger.Debug("no fargate role found")
	return nil
}

func (c *UnownedCluster) checkClusterExists(ctx context.Context, clusterName string) error {
	_, err := c.ctl.AWSProvider.EKS().DescribeCluster(ctx, &awseks.DescribeClusterInput{
		Name: &c.cfg.Metadata.Name,
	})
	if err != nil {
		if isNotFound(err) {
			return errors.Errorf("cluster %q not found", clusterName)
		}
		return errors.Wrapf(err, "error describing cluster %q", clusterName)
	}
	return nil
}

func (c *UnownedCluster) deleteIAMAndOIDC(ctx context.Context, wait bool, clusterOperable bool, clientSet kubernetes.Interface, force bool) error {
	tasksTree := &tasks.TaskTree{Parallel: false}

	if clusterOperable {
		clientSetGetter := kubernetes.NewCachedClientSet(clientSet)
		newOIDCManager := func() (*iamoidc.OpenIDConnectManager, error) {
			return c.ctl.NewOpenIDConnectManager(ctx, c.cfg)
		}
		serviceAccountAndOIDCTasks, err := c.stackManager.NewTasksToDeleteOIDCProviderWithIAMServiceAccounts(ctx, newOIDCManager, c.ctl.Status.ClusterInfo.Cluster, clientSetGetter, force)
		if err != nil {
			return err
		}

		if serviceAccountAndOIDCTasks.Len() > 0 {
			serviceAccountAndOIDCTasks.IsSubTask = true
			tasksTree.Append(serviceAccountAndOIDCTasks)
		}
	}

	deleteAddonIAMTasks, err := addon.NewRemover(c.stackManager).DeleteAddonIAMTasks(ctx, wait)
	if err != nil {
		return err
	}

	if deleteAddonIAMTasks.Len() > 0 {
		deleteAddonIAMTasks.IsSubTask = true
		tasksTree.Append(deleteAddonIAMTasks)
	}

	if tasksTree.Len() == 0 {
		logger.Warning("no IAM and OIDC resources were found for %q", c.cfg.Metadata.Name)
		return nil
	}

	logger.Info(tasksTree.Describe())
	if errs := tasksTree.DoAllSync(); len(errs) > 0 {
		return handleErrors(errs, "cluster IAM and OIDC")
	}

	logger.Info("all IAM and OIDC resources were deleted")
	return nil
}

func (c *UnownedCluster) deleteCluster(ctx context.Context, wait bool) error {
	clusterName := c.cfg.Metadata.Name

	out, err := c.ctl.AWSProvider.EKS().DeleteCluster(ctx, &awseks.DeleteClusterInput{
		Name: &clusterName,
	})

	if err != nil {
		return err
	}

	logger.Info("initiated deletion of cluster %q", clusterName)
	if out != nil {
		logger.Debug("delete cluster response: %+v", out.Cluster)
	}

	if !wait {
		logger.Info("to see the status of the deletion run `eksctl get cluster --name %s --region %s`", clusterName, c.cfg.Metadata.Region)
		return nil
	}

	logger.Info("waiting for cluster %q to be deleted", clusterName)
	waiter := awseks.NewClusterDeletedWaiter(c.ctl.AWSProvider.EKS())
	return waiter.Wait(ctx, &awseks.DescribeClusterInput{
		Name: &clusterName,
	}, c.ctl.AWSProvider.WaitTimeout())
}

func (c *UnownedCluster) deleteAndWaitForNodegroupsDeletion(ctx context.Context, waitInterval time.Duration, allStacks []manager.NodeGroupStack) error {
	clusterName := c.cfg.Metadata.Name
	eksAPI := c.ctl.AWSProvider.EKS()

	// get all managed nodegroups for this cluster
	nodeGroups, err := eksAPI.ListNodegroups(ctx, &awseks.ListNodegroupsInput{
		ClusterName: &clusterName,
	})
	if err != nil {
		return err
	}

	if len(allStacks) == 0 && len(nodeGroups.Nodegroups) == 0 {
		logger.Warning("no nodegroups found for %s", clusterName)
		return nil
	}

	// we kill every nodegroup with a stack the standard way. wait is always true
	tasks, err := c.stackManager.NewTasksToDeleteNodeGroups(allStacks, func(_ string) bool { return true }, true, nil)
	if err != nil {
		return err
	}

	for _, n := range nodeGroups.Nodegroups {
		isUnowned := func() bool {
			for _, stack := range allStacks {
				if stack.NodeGroupName == n {
					return false
				}
			}
			return true
		}

		if isUnowned() {
			// if a managed ng does not have a stack, we queue it for deletion via api
			tasks.Append(c.stackManager.NewTaskToDeleteUnownedNodeGroup(ctx, clusterName, n, eksAPI, c.waitForUnownedNgsDeletion(ctx, waitInterval)))
		}
	}

	// TODO what dis?
	tasks.PlanMode = false
	logger.Info(tasks.Describe())
	if errs := tasks.DoAllSync(); len(errs) > 0 {
		return handleErrors(errs, "nodegroup(s)")
	}
	return nil
}

func isNotFound(err error) bool {
	var notFoundErr *ekstypes.ResourceNotFoundException
	return errors.As(err, &notFoundErr)
}

func (c *UnownedCluster) waitForUnownedNgsDeletion(ctx context.Context, interval time.Duration) *manager.DeleteWaitCondition {
	condition := func() (bool, error) {
		nodeGroups, err := c.ctl.AWSProvider.EKS().ListNodegroups(ctx, &awseks.ListNodegroupsInput{
			ClusterName: &c.cfg.Metadata.Name,
		})
		if err != nil {
			return false, err
		}
		if len(nodeGroups.Nodegroups) == 0 {
			return true, nil
		}

		logger.Info("waiting for all non eksctl-owned nodegroups to be deleted")
		return false, nil
	}

	return &manager.DeleteWaitCondition{
		Condition: condition,
		Timeout:   c.ctl.AWSProvider.WaitTimeout(),
		Interval:  interval,
	}
}
