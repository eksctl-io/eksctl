package nodegroup

import (
	"context"
	"fmt"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"

	"github.com/kris-nova/logger"
)

func (m *Manager) Delete(ctx context.Context, nodeGroups []*api.NodeGroup, managedNodeGroups []*api.ManagedNodeGroup, wait, plan bool) error {
	var nodeGroupsWithStacks []eks.KubeNodeGroup

	for _, n := range nodeGroups {
		nodeGroupsWithStacks = append(nodeGroupsWithStacks, n)
	}

	tasks := &tasks.TaskTree{Parallel: true}
	stacks, err := m.stackManager.ListNodeGroupStacksWithStatuses(ctx)
	if err != nil {
		return err
	}

	for _, n := range managedNodeGroups {
		if m.hasStacks(stacks, n.Name) != nil {
			nodeGroupsWithStacks = append(nodeGroupsWithStacks, n)
		} else {
			tasks.Append(m.stackManager.NewTaskToDeleteUnownedNodeGroup(ctx, m.cfg.Metadata.Name, n.Name, m.ctl.AWSProvider.EKS(), nil))
		}
	}

	shouldDelete := func(ngName string) bool {
		for _, n := range nodeGroupsWithStacks {
			if n.NameString() == ngName {
				return true
			}
		}
		return false
	}

	deleteTasks, err := m.stackManager.NewTasksToDeleteNodeGroups(stacks, shouldDelete, wait, nil)
	if err != nil {
		return err
	}
	tasks.Append(deleteTasks)

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
