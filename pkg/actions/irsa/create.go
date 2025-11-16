package irsa

import (
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

func (m *Manager) CreateIAMServiceAccount(iamServiceAccounts []*api.ClusterIAMServiceAccount, plan bool) error {
	taskTree := m.stackManager.NewTasksToCreateIAMServiceAccounts(iamServiceAccounts, m.oidcManager, kubernetes.NewCachedClientSet(m.clientSet))
	taskTree.PlanMode = plan

	err := doTasks(taskTree, actionCreate)

	logPlanModeWarning(plan && len(iamServiceAccounts) > 0)

	return err
}
