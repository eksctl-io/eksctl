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
	DeleteTasks(serviceAccounts []string) (*tasks.TaskTree, error)
	UpdateTasks(iamServiceAccounts []*api.ClusterIAMServiceAccount) (*tasks.TaskTree, error)
	IsUpToDate(iamServiceAccount *api.ClusterIAMServiceAccount, stack *manager.Stack) (bool, error)
}

func (r *Reconciler) ReconcileIAMServiceAccounts() (*tasks.TaskTree, *tasks.TaskTree, *tasks.TaskTree, error) {
	stacks, err := r.stackManager.DescribeIAMServiceAccountStacks()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to discover existing service accounts: %w", err)
	}

	existingServiceAccountsNameToStackMap := listOfStacksToServiceAccountMap(stacks)
	desiredServiceAccountsNameMapSpec := listOfServiceAccountsToMap(r.cfg.IAM.ServiceAccounts)

	var toDelete []string
	var toCreate, toUpdate []*api.ClusterIAMServiceAccount

	for saName, saSpec := range desiredServiceAccountsNameMapSpec {
		if stack, ok := existingServiceAccountsNameToStackMap[saName]; ok {
			isUpToDate, err := r.irsaManager.IsUpToDate(saSpec, stack)
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

	for saName, _ := range existingServiceAccountsNameToStackMap {
		if _, ok := desiredServiceAccountsNameMapSpec[saName]; !ok {
			toDelete = append(toDelete, saName)
		}
	}

	createTasks := r.irsaManager.CreateTasks(toCreate)
	updateTasks, err := r.irsaManager.UpdateTasks(toUpdate)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate update tasks: %w", err)
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
