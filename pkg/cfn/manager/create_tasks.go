package manager

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"
)

// NewTasksToCreateClusterWithNodeGroups defines all tasks required to create a cluster along
// with some nodegroups; see CreateAllNodeGroups for how onlyNodeGroupSubset works
func (c *StackCollection) NewTasksToCreateClusterWithNodeGroups(onlyNodeGroupSubset sets.String) *TaskTree {
	tasks := &TaskTree{Parallel: false}

	tasks.Append(
		&taskWithoutParams{
			info: fmt.Sprintf("create cluster control plane %q", c.spec.Metadata.Name),
			call: c.createClusterTask,
		},
	)

	nodeGroupTasks := c.NewTasksToCreateNodeGroups(onlyNodeGroupSubset)
	if nodeGroupTasks.Len() > 0 {
		nodeGroupTasks.IsSubTask = true
		tasks.Append(nodeGroupTasks)
	}

	return tasks
}

// NewTasksToCreateNodeGroups defines tasks required to create all of the nodegroups if
// onlySubset is nil, otherwise just the tasks for nodegroups that are in onlySubset
// will be defined
func (c *StackCollection) NewTasksToCreateNodeGroups(onlySubset sets.String) *TaskTree {
	tasks := &TaskTree{Parallel: true}

	for i := range c.spec.NodeGroups {
		ng := c.spec.NodeGroups[i]
		if onlySubset != nil && !onlySubset.Has(ng.Name) {
			continue
		}
		tasks.Append(&taskWithNodeGroupSpec{
			info:      fmt.Sprintf("create nodegroup %q", ng.Name),
			nodeGroup: ng,
			call:      c.createNodeGroupTask,
		})
	}

	return tasks
}
