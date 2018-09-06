package manager

import (
	"sync"

	"github.com/kubicorn/kubicorn/pkg/logger"
)

type task func(chan error) error

// Run a set of tests in parallel and wait for them to complete;
// passError should take any errors and do what it needs to in
// a given context, e.g. duing serial CLI-driven execution one
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
			if err := tasks[t](errs); err != nil {
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

// CreateClusterWithInitialNodeGroup runs two tasks to create
// the stacks for use with CLI; any errors will be returned
// as a slice on completion of one of the two tasks
func (s *StackCollection) CreateClusterWithInitialNodeGroup() []error {
	errs := []error{}
	appendErr := func(err error) {
		errs = append(errs, err)
	}
	if Run(appendErr, s.CreateCluster); len(errs) > 0 {
		return errs
	}
	if Run(appendErr, s.CreateInitialNodeGroup); len(errs) > 0 {
		return errs
	}
	return nil
}

func (s *StackCollection) ScaleInitialNodeGroup() []error {
	errs := []error{}
	appendErr := func(err error) {
		errs = append(errs, err)
	}
	if Run(appendErr, s.ScaleNodeGroup); len(errs) > 0 {
		return errs
	}

	return nil
}
