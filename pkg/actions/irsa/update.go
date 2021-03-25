package irsa

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/weaveworks/eksctl/pkg/cfn/builder"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"

	"github.com/kris-nova/logger"

	"github.com/weaveworks/eksctl/pkg/utils/tasks"

	"strings"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (a *Manager) IsUpToDate(sa api.ClusterIAMServiceAccount, stack *manager.Stack) (bool, error) {
	rs := builder.NewIAMRoleResourceSetForServiceAccount(&sa, a.oidcManager)
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

	var originalJSON, newJSON interface{}
	err = json.Unmarshal(template, &originalJSON)
	if err != nil {
		logger.Info(string(template))
		return false, err
	}

	err = json.Unmarshal([]byte(existingTemplate), &newJSON)
	if err != nil {
		return false, err
	}

	templateUpToDate := reflect.DeepEqual(originalJSON, newJSON)
	if !templateUpToDate {
		return false, nil
	}

	existingSA, err := a.clientSet.CoreV1().ServiceAccounts(sa.Namespace).Get(context.TODO(), sa.Name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}

	labelsUpToDate := reflect.DeepEqual(existingSA.Labels, sa.Labels)

	if !labelsUpToDate {
		return false, nil
	}

	var roleARN string
	for _, output := range stack.Outputs {
		if *output.OutputKey == "Role1" {
			roleARN = *output.OutputValue
			break
		}
	}

	//the EKSRoleARN annotation is not set on the provided spec, so we have to set it ourselves in order to compare annotations
	if sa.Annotations == nil {
		sa.Annotations = make(map[string]string)
	}
	sa.Annotations[api.AnnotationEKSRoleARN] = roleARN

	return reflect.DeepEqual(sa.Annotations, existingSA.Annotations), nil
}

func (a *Manager) UpdateRolePoliciesForIAMServiceAccounts(iamServiceAccounts []*api.ClusterIAMServiceAccount, plan bool) error {
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
	updateTasks := &tasks.TaskTree{Parallel: true}

	existingIAMStacks, err := a.stackManager.ListStacksMatching("eksctl-.*-addon-iamserviceaccount")
	if err != nil {
		return nil, err
	}

	existingIAMStacksMap := listToSet(existingIAMStacks)

	for _, iamServiceAccount := range iamServiceAccounts {
		stackName := makeIAMServiceAccountStackName(a.clusterName, iamServiceAccount.Namespace, iamServiceAccount.Name)

		if _, ok := existingIAMStacksMap[stackName]; !ok {
			logger.Debug("cannot update IAMServiceAccount %s/%s as it does not exist", iamServiceAccount.Namespace, iamServiceAccount.Name)
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
