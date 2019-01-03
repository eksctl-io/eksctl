package manager

import (
	"sync"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha3"
)

type taskFunc func(chan error, interface{}) error

// Task has a function with an opaque payload
type Task struct {
	Call taskFunc
	Data interface{}
}

// Run a set of tests in parallel and wait for them to complete;
// passError should take any errors and do what it needs to in
// a given context, e.g. during serial CLI-driven execution one
// can keep errors in a slice, while in a daemon channel maybe
// more suitable
func Run(passError func(error), tasks ...Task) {
	wg := &sync.WaitGroup{}
	wg.Add(len(tasks))
	for t := range tasks {
		go func(t int) {
			defer wg.Done()
			logger.Debug("task %d started", t)
			errs := make(chan error)
			if err := tasks[t].Call(errs, tasks[t].Data); err != nil {
				passError(err)
				return
			}
			if err := <-errs; err != nil {
				passError(err)
				return
			}
			logger.Debug("task %d returned without errors", t)
		}(t)
	}
	logger.Debug("waiting for %d tasks to complete", len(tasks))
	wg.Wait()
}

// RunSingleTask runs a task with a proper error handling
func (c *StackCollection) RunSingleTask(t Task) []error {
	errs := []error{}
	appendErr := func(err error) {
		errs = append(errs, err)
	}
	if Run(appendErr, t); len(errs) > 0 {
		return errs
	}
	return nil
}

// CreateClusterWithNodeGroups runs all tasks required to create
// the stacks (a cluster and one or more nodegroups); any errors
// will be returned as a slice as soon as one of the tasks or group
// of tasks is completed
func (c *StackCollection) CreateClusterWithNodeGroups() []error {
	if errs := c.RunSingleTask(Task{c.CreateCluster, nil}); len(errs) > 0 {
		return errs
	}

	return c.CreateAllNodeGroups()
}

// CreateAllNodeGroups runs all tasks required to create the node groups;
// any errors will be returned as a slice as soon as one of the tasks
// or group of tasks is completed
func (c *StackCollection) CreateAllNodeGroups() []error {
	errs := []error{}
	appendErr := func(err error) {
		errs = append(errs, err)
	}

	createAllNodeGroups := []Task{}
	for i := range c.spec.NodeGroups {
		t := Task{
			Call: c.CreateNodeGroup,
			Data: c.spec.NodeGroups[i],
		}
		createAllNodeGroups = append(createAllNodeGroups, t)
	}
	if Run(appendErr, createAllNodeGroups...); len(errs) > 0 {
		return errs
	}

	return nil
}

// CreateOneNodeGroup runs a task to create a single node groups;
// any errors will be returned as a slice as soon as the tasks is
// completed
func (c *StackCollection) CreateOneNodeGroup(ng *api.NodeGroup) []error {
	return c.RunSingleTask(Task{
		Call: c.CreateNodeGroup,
		Data: ng,
	})
}

// DeleteAllNodeGroups deletes all nodegroups without waiting
func (c *StackCollection) DeleteAllNodeGroups(call taskFunc) []error {
	nodeGroupStacks, err := c.DescribeNodeGroupStacks()
	if err != nil {
		return []error{err}
	}

	errs := []error{}
	for _, s := range nodeGroupStacks {
		if err := c.DeleteNodeGroup(getNodeGroupName(s)); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

// WaitDeleteAllNodeGroups runs all tasks required to delete all the nodegroup
// stacks and wait for all nodegroups to be deleted; any errors will be returned
// as a slice as soon as the group of tasks is completed
func (c *StackCollection) WaitDeleteAllNodeGroups() []error {
	nodeGroupStacks, err := c.DescribeNodeGroupStacks()
	if err != nil {
		return []error{err}
	}

	errs := []error{}
	appendErr := func(err error) {
		errs = append(errs, err)
	}

	deleteAllNodeGroups := []Task{}
	for i := range nodeGroupStacks {
		t := Task{
			Call: c.WaitDeleteNodeGroup,
			Data: getNodeGroupName(nodeGroupStacks[i]),
		}
		deleteAllNodeGroups = append(deleteAllNodeGroups, t)
	}
	if Run(appendErr, deleteAllNodeGroups...); len(errs) > 0 {
		return errs
	}

	return nil
}
