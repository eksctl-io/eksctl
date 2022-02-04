package irsa

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
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

		stack, ok := existingIAMStacksMap[stackName]
		if !ok {
			logger.Info("cannot update IAMServiceAccount %s/%s as it does not exist", iamServiceAccount.Namespace, iamServiceAccount.Name)
			nonExistingSAs = append(nonExistingSAs, fmt.Sprintf("%s/%s", iamServiceAccount.Namespace, iamServiceAccount.Name))
			continue
		}

		roleName, err := getRoleNameFromStack(stack)
		if err != nil {
			return err
		}

		iamServiceAccount.RoleName = roleName
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

func getRoleNameFromStack(stack *manager.Stack) (string, error) {
	var roleName string
	for _, o := range stack.Outputs {
		if aws.StringValue(o.OutputKey) != outputs.IAMServiceAccountRoleName {
			continue
		}
		roleArn, err := arn.Parse(aws.StringValue(o.OutputValue))
		if err != nil {
			return "", fmt.Errorf("failed to parse role arn %q: %w", aws.StringValue(o.OutputValue), err)
		}
		start := strings.IndexRune(roleArn.Resource, '/')
		if start == -1 {
			return "", fmt.Errorf("failed to parse resource: %s", roleArn.Resource)
		}
		roleName = roleArn.Resource[start+1:]
		logger.Info("found role name %q for service account role", roleName)
		break
	}
	if roleName == "" {
		return "", fmt.Errorf("failed to find role name service account")
	}
	return roleName, nil
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
