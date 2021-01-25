package iam

import (
	"fmt"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"

	"github.com/kris-nova/logger"

	"github.com/weaveworks/eksctl/pkg/utils/tasks"

	"strings"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

func (a *Manager) UpdateIAMServiceAccounts(iamServiceAccounts []*api.ClusterIAMServiceAccount, plan bool) error {
	var nonExistingSAs []string
	updateTasks := &tasks.TaskTree{Parallel: true}

	existingIAMStacks, err := a.stackManager.ListStacksMatching("eksctl-.*-addon-iamserviceaccount")
	if err != nil {
		return err
	}

	existingIAMStacksMap := listToSet(existingIAMStacks)

	for _, iamServiceAccount := range iamServiceAccounts {
		stackName := makeIAMServiceAccountStackName(a.clusterName, iamServiceAccount.Namespace, iamServiceAccount.Name)

		if _, ok := existingIAMStacksMap[stackName]; !ok {
			logger.Info("cannot update IAMServiceAccount %s/%s as it does not exist", iamServiceAccount.Namespace, iamServiceAccount.Name)
			nonExistingSAs = append(nonExistingSAs, fmt.Sprintf("%s/%s", iamServiceAccount.Namespace, iamServiceAccount.Name))
			continue
		}

		taskTree, err := NewUpdateIAMServiceAccountTask(a.clusterName, iamServiceAccount, a.stackManager, iamServiceAccount, a.oidcManager)
		if err != nil {
			return err
		}
		taskTree.PlanMode = plan
		updateTasks.Append(taskTree)
	}
	if len(nonExistingSAs) > 0 {
		logger.Info("the following IAMServiceAccounts will not be updated as they do not exist: %v", strings.Join(nonExistingSAs, ", "))
	}

	defer logPlanModeWarning(plan && len(iamServiceAccounts) > 0)
	return doTasks(updateTasks)

}

func listToSet(stacks []*manager.Stack) map[string]struct{} {
	stacksMap := make(map[string]struct{})
	for _, stack := range stacks {
		stacksMap[*stack.StackName] = struct{}{}
	}
	return stacksMap
}
func makeIAMServiceAccountStackName(clusterName, namespace, name string) string {
	return fmt.Sprintf("eksctl-%s-addon-iamserviceaccount-%s-%s", clusterName, namespace, name)
}
