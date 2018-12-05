package manager

import (
	"sync"

	"github.com/kris-nova/logger"
)

type task struct {
	call func(chan error, interface{}) error
	data interface{}
}

// Run a set of tests in parallel and wait for them to complete;
// passError should take any errors and do what it needs to in
// a given context, e.g. during serial CLI-driven execution one
// can keep errors in a slice, while in a daemon channel maybe
// more suitable
func Run(passError func(error), tasks ...task) {
	wg := &sync.WaitGroup{}
	wg.Add(len(tasks))
	for t := range tasks {
		go func(t int) {
			defer wg.Done()
			logger.Debug("task %d started", t)
			errs := make(chan error)
			if err := tasks[t].call(errs, tasks[t].data); err != nil {
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

// CreateClusterWithNodeGroups runs all tasks required to create
// the stacks (a cluster and one or more nodegroups); any errors
// will be returned as a slice as soon as one of the tasks or group
// of tasks is completed
func (s *StackCollection) CreateClusterWithNodeGroups() []error {
	errs := []error{}
	appendErr := func(err error) {
		errs = append(errs, err)
	}
	if Run(appendErr, task{s.CreateCluster, nil}); len(errs) > 0 {
		return errs
	}

	createAllNodeGroups := []task{}
	for i := range s.spec.NodeGroups {
		t := task{
			call: s.CreateNodeGroup,
			data: s.spec.NodeGroups[i],
		}
		createAllNodeGroups = append(createAllNodeGroups, t)
	}
	if Run(appendErr, createAllNodeGroups...); len(errs) > 0 {
		return errs
	}

	return nil
}

// deleteAllNodeGroupsTasks returns a list of tasks for deleting all the
// nodegroup stacks
func (s *StackCollection) deleteAllNodeGroupsTasks(call func(chan error, interface{}) error) ([]task, error) {
	stacks, err := s.listAllNodeGroups()
	if err != nil {
		return nil, err
	}
	deleteAllNodeGroups := []task{}
	for i := range stacks {
		t := task{
			call: call,
			data: stacks[i],
		}
		deleteAllNodeGroups = append(deleteAllNodeGroups, t)
	}
	return deleteAllNodeGroups, nil
}

// DeleteAllNodeGroups runs all tasks required to delete all the nodegroup
// stacks; any errors will be returned as a slice as soon as the group
// of tasks is completed
func (s *StackCollection) DeleteAllNodeGroups() []error {
	errs := []error{}
	appendErr := func(err error) {
		errs = append(errs, err)
	}

	deleteAllNodeGroups, err := s.deleteAllNodeGroupsTasks(s.DeleteNodeGroup)
	if err != nil {
		appendErr(err)
		return errs
	}

	if Run(appendErr, deleteAllNodeGroups...); len(errs) > 0 {
		return errs
	}

	return nil
}

// WaitDeleteAllNodeGroups runs all tasks required to delete all the nodegroup
// stacks, it waits for each nodegroup to get deleted; any errors will be
// returned as a slice as soon as the group of tasks is completed
func (s *StackCollection) WaitDeleteAllNodeGroups() []error {
	errs := []error{}
	appendErr := func(err error) {
		errs = append(errs, err)
	}

	deleteAllNodeGroups, err := s.deleteAllNodeGroupsTasks(s.WaitDeleteNodeGroup)
	if err != nil {
		appendErr(err)
		return errs
	}

	if Run(appendErr, deleteAllNodeGroups...); len(errs) > 0 {
		return errs
	}

	return nil
}
