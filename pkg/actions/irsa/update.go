package irsa

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kris-nova/logger"
	"github.com/tidwall/gjson"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

const (
	resourcesPath  = "Resources"
	propertiesPath = "Properties"
	roleNamePath   = "RoleName"
)

func (a *Manager) UpdateIAMServiceAccounts(ctx context.Context, iamServiceAccounts []*api.ClusterIAMServiceAccount, existingIAMStacks []*manager.Stack, plan bool) error {
	var nonExistingSAs []string
	updateTasks := &tasks.TaskTree{Parallel: true}

	existingIAMStacksMap := listToSet(existingIAMStacks)

	for _, iamServiceAccount := range iamServiceAccounts {
		stackName := makeIAMServiceAccountStackName(a.clusterName, iamServiceAccount.Namespace, iamServiceAccount.Name)

		stack, ok := existingIAMStacksMap[stackName]
		if !ok {
			logger.Info("cannot update IAMServiceAccount %s/%s as it does not exist", iamServiceAccount.Namespace, iamServiceAccount.Name)
			nonExistingSAs = append(nonExistingSAs, fmt.Sprintf("%s/%s", iamServiceAccount.Namespace, iamServiceAccount.Name))
			continue
		}

		roleName, err := a.getRoleNameFromStackTemplate(ctx, stack)
		if err != nil {
			return err
		}
		if roleName != "" {
			logger.Info("found set role name during creation %s for account %s", roleName, iamServiceAccount.Name)
			iamServiceAccount.RoleName = roleName
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
	return doTasks(updateTasks, actionUpdate)
}

// getRoleNameFromStackTemplate returns the role if the initial stack's template contained it.
// That means it was defined upon creation, and we need to re-use that same name.
func (a *Manager) getRoleNameFromStackTemplate(ctx context.Context, stack *manager.Stack) (string, error) {
	template, err := a.stackManager.GetStackTemplate(ctx, aws.ToString(stack.StackName))
	if err != nil {
		return "", fmt.Errorf("failed to get stack template: %w", err)
	}
	resources := gjson.Get(template, resourcesPath)
	if !resources.Get(outputs.IAMServiceAccountRoleName).Exists() {
		return "", nil
	}
	return resources.Get(outputs.IAMServiceAccountRoleName).Get(propertiesPath).Get(roleNamePath).String(), nil
}

func listToSet(stacks []*manager.Stack) map[string]*manager.Stack {
	stacksMap := make(map[string]*manager.Stack)
	for _, stack := range stacks {
		stacksMap[*stack.StackName] = stack
	}
	return stacksMap
}
func makeIAMServiceAccountStackName(clusterName, namespace, name string) string {
	return fmt.Sprintf("eksctl-%s-addon-iamserviceaccount-%s-%s", clusterName, namespace, name)
}
