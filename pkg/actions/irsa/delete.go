package irsa

import (
	"context"

	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

func (m *Manager) Delete(ctx context.Context, serviceAccounts []string, plan, wait bool) error {
	taskTree, err := m.stackManager.NewTasksToDeleteIAMServiceAccounts(ctx, serviceAccounts, kubernetes.NewCachedClientSet(m.clientSet), wait)
	if err != nil {
		return err
	}
	taskTree.PlanMode = plan

	doTasks(taskTree, "delete")

	logPlanModeWarning(plan && taskTree.Len() > 0)
	return nil
}
