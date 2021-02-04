package manager

import (
	"fmt"

	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

// NewTasksToCreateClusterWithNodeGroups defines all tasks required to create a cluster along
// with some nodegroups; see CreateAllNodeGroups for how onlyNodeGroupSubset works
func (c *StackCollection) NewTasksToCreateClusterWithNodeGroups(nodeGroups []*api.NodeGroup,
	managedNodeGroups []*api.ManagedNodeGroup, supportsManagedNodes bool, postClusterCreationTasks ...tasks.Task) *tasks.TaskTree {

	taskTree := tasks.TaskTree{Parallel: false}

	taskTree.Append(
		&createClusterTask{
			info:                 fmt.Sprintf("create cluster control plane %q", c.spec.Metadata.Name),
			stackCollection:      c,
			supportsManagedNodes: supportsManagedNodes,
		},
	)

	appendNodeGroupTasksTo := func(taskTree *tasks.TaskTree) {
		nodeGroupTasks := c.NewUnmanagedNodeGroupTask(nodeGroups, supportsManagedNodes, false)

		managedNodeGroupTasks := c.NewManagedNodeGroupTask(managedNodeGroups, false)
		if managedNodeGroupTasks.Len() > 0 {
			nodeGroupTasks.Append(managedNodeGroupTasks.Tasks...)
		}

		if nodeGroupTasks.Len() > 0 {
			nodeGroupTasks.IsSubTask = true
			taskTree.Append(nodeGroupTasks)
		}
	}

	if len(postClusterCreationTasks) > 0 {
		postClusterCreationTaskTree := tasks.TaskTree{
			Parallel:  false,
			IsSubTask: true,
		}
		postClusterCreationTaskTree.Append(postClusterCreationTasks...)
		appendNodeGroupTasksTo(&postClusterCreationTaskTree)
		taskTree.Append(&postClusterCreationTaskTree)
	} else {
		appendNodeGroupTasksTo(&taskTree)
	}

	return &taskTree
}

// NewUnmanagedNodeGroupTask defines tasks required to create all of the nodegroups
func (c *StackCollection) NewUnmanagedNodeGroupTask(nodeGroups []*api.NodeGroup, supportsManagedNodes bool, forceAddCNIPolicy bool) *tasks.TaskTree {
	taskTree := &tasks.TaskTree{Parallel: true}

	for _, ng := range nodeGroups {
		taskTree.Append(&nodeGroupTask{
			info:                 fmt.Sprintf("create nodegroup %q", ng.NameString()),
			nodeGroup:            ng,
			stackCollection:      c,
			supportsManagedNodes: supportsManagedNodes,
			forceAddCNIPolicy:    forceAddCNIPolicy,
		})
		// TODO: move authconfigmap tasks here using kubernetesTask and kubernetes.CallbackClientSet
	}

	return taskTree
}

// NewManagedNodeGroupTask defines tasks required to create managed nodegroups
func (c *StackCollection) NewManagedNodeGroupTask(nodeGroups []*api.ManagedNodeGroup, forceAddCNIPolicy bool) *tasks.TaskTree {
	taskTree := &tasks.TaskTree{Parallel: true}
	for _, ng := range nodeGroups {
		taskTree.Append(&managedNodeGroupTask{
			stackCollection:   c,
			nodeGroup:         ng,
			forceAddCNIPolicy: forceAddCNIPolicy,
			info:              fmt.Sprintf("create managed nodegroup %q", ng.Name),
		})
	}
	return taskTree
}

// NewClusterCompatTask creates a new task that checks for cluster compatibility with new features like
// Managed Nodegroups and Fargate, and updates the CloudFormation cluster stack if the required resources are missing
func (c *StackCollection) NewClusterCompatTask() tasks.Task {
	return &clusterCompatTask{
		stackCollection: c,
		info:            "fix cluster compatibility",
	}
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

		saTasks.Append(&taskWithClusterIAMServiceAccountSpec{
			info:            fmt.Sprintf("create IAM role for serviceaccount %q", sa.NameString()),
			stackCollection: c,
			serviceAccount:  sa,
			oidc:            oidc,
		})

		saTasks.Append(&kubernetesTask{
			info:       fmt.Sprintf("create serviceaccount %q", sa.NameString()),
			kubernetes: clientSetGetter,
			call: func(clientSet kubernetes.Interface) error {
				sa.SetAnnotations()
				if err := kubernetes.MaybeCreateServiceAccountOrUpdateMetadata(clientSet, sa.ClusterIAMMeta.AsObjectMeta()); err != nil {
					return errors.Wrapf(err, "failed to create service account %s", sa.NameString())
				}
				return nil
			},
		})

		taskTree.Append(saTasks)
	}
	return taskTree
}
