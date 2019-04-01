package manager

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/cloudformation"

	"k8s.io/apimachinery/pkg/util/sets"
)

// DeleteTasksForClusterWithNodeGroups defines tasks required to delete all the nodegroup
// stacks and the cluster
func (c *StackCollection) DeleteTasksForClusterWithNodeGroups(wait bool, cleanup func(chan error, string) error) (*TaskTree, error) {
	tasks := &TaskTree{Parallel: false}

	nodeGroupTasks, err := c.DeleteTasksForNodeGroups(nil, true, cleanup)
	if err != nil {
		return nil, err
	}
	if nodeGroupTasks.Len() > 0 {
		nodeGroupTasks.Sub = true
		tasks.Append(nodeGroupTasks)
	}

	clusterStack, err := c.DescribeClusterStack()
	if err != nil {
		return nil, err
	}

	info := fmt.Sprintf("delete cluster control plane %q", c.spec.Metadata.Name)
	if wait {
		tasks.Append(&taskWithStackSpec{
			info:  info,
			stack: clusterStack,
			call:  c.WaitDeleteStackBySpec,
		})
	} else {
		tasks.Append(&asyncTaskWithStackSpec{
			info:  info,
			stack: clusterStack,
			call:  c.DeleteStackBySpec,
		})
	}

	return tasks, nil
}

// DeleteTasksForNodeGroups defines tasks required to delete all of the nodegroups if
// onlySubset is nil, otherwise just the tasks for nodegroups that are in onlySubset
// will be defined
func (c *StackCollection) DeleteTasksForNodeGroups(onlySubset sets.String, wait bool, cleanup func(chan error, string) error) (*TaskTree, error) {
	nodeGroupStacks, err := c.DescribeNodeGroupStacks()
	if err != nil {
		return nil, err
	}

	tasks := &TaskTree{Parallel: true}

	for _, s := range nodeGroupStacks {
		name := getNodeGroupName(s)
		if onlySubset != nil && !onlySubset.Has(name) {
			continue
		}
		if *s.StackStatus == cloudformation.StackStatusDeleteFailed && cleanup != nil {
			tasks.Append(&taskWithNameParam{
				info: fmt.Sprintf("cleanup for nodegroup %q", name),
				call: cleanup,
			})
		}
		info := fmt.Sprintf("delete nodegroup %q", name)
		if wait {
			tasks.Append(&taskWithStackSpec{
				info:  info,
				stack: s,
				call:  c.WaitDeleteStackBySpec,
			})
		} else {
			tasks.Append(&asyncTaskWithStackSpec{
				info:  info,
				stack: s,
				call:  c.DeleteStackBySpec,
			})
		}
	}

	return tasks, nil
}
