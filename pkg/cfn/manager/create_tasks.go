package manager

import (
	"fmt"

	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// NewTasksToCreateClusterWithNodeGroups defines all tasks required to create a cluster along
// with some nodegroups; see CreateAllNodeGroups for how onlyNodeGroupSubset works
func (c *StackCollection) NewTasksToCreateClusterWithNodeGroups(nodeGroups []*v1alpha5.NodeGroup) *TaskTree {
	tasks := &TaskTree{Parallel: false}

	tasks.Append(
		&taskWithoutParams{
			info: fmt.Sprintf("create cluster control plane %q", c.spec.Metadata.Name),
			call: c.createClusterTask,
		},
	)

	nodeGroupTasks := c.NewTasksToCreateNodeGroups(nodeGroups)
	if nodeGroupTasks.Len() > 0 {
		nodeGroupTasks.IsSubTask = true
		tasks.Append(nodeGroupTasks)
	}

	return tasks
}

// NewTasksToCreateNodeGroups defines tasks required to create all of the nodegroups
func (c *StackCollection) NewTasksToCreateNodeGroups(nodeGroups []*v1alpha5.NodeGroup) *TaskTree {
	tasks := &TaskTree{Parallel: true}

	for _, ng := range nodeGroups {
		tasks.Append(&taskWithNodeGroupSpec{
			info:      fmt.Sprintf("create nodegroup %q", ng.Name),
			nodeGroup: ng,
			call:      c.createNodeGroupTask,
		})
	}

	return tasks
}
