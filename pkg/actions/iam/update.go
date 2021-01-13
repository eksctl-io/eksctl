package iam

import (
	"fmt"

	"github.com/weaveworks/eksctl/pkg/utils/tasks"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

func (a *Manager) UpdateIAMServiceAccounts(iamServiceAccounts []*api.ClusterIAMServiceAccount, plan bool) error {
	updateTasks := &tasks.TaskTree{Parallel: true}
	for _, iamServiceAccount := range iamServiceAccounts {
		stackName := makeIAMServiceAccountStackName(a.clusterName, iamServiceAccount.Namespace, iamServiceAccount.Name)
		stacks, err := a.stackManager.ListStacksMatching(stackName)
		if err != nil {
			return err
		}

		if len(stacks) == 0 {
			return fmt.Errorf("IAMServiceAccount %s/%s does not exist", iamServiceAccount.Namespace, iamServiceAccount.Name)
		}

		rs := builder.NewIAMServiceAccountResourceSet(iamServiceAccount, a.oidcManager)
		err = rs.AddAllResources()
		if err != nil {
			return err
		}

		template, err := rs.RenderJSON()
		if err != nil {
			return err
		}

		var templateBody manager.TemplateBody = template
		taskTree := UpdateIAMServiceAccountTask(a.clusterName, iamServiceAccount, a.stackManager, templateBody)
		taskTree.PlanMode = plan
		updateTasks.Append(taskTree)
	}

	return doTasks(updateTasks)
}

func makeIAMServiceAccountStackName(clusterName, namespace, name string) string {
	return fmt.Sprintf("eksctl-%s-addon-iamserviceaccount-%s-%s", clusterName, namespace, name)
}
