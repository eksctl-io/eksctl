package manager

import (
	"fmt"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"strings"
	"sync"

	"github.com/kris-nova/logger"
	"k8s.io/client-go/kubernetes"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	kubewrapper "github.com/weaveworks/eksctl/pkg/kubernetes"
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

type taskWithoutParams struct {
	info string
	call func(chan error) error
}

func (t *taskWithoutParams) Describe() string         { return t.info }
func (t *taskWithoutParams) Do(errs chan error) error { return t.call(errs) }

type taskWithNameParam struct {
	info string
	name string
	call func(chan error, string) error
}

func (t *taskWithNameParam) Describe() string         { return t.info }
func (t *taskWithNameParam) Do(errs chan error) error { return t.call(errs, t.name) }

type createClusterTask struct {
	info                 string
	stackCollection      *StackCollection
	supportsManagedNodes bool
}

func (t *createClusterTask) Describe() string { return t.info }

func (t *createClusterTask) Do(errorCh chan error) error {
	return t.stackCollection.createClusterTask(errorCh, t.supportsManagedNodes)
}

type nodeGroupTask struct {
	info                 string
	nodeGroup            *api.NodeGroup
	supportsManagedNodes bool
	stackCollection      *StackCollection
}

func (t *nodeGroupTask) Describe() string { return t.info }
func (t *nodeGroupTask) Do(errs chan error) error {
	return t.stackCollection.createNodeGroupTask(errs, t.nodeGroup, t.supportsManagedNodes)
}

type managedNodeGroupTask struct {
	info            string
	nodeGroup       *api.ManagedNodeGroup
	stackCollection *StackCollection
}

func (t *managedNodeGroupTask) Describe() string { return t.info }

func (t *managedNodeGroupTask) Do(errorCh chan error) error {
	return t.stackCollection.createManagedNodeGroupTask(errorCh, t.nodeGroup)
}

type clusterCompatTask struct {
	info            string
	stackCollection *StackCollection
}

func (t *clusterCompatTask) Describe() string { return t.info }

func (t *clusterCompatTask) Do(errorCh chan error) error {
	defer close(errorCh)
	return t.stackCollection.FixClusterCompatibility()
}

type taskWithClusterIAMServiceAccountSpec struct {
	info           string
	serviceAccount *api.ClusterIAMServiceAccount
	oidc           *iamoidc.OpenIDConnectManager
	call           func(chan error, *api.ClusterIAMServiceAccount, *iamoidc.OpenIDConnectManager) error
}

func (t *taskWithClusterIAMServiceAccountSpec) Describe() string { return t.info }
func (t *taskWithClusterIAMServiceAccountSpec) Do(errs chan error) error {
	return t.call(errs, t.serviceAccount, t.oidc)
}

type taskWithStackSpec struct {
	info  string
	stack *Stack
	call  func(*Stack, chan error) error
}

func (t *taskWithStackSpec) Describe() string { return t.info }
func (t *taskWithStackSpec) Do(errs chan error) error {
	return t.call(t.stack, errs)
}

type asyncTaskWithStackSpec struct {
	info  string
	stack *Stack
	call  func(*Stack) (*Stack, error)
}

func (t *asyncTaskWithStackSpec) Describe() string { return t.info + " [async]" }
func (t *asyncTaskWithStackSpec) Do(errs chan error) error {
	_, err := t.call(t.stack)
	close(errs)
	return err
}

type asyncTaskWithoutParams struct {
	info string
	call func() error
}

func (t *asyncTaskWithoutParams) Describe() string { return t.info }
func (t *asyncTaskWithoutParams) Do(errs chan error) error {
	err := t.call()
	close(errs)
	return err
}

type kubernetesTask struct {
	info       string
	kubernetes kubewrapper.ClientSetGetter
	call       func(kubernetes.Interface) error
}

func (t *kubernetesTask) Describe() string { return t.info }
func (t *kubernetesTask) Do(errs chan error) error {
	if t.kubernetes == nil {
		return fmt.Errorf("cannot start task %q as Kubernetes client configuration wasn't provided", t.Describe())
	}
	clientSet, err := t.kubernetes.ClientSet()
	if err != nil {
		return err
	}
	err = t.call(clientSet)
	close(errs)
	return err
}

type deleteFromAuthConfigMapTask struct {
	info            string
	clientSet       kubernetes.Interface
	instanceRoleARN string
}

func (t *deleteFromAuthConfigMapTask) Describe() string { return t.info }
func (t *deleteFromAuthConfigMapTask) Do(errs chan error) error {
	defer close(errs)
	return authconfigmap.RemoveARNIdentity(t.clientSet, t.instanceRoleARN, false)
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
