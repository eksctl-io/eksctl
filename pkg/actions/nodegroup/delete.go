package nodegroup

import (
	"fmt"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"

	"github.com/kris-nova/logger"
)

func (m *Manager) Delete(nodeGroups []*api.NodeGroup, managedNodeGroups []*api.ManagedNodeGroup, wait, plan bool) error {
	tasks, err := m.stackManager.NewTasksToDeleteNodeGroups(wait, nil)
	if err != nil {
		return err
	}

	for _, n := range managedNodeGroups {
		hasStacks, err := m.hasStacks(n.Name)
		if err != nil {
			return err
		}

		if !hasStacks {
			tasks.Append(m.stackManager.NewTaskToDeleteUnownedNodeGroup(m.cfg.Metadata.Name, n.Name, m.ctl.Provider.EKS(), nil))
		}
	}

	tasks.PlanMode = plan
	logger.Info(tasks.Describe())
	if errs := tasks.DoAllSync(); len(errs) > 0 {
		return handleErrors(errs, "nodegroup(s)")
	}
	return nil
}

func handleErrors(errs []error, subject string) error {
	logger.Info("%d error(s) occurred while deleting %s", len(errs), subject)
	for _, err := range errs {
		logger.Critical("%s\n", err.Error())
	}
	return fmt.Errorf("failed to delete %s", subject)
}
