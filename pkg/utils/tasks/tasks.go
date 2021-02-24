package tasks

import (
	"fmt"
	"strings"
	"sync"

	"github.com/kris-nova/logger"
)

// Task is a common interface for the stack manager tasks
type Task interface {
	Describe() string
	Do(chan error) error
}

type GenericTask struct {
	Description string
	Doer        func() error
}

func (t *GenericTask) Describe() string {
	return t.Description
}
func (t *GenericTask) Do(errCh chan error) error {
	close(errCh)
	return t.Doer()
}

type SynchronousTaskIface interface {
	Describe() string
	Do() error
}

type SynchronousTask struct {
	SynchronousTaskIface
}

func (t SynchronousTask) Do(errCh chan error) error {
	defer close(errCh)
	return t.SynchronousTaskIface.Do()
}

// TaskTree wraps a set of tasks
type TaskTree struct {
	Tasks     []Task
	Parallel  bool
	PlanMode  bool
	IsSubTask bool
}

// Append new tasks to the set
func (t *TaskTree) Append(newTasks ...Task) {
	t.Tasks = append(t.Tasks, newTasks...)
}

// Len returns number of tasks in the set
func (t *TaskTree) Len() int {
	if t == nil {
		return 0
	}
	return len(t.Tasks)
}

// Describe the set
func (t *TaskTree) Describe() string {
	descriptions := []string{}
	for _, task := range t.Tasks {
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

// Do will run through the set in the background, it may return an error immediately,
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
		go doParallelTasks(errs, t.Tasks)
	} else {
		go doSequentialTasks(errs, t.Tasks)
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
		go doParallelTasks(errs, t.Tasks)
	} else {
		go doSequentialTasks(errs, t.Tasks)
	}

	allErrs := []error{}
	for err := range errs {
		allErrs = append(allErrs, err)
	}
	return allErrs
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

type TaskWithoutParams struct {
	Info string
	Call func(chan error) error
}

func (t *TaskWithoutParams) Describe() string         { return t.Info }
func (t *TaskWithoutParams) Do(errs chan error) error { return t.Call(errs) }

type TaskWithNameParam struct {
	Info string
	Name string
	Call func(chan error, string) error
}

func (t *TaskWithNameParam) Describe() string         { return t.Info }
func (t *TaskWithNameParam) Do(errs chan error) error { return t.Call(errs, t.Name) }
