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

// Describe collects all tasks which have been added to the task tree.
// This is a lazy tree which does not track its nodes in any form. This function
// is recursively called from the rest of the task Describes and eventually
// returns a collection of all the tasks' `Info` value.
func (t *TaskTree) Describe() string {
	if t.Len() == 0 {
		return "no tasks"
	}
	var descriptions []string
	for _, task := range t.Tasks {
		descriptions = append(descriptions, strings.TrimSuffix(task.Describe(), "\n"))
	}
	noun := "task"
	if t.IsSubTask {
		noun = "sub-task"
	}
	if len(descriptions) == 1 {
		msg := fmt.Sprintf("1 %s: { %s }", noun, descriptions[0])
		if t.IsSubTask {
			msg = descriptions[0]
		}
		return msg
	}
	count := len(descriptions)
	mode := "sequential"
	if t.Parallel {
		mode = "parallel"
	}
	noun += "s"
	head := fmt.Sprintf("\n%d %s %s: { ", count, mode, noun)
	var tail string
	if t.IsSubTask {
		// Only add a linebreak at the end if we have multiple subtasks as well. Otherwise, leave it
		// as single line.
		head = fmt.Sprintf("\n%s%d %s %s: { ", strings.Repeat(" ", 4), count, mode, noun)
		tail = "\n"
		for _, d := range descriptions {
			// all tasks are sub-tasks if they are inside a task.
			// which means we don't have to care about sequential tasks.
			if strings.Contains(d, "sub-task") {
				// trim the previous leading tail new line...
				d = strings.TrimPrefix(d, "\n")
				split := strings.Split(d, "\n")
				// indent all lines of the subtask one deepness more
				var result []string
				for _, s := range split {
					result = append(result, strings.Repeat(" ", 4)+s)
				}
				// join it back up with line breaks
				d = strings.Join(result, "\n")
			} else {
				d = strings.Repeat(" ", 8) + d
			}
			tail += fmt.Sprintf("%s,\n", d)
		}
		// closing the final bracket
		tail += fmt.Sprintf("%s}", strings.Repeat(" ", 4))
	} else {
		// if it isn't a subtask, we just add the descriptions as is joined by new line.
		// this results in line like `1 task: { t1.1 }` which are more readable this way.
		tail = fmt.Sprintf("%s \n}", strings.Join(descriptions, ", "))
	}
	msg := head + tail
	if t.PlanMode {
		msg = "(plan) " + msg
	}
	return msg + "\n"
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
