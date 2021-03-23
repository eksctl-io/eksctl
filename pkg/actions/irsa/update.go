package irsa

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/weaveworks/eksctl/pkg/cfn/builder"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"

	"github.com/kris-nova/logger"

	"github.com/weaveworks/eksctl/pkg/utils/tasks"

	"strings"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

func (a *Manager) IsUpToDate(sa *api.ClusterIAMServiceAccount, stack *manager.Stack) (bool, error) {
	rs := builder.NewIAMRoleResourceSetForServiceAccount(sa, a.oidcManager)
	err := rs.AddAllResources()
	if err != nil {
		return false, err
	}

	template, err := rs.RenderJSON()
	if err != nil {
		return false, err
	}

	existingTemplate, err := a.stackManager.GetStackTemplate(*stack.StackName)
	if err != nil {
		return false, err
	}

	//logger.Info("existing stack:\n%s", existingTemplate)
	//logger.Info("would be created stack:\n%s", template)

	var j, j2 interface{}
	err = json.Unmarshal(template, &j)
	if err != nil {
		logger.Info("marshal1")
		logger.Info(string(template))
		return false, err
	}

	err = json.Unmarshal([]byte(existingTemplate), &j2)
	if err != nil {
		logger.Info("marshal2")
		return false, err
	}
	return reflect.DeepEqual(j2, j), nil
}

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

		taskTree, err := NewUpdateIAMServiceAccountTask(a.clusterName, iamServiceAccount, a.stackManager, a.oidcManager)
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

func (a *Manager) UpdateTasks(iamServiceAccounts []*api.ClusterIAMServiceAccount) (*tasks.TaskTree, error) {
	var nonExistingSAs []string
	updateTasks := &tasks.TaskTree{Parallel: true}

	existingIAMStacks, err := a.stackManager.ListStacksMatching("eksctl-.*-addon-iamserviceaccount")
	if err != nil {
		return nil, err
	}

	existingIAMStacksMap := listToSet(existingIAMStacks)

	for _, iamServiceAccount := range iamServiceAccounts {
		stackName := makeIAMServiceAccountStackName(a.clusterName, iamServiceAccount.Namespace, iamServiceAccount.Name)

		if _, ok := existingIAMStacksMap[stackName]; !ok {
			logger.Info("cannot update IAMServiceAccount %s/%s as it does not exist", iamServiceAccount.Namespace, iamServiceAccount.Name)
			nonExistingSAs = append(nonExistingSAs, fmt.Sprintf("%s/%s", iamServiceAccount.Namespace, iamServiceAccount.Name))
			continue
		}

		taskTree, err := NewUpdateIAMServiceAccountTask(a.clusterName, iamServiceAccount, a.stackManager, a.oidcManager)
		if err != nil {
			return nil, err
		}
		taskTree.PlanMode = false
		updateTasks.Append(taskTree)
	}
	return updateTasks, nil

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
