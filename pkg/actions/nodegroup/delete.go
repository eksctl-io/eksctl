package nodegroup

import (
	"fmt"

	"github.com/kris-nova/logger"
	"github.com/weaveworks/eksctl/pkg/eks"
)

func (ng *NodeGroup) Delete(allNodeGroups []eks.KubeNodeGroup, wait, plan bool) error {
	shouldDelete := func(ngName string) bool {
		for _, n := range allNodeGroups {
			if n.NameString() == ngName {
				return true
			}
		}
		return false
	}

	tasks, err := ng.manager.NewTasksToDeleteNodeGroups(shouldDelete, wait, nil)
	if err != nil {
		return err
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
