package manager

import (
	"fmt"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"k8s.io/apimachinery/pkg/util/sets"
)

// NewTasksToCreateClusterWithNodeGroups defines all tasks required to create a cluster along
// with some nodegroups; see CreateAllNodeGroups for how onlyNodeGroupSubset works
func (c *StackCollection) NewTasksToCreateClusterWithNodeGroups(onlyNodeGroupSubset sets.String) (*TaskTree, error) {
	tasks := &TaskTree{Parallel: false}

	// Control plane.
	{
		tasks.Append(
			&taskWithoutParams{
				info: fmt.Sprintf("create cluster control plane %q", c.spec.Metadata.Name),
				call: c.createClusterTask,
			},
		)
	}

	// Nodegroups.
	{
		nodeGroupTasks, err := c.NewTasksToCreateNodeGroups(onlyNodeGroupSubset)
		if err != nil {
			return nil, err
		}

		if nodeGroupTasks.Len() > 0 {
			nodeGroupTasks.IsSubTask = true
			tasks.Append(nodeGroupTasks)
		}
	}

	return tasks, nil
}

// NewTasksToCreateNodeGroups defines tasks required to create all of the nodegroups if
// onlySubset is nil, otherwise just the tasks for nodegroups that are in onlySubset
// will be defined
func (c *StackCollection) NewTasksToCreateNodeGroups(onlySubset sets.String) (*TaskTree, error) {
	tasks := &TaskTree{Parallel: true}

	// Spotinst.
	{
		oceanTasks, err := c.newTasksToCreateNodeGroupSpotinstOcean()
		if err != nil {
			return nil, err
		}

		if oceanTasks.Len() > 0 {
			oceanTasks.IsSubTask = true
			tasks.Append(oceanTasks)
		}
	}

	// Nodegroups.
	{
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
	}

	return tasks, nil
}

func (c *StackCollection) newTasksToCreateNodeGroupSpotinstOcean() (*TaskTree, error) {
	tasks := &TaskTree{Parallel: true}
	var ng *api.NodeGroup

	// Single node group.
	if len(c.spec.NodeGroups) == 1 && c.spec.NodeGroups[0].Spotinst != nil {
		ng = c.spec.NodeGroups[0]
	}

	// Multiple node groups.
	if len(c.spec.NodeGroups) > 1 {
		for _, g := range c.spec.NodeGroups {
			if g.Spotinst.Ocean == nil || g.Spotinst.Ocean.DefaultLaunchSpec == nil {
				continue
			}
			if *g.Spotinst.Ocean.DefaultLaunchSpec {
				if ng != nil {
					return nil, fmt.Errorf("unable to detect default ocean launch spec: " +
						"multiple nodegroups configured with `ocean.defaultLaunchSpec: \"true\"`")
				}

				ng = g
			}
		}

		if ng == nil {
			return nil, fmt.Errorf("unable to detect default ocean launch spec: " +
				"please configure the desired default nodegroup with `ocean.defaultLaunchSpec: \"true\"`")
		}
	}

	// Configure all tasks.
	if ng != nil {
		ng.Name = "ocean"

		tasks.Append(&taskWithNodeGroupSpec{
			info:      fmt.Sprintf("create spotinst ocean %q", c.spec.Metadata.Name),
			nodeGroup: ng,
			call:      c.createNodeGroupTask,
		})
	}

	return tasks, nil
}
