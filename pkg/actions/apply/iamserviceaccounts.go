package apply

import (
	"fmt"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

//go:generate counterfeiter -o fakes/fake_irsa_manager.go . IRSAManager
type IRSAManager interface {
	CreateTasks(iamServiceAccounts []*api.ClusterIAMServiceAccount) *tasks.TaskTree
	DeleteTasks(serviceAccounts map[string]*manager.Stack) (*tasks.TaskTree, error)
	UpdateTask(iamServiceAccount *api.ClusterIAMServiceAccount, stack *manager.Stack) (*tasks.TaskTree, error)
	IsUpToDate(iamServiceAccount api.ClusterIAMServiceAccount, stack *manager.Stack) (bool, error)
}

func (r *Reconciler) ReconcileIAMServiceAccounts() (*tasks.TaskTree, *tasks.TaskTree, *tasks.TaskTree, error) {
	stacks, err := r.stackManager.DescribeIAMServiceAccountStacks()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to discover existing service accounts: %w", err)
	}

	//map[serviceAccountName string]*manager.Stack
	existingServiceAccountsNameToStackMap := listOfStacksToServiceAccountMap(stacks)
	//map[serviceAccountName string]*api.ClusterIAMServiceAccount
	desiredServiceAccountsNameMapSpec := listOfServiceAccountsToMap(r.cfg.IAM.ServiceAccounts)

	toDelete := make(map[string]*manager.Stack)
	var toCreate, toUpdate []*api.ClusterIAMServiceAccount

	for saName, saSpec := range desiredServiceAccountsNameMapSpec {
		if stack, ok := existingServiceAccountsNameToStackMap[saName]; ok {
			isUpToDate, err := r.irsaManager.IsUpToDate(*saSpec, stack)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("failed to check if service account is up to date: %w", err)
			}
			if !isUpToDate {
				toUpdate = append(toUpdate, saSpec)
			} else {
				logger.Debug("IAMServiceAccount %s is already up to date", saName)
			}
		} else {
			toCreate = append(toCreate, saSpec)
		}
	}

	for saName, saStack := range existingServiceAccountsNameToStackMap {
		if _, ok := desiredServiceAccountsNameMapSpec[saName]; !ok {
			toDelete[saName] = saStack
		}
	}

	createTasks := r.irsaManager.CreateTasks(toCreate)
	updateTasks := &tasks.TaskTree{Parallel: false}
	for _, sa := range toUpdate {
		task, err := r.irsaManager.UpdateTask(sa, existingServiceAccountsNameToStackMap[sa.NameString()])
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to generate update tasks: %w", err)
		}
		updateTasks.Append(task)
	}
	deleteTasks, err := r.irsaManager.DeleteTasks(toDelete)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate delete tasks: %w", err)
	}

	return createTasks, updateTasks, deleteTasks, nil
}

func listOfStacksToServiceAccountMap(stacks []*manager.Stack) map[string]*manager.Stack {
	m := make(map[string]*manager.Stack)
	for _, stack := range stacks {
		m[manager.GetIAMServiceAccountName(stack)] = stack
	}
	return m
}

func listOfServiceAccountsToMap(serviceAccounts []*api.ClusterIAMServiceAccount) map[string]*api.ClusterIAMServiceAccount {
	m := make(map[string]*api.ClusterIAMServiceAccount)
	for _, sa := range serviceAccounts {
		m[sa.NameString()] = sa
	}
	return m
}
