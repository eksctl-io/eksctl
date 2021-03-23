package irsa

import (
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

func (a *Manager) CreateIAMServiceAccount(iamServiceAccounts []*api.ClusterIAMServiceAccount, plan bool) error {
	taskTree := a.CreateTasks(iamServiceAccounts)
	taskTree.PlanMode = plan

	err := doTasks(taskTree)

	logPlanModeWarning(plan && len(iamServiceAccounts) > 0)

	return err
}

func (a *Manager) CreateTasks(iamServiceAccounts []*api.ClusterIAMServiceAccount) *tasks.TaskTree {
	return a.stackManager.NewTasksToCreateIAMServiceAccounts(iamServiceAccounts, a.oidcManager, kubernetes.NewCachedClientSet(a.clientSet))
}
