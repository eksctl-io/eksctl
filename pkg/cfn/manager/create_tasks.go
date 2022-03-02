package manager

import (
	"context"
	"fmt"

	"github.com/weaveworks/eksctl/pkg/spot"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

const (
	managedByKubernetesLabelKey   = "app.kubernetes.io/managed-by"
	managedByKubernetesLabelValue = "eksctl"
)

// NewTasksToCreateClusterWithNodeGroups defines all tasks required to create a cluster along
// with some nodegroups; see CreateAllNodeGroups for how onlyNodeGroupSubset works.
func (c *StackCollection) NewTasksToCreateClusterWithNodeGroups(ctx context.Context, nodeGroups []*api.NodeGroup,
	managedNodeGroups []*api.ManagedNodeGroup, postClusterCreationTasks ...tasks.Task) (*tasks.TaskTree, error) {

	taskTree := tasks.TaskTree{Parallel: false}

	// Control plane.
	{
		taskTree.Append(
			&createClusterTask{
				info:                 fmt.Sprintf("create cluster control plane %q", c.spec.Metadata.Name),
				stackCollection:      c,
				supportsManagedNodes: true,
				ctx:                  ctx,
			},
		)
	}

	// Nodegroups.
	{
		vpcImporter := vpc.NewStackConfigImporter(c.MakeClusterStackName())
		nodeGroupTaskTree, err := c.NewNodeGroupTask(ctx, nodeGroups, managedNodeGroups, false, vpcImporter)
		if err != nil {
			return nil, err
		}

		if nodeGroupTaskTree.Len() > 0 {
			nodeGroupTaskTree.IsSubTask = true
			taskTree.Append(nodeGroupTaskTree)
		}
	}

	// Post creation tasks.
	{
		if len(postClusterCreationTasks) > 0 {
			postTaskTree := &tasks.TaskTree{
				Parallel:  false,
				IsSubTask: true,
			}
			postTaskTree.Append(postClusterCreationTasks...)
			taskTree.Append(postTaskTree)
		}
	}

	return &taskTree, nil
}

// NewNodeGroupTask defines tasks required to create all of the nodegroups
func (c *StackCollection) NewNodeGroupTask(ctx context.Context, nodeGroups []*api.NodeGroup, managedNodeGroups []*api.ManagedNodeGroup,
	forceAddCNIPolicy bool, vpcImporter vpc.Importer) (*tasks.TaskTree, error) {
	taskTree := &tasks.TaskTree{Parallel: true}

	// Spot Ocean.
	{
		oceanTaskTree, err := c.NewSpotOceanNodeGroupTask(ctx, vpcImporter)
		if err != nil {
			return nil, err
		}
		if oceanTaskTree.Len() > 0 {
			oceanTaskTree.IsSubTask = true
			taskTree.Parallel = false
			taskTree.Append(oceanTaskTree)
		}
	}

	// Managed.
	{
		managedNodeGroupTaskTree := c.NewManagedNodeGroupTask(ctx, managedNodeGroups, forceAddCNIPolicy, vpcImporter)
		if managedNodeGroupTaskTree.Len() > 0 {
			managedNodeGroupTaskTree.IsSubTask = true
			taskTree.Append(managedNodeGroupTaskTree)
		}
	}

	// Unmanaged.
	{
		nodeGroupTaskTree := c.NewUnmanagedNodeGroupTask(ctx, nodeGroups, forceAddCNIPolicy, vpcImporter)
		if nodeGroupTaskTree.Len() > 0 {
			nodeGroupTaskTree.IsSubTask = true
			taskTree.Append(nodeGroupTaskTree)
		}
	}

	return taskTree, nil
}

// NewUnmanagedNodeGroupTask defines tasks required to create all of the nodegroups
func (c *StackCollection) NewUnmanagedNodeGroupTask(ctx context.Context, nodeGroups []*api.NodeGroup, forceAddCNIPolicy bool, vpcImporter vpc.Importer) *tasks.TaskTree {
	taskTree := &tasks.TaskTree{Parallel: true}

	for _, ng := range nodeGroups {
		taskTree.Append(&nodeGroupTask{
			info:              fmt.Sprintf("create nodegroup %q", ng.NameString()),
			ctx:               ctx,
			nodeGroup:         ng,
			stackCollection:   c,
			forceAddCNIPolicy: forceAddCNIPolicy,
			vpcImporter:       vpcImporter,
		})
		// TODO: move authconfigmap tasks here using kubernetesTask and kubernetes.CallbackClientSet
	}

	return taskTree
}

// NewManagedNodeGroupTask defines tasks required to create managed nodegroups
func (c *StackCollection) NewManagedNodeGroupTask(ctx context.Context, nodeGroups []*api.ManagedNodeGroup, forceAddCNIPolicy bool, vpcImporter vpc.Importer) *tasks.TaskTree {
	taskTree := &tasks.TaskTree{Parallel: true}
	for _, ng := range nodeGroups {
		// Disable parallelisation if any tags propagation is done
		// since nodegroup must be created to propagate tags to its ASGs.
		subTask := &tasks.TaskTree{
			Parallel:  false,
			IsSubTask: true,
		}
		subTask.Append(&managedNodeGroupTask{
			stackCollection:   c,
			nodeGroup:         ng,
			forceAddCNIPolicy: forceAddCNIPolicy,
			vpcImporter:       vpcImporter,
			info:              fmt.Sprintf("create managed nodegroup %q", ng.Name),
			ctx:               ctx,
		})
		if api.IsEnabled(ng.PropagateASGTags) {
			subTask.Append(&managedNodeGroupTagsToASGPropagationTask{
				stackCollection: c,
				nodeGroup:       ng,
				info:            fmt.Sprintf("propagate tags to ASG for managed nodegroup %q", ng.Name),
				ctx:             ctx,
			})
		}
		taskTree.Append(subTask)
	}
	return taskTree
}

// NewTasksToCreateIAMServiceAccounts defines tasks required to create all of the IAM ServiceAccounts
func (c *StackCollection) NewTasksToCreateIAMServiceAccounts(serviceAccounts []*api.ClusterIAMServiceAccount, oidc *iamoidc.OpenIDConnectManager, clientSetGetter kubernetes.ClientSetGetter) *tasks.TaskTree {
	taskTree := &tasks.TaskTree{Parallel: true}

	for i := range serviceAccounts {
		sa := serviceAccounts[i]
		saTasks := &tasks.TaskTree{
			Parallel:  false,
			IsSubTask: true,
		}

		if sa.AttachRoleARN == "" {
			saTasks.Append(&taskWithClusterIAMServiceAccountSpec{
				info:            fmt.Sprintf("create IAM role for serviceaccount %q", sa.NameString()),
				stackCollection: c,
				serviceAccount:  sa,
				oidc:            oidc,
			})
		} else {
			logger.Debug("attachRoleARN was provided, skipping role creation")
			sa.Status = &api.ClusterIAMServiceAccountStatus{
				RoleARN: &sa.AttachRoleARN,
			}
		}

		if sa.Labels == nil {
			sa.Labels = make(map[string]string)
		}
		sa.Labels[managedByKubernetesLabelKey] = managedByKubernetesLabelValue
		if !api.IsEnabled(sa.RoleOnly) {
			saTasks.Append(&kubernetesTask{
				info:       fmt.Sprintf("create serviceaccount %q", sa.NameString()),
				kubernetes: clientSetGetter,
				objectMeta: sa.ClusterIAMMeta.AsObjectMeta(),
				call: func(clientSet kubernetes.Interface, objectMeta v1.ObjectMeta) error {
					sa.SetAnnotations()
					objectMeta.SetAnnotations(sa.AsObjectMeta().Annotations)
					objectMeta.SetLabels(sa.AsObjectMeta().Labels)
					if err := kubernetes.MaybeCreateServiceAccountOrUpdateMetadata(clientSet, objectMeta); err != nil {
						return errors.Wrapf(err, "failed to create service account %s/%s", objectMeta.GetNamespace(), objectMeta.GetName())
					}
					return nil
				},
			})
		}

		taskTree.Append(saTasks)
	}
	return taskTree
}

// NewSpotOceanNodeGroupTask defines tasks required to create Ocean Cluster.
func (c *StackCollection) NewSpotOceanNodeGroupTask(ctx context.Context, vpcImporter vpc.Importer) (*tasks.TaskTree, error) {
	taskTree := &tasks.TaskTree{Parallel: true}

	// Check whether the Ocean Cluster should be created.
	stacks, err := c.ListNodeGroupStacks(ctx)
	if err != nil {
		return nil, err
	}
	ng := spot.ShouldCreateOceanCluster(c.spec, stacks)
	if ng == nil { // already exists OR --without-nodegroup
		return taskTree, nil
	}

	// Allow post-create actions on this nodegroup.
	c.spec.NodeGroups = append(c.spec.NodeGroups, ng)

	// Add a new task.
	taskTree.Append(&nodeGroupTask{
		info:            "create ocean cluster",
		nodeGroup:       ng,
		stackCollection: c,
		vpcImporter:     vpcImporter,
		ctx:             ctx,
	})

	return taskTree, nil
}
