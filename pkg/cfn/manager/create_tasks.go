package manager

import (
	"context"
	"fmt"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

// NewTasksToCreateClusterWithNodeGroups defines all tasks required to create a cluster along
// with some nodegroups; see CreateAllNodeGroups for how onlyNodeGroupSubset works.
func (c *StackCollection) NewTasksToCreateClusterWithNodeGroups(ctx context.Context, nodeGroups []*api.NodeGroup,
	managedNodeGroups []*api.ManagedNodeGroup, postClusterCreationTasks ...tasks.Task) *tasks.TaskTree {

	taskTree := tasks.TaskTree{Parallel: false}

	taskTree.Append(
		&createClusterTask{
			info:                 fmt.Sprintf("create cluster control plane %q", c.spec.Metadata.Name),
			stackCollection:      c,
			supportsManagedNodes: true,
			ctx:                  ctx,
		},
	)

	appendNodeGroupTasksTo := func(taskTree *tasks.TaskTree) {
		vpcImporter := vpc.NewStackConfigImporter(c.MakeClusterStackName())
		nodeGroupTasks := &tasks.TaskTree{
			Parallel:  true,
			IsSubTask: true,
		}
		if unmanagedNodeGroupTasks := c.NewUnmanagedNodeGroupTask(ctx, nodeGroups, false, false, vpcImporter); unmanagedNodeGroupTasks.Len() > 0 {
			unmanagedNodeGroupTasks.IsSubTask = true
			nodeGroupTasks.Append(unmanagedNodeGroupTasks)
		}
		if managedNodeGroupTasks := c.NewManagedNodeGroupTask(ctx, managedNodeGroups, false, vpcImporter); managedNodeGroupTasks.Len() > 0 {
			managedNodeGroupTasks.IsSubTask = true
			nodeGroupTasks.Append(managedNodeGroupTasks)
		}

		if nodeGroupTasks.Len() > 0 {
			taskTree.Append(nodeGroupTasks)
		}
	}

	if len(postClusterCreationTasks) > 0 {
		postClusterCreationTaskTree := &tasks.TaskTree{
			Parallel:  false,
			IsSubTask: true,
		}
		postClusterCreationTaskTree.Append(postClusterCreationTasks...)
		appendNodeGroupTasksTo(postClusterCreationTaskTree)
		taskTree.Append(postClusterCreationTaskTree)
	} else {
		appendNodeGroupTasksTo(&taskTree)
	}

	return &taskTree
}

// NewUnmanagedNodeGroupTask defines tasks required to create all of the nodegroups
func (c *StackCollection) NewUnmanagedNodeGroupTask(ctx context.Context, nodeGroups []*api.NodeGroup, forceAddCNIPolicy, skipEgressRules bool, vpcImporter vpc.Importer) *tasks.TaskTree {
	taskTree := &tasks.TaskTree{Parallel: true}

	for _, ng := range nodeGroups {
		taskTree.Append(&nodeGroupTask{
			info:              fmt.Sprintf("create nodegroup %q", ng.NameString()),
			ctx:               ctx,
			nodeGroup:         ng,
			stackCollection:   c,
			forceAddCNIPolicy: forceAddCNIPolicy,
			vpcImporter:       vpcImporter,
			skipEgressRules:   skipEgressRules,
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
