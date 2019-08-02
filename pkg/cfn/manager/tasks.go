package manager

import (
	"fmt"
	"strings"
	"sync"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// Task is a common interface for the stack manager tasks
type Task interface {
	Describe() string
	Do(chan error) error
}

// TaskTree wraps a set of tasks
type TaskTree struct {
	tasks     []Task
	Parallel  bool
	PlanMode  bool
	IsSubTask bool
}

// Append new tasks to the set
func (t *TaskTree) Append(newTasks ...Task) {
	t.tasks = append(t.tasks, newTasks...)
}

// Len returns number of tasks in the set
func (t *TaskTree) Len() int {
	if t == nil {
		return 0
	}
	return len(t.tasks)
}

// Describe the set
func (t *TaskTree) Describe() string {
	descriptions := []string{}
	for _, task := range t.tasks {
		descriptions = append(descriptions, task.Describe())
	}
	mode := "sequential"
	if t.Parallel {
		mode = "parallel"
	}
	count := len(descriptions)
	var msg string
	noun := "task"
	if t.IsSubTask {
		noun = "sub-task"
	}
	switch count {
	case 0:
		msg = "no tasks"
	case 1:
		msg = fmt.Sprintf("1 %s: { %s }", noun, descriptions[0])
		if t.IsSubTask {
			msg = descriptions[0] // simple description for single sub-task
		}
	default:
		noun += "s"
		msg = fmt.Sprintf("%d %s %s: { %s }", count, mode, noun, strings.Join(descriptions, ", "))
	}
	if t.PlanMode {
		return "(plan) " + msg
	}
	return msg
}

// Do will run through the set in the backround, it may return an error immediately,
// or eventually write to the errs channel; it will close the channel once all tasks
// are completed
func (t *TaskTree) Do(allErrs chan error) error {
	if t.Len() == 0 || t.PlanMode {
		logger.Debug("no actual tasks")
		close(allErrs)
		return nil
	}

	errs := make(chan error)

	if t.Parallel {
		go doParallelTasks(errs, t.tasks)
	} else {
		go doSequentialTasks(errs, t.tasks)
	}

	go func() {
		defer close(allErrs)
		for err := range errs {
			allErrs <- err
		}
	}()

	return nil
}

// DoAllSync will run through the set in the foregounds and return all the errors
// in a slice
func (t *TaskTree) DoAllSync() []error {
	if t.Len() == 0 || t.PlanMode {
		logger.Debug("no actual tasks")
		return nil
	}

	errs := make(chan error)

	if t.Parallel {
		go doParallelTasks(errs, t.tasks)
	} else {
		go doSequentialTasks(errs, t.tasks)
	}

	allErrs := []error{}
	for err := range errs {
		allErrs = append(allErrs, err)
	}
	return allErrs
}

func info(skip bool, info string) string {
	if skip {
		return "(skip) " + info
	}
	return info
}

type taskWithoutParams struct {
	info string
	skip bool
	call func(chan error) error
}

func (t *taskWithoutParams) Skip() bool               { return t.skip }
func (t *taskWithoutParams) Describe() string         { return info(t.skip, t.info) }
func (t *taskWithoutParams) Do(errs chan error) error { return t.call(errs) }

type taskWithNameParam struct {
	info string
	skip bool
	name string
	call func(chan error, string) error
}

func (t *taskWithNameParam) Skip() bool               { return t.skip }
func (t *taskWithNameParam) Describe() string         { return info(t.skip, t.info) }
func (t *taskWithNameParam) Do(errs chan error) error { return t.call(errs, t.name) }

type taskWithNodeGroupSpec struct {
	info      string
	skip      bool
	nodeGroup *api.NodeGroup
	call      func(chan error, *api.NodeGroup) error
}

func (t *taskWithNodeGroupSpec) Skip() bool       { return t.skip }
func (t *taskWithNodeGroupSpec) Describe() string { return info(t.skip, t.info) }
func (t *taskWithNodeGroupSpec) Do(errs chan error) error {
	return t.call(errs, t.nodeGroup)
}

type taskWithStackSpec struct {
	info  string
	skip  bool
	stack *Stack
	call  func(*Stack, chan error) error
}

func (t *taskWithStackSpec) Skip() bool       { return t.skip }
func (t *taskWithStackSpec) Describe() string { return info(t.skip, t.info) }
func (t *taskWithStackSpec) Do(errs chan error) error {
	return t.call(t.stack, errs)
}

type asyncTaskWithStackSpec struct {
	info  string
	skip  bool
	stack *Stack
	call  func(*Stack) (*Stack, error)
}

func (t *asyncTaskWithStackSpec) Skip() bool       { return t.skip }
func (t *asyncTaskWithStackSpec) Describe() string { return info(t.skip, t.info) + " [async]" }
func (t *asyncTaskWithStackSpec) Do(errs chan error) error {
	_, err := t.call(t.stack)
	close(errs)
	return err
}

func doSingleTask(allErrs chan error, task Task) bool {
	desc := task.Describe()
	logger.Debug("started task: %s", desc)
	errs := make(chan error)
	if err := task.Do(errs); err != nil {
		allErrs <- err
		return false
	}
	if err := <-errs; err != nil {
		allErrs <- err
		return false
	}
	logger.Debug("completed task: %s", desc)
	return true
}

func doParallelTasks(allErrs chan error, tasks []Task) {
	wg := &sync.WaitGroup{}
	wg.Add(len(tasks))
	for t := range tasks {
		go func(t int) {
			defer wg.Done()
			if ok := doSingleTask(allErrs, tasks[t]); !ok {
				logger.Debug("failed task: %s (will continue until other parallel tasks are completed)", tasks[t].Describe())
			}
		}(t)
	}
	logger.Debug("waiting for %d parallel tasks to complete", len(tasks))
	wg.Wait()
	close(allErrs)
}

func doSequentialTasks(allErrs chan error, tasks []Task) {
	for t := range tasks {
		if ok := doSingleTask(allErrs, tasks[t]); !ok {
			logger.Debug("failed task: %s (will not run other sequential tasks)", tasks[t].Describe())
			break
		}
	}
	close(allErrs)
}
