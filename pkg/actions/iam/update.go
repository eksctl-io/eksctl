package iam

import (
	"fmt"

	"github.com/kris-nova/logger"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"

	"github.com/weaveworks/eksctl/pkg/utils/tasks"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

func (a *Manager) UpdateIAMServiceAccounts(iamServiceAccounts []*api.ClusterIAMServiceAccount, plan bool) error {
	var nonExistingSAs []string
	updateTasks := &tasks.TaskTree{Parallel: true}

	for _, iamServiceAccount := range iamServiceAccounts {
		stackName := makeIAMServiceAccountStackName(a.clusterName, iamServiceAccount.Namespace, iamServiceAccount.Name)
		stacks, err := a.stackManager.ListStacksMatching(stackName)
		if err != nil {
			return err
		}

		if len(stacks) == 0 {
			logger.Info("Cannot update IAMServiceAccount %s/%s as it does not exist", iamServiceAccount.Namespace, iamServiceAccount.Name)
			nonExistingSAs = append(nonExistingSAs, fmt.Sprintf("%s/%s", iamServiceAccount.Namespace, iamServiceAccount.Name))
			continue
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
	if len(nonExistingSAs) > 0 {
		logger.Info("the following IAMServiceAccounts will not be updated as they do not exist: %v", nonExistingSAs)
	}

	err := doTasks(updateTasks)
	cmdutils.LogPlanModeWarning(plan && len(iamServiceAccounts) > 0)
	return err

}

func makeIAMServiceAccountStackName(clusterName, namespace, name string) string {
	return fmt.Sprintf("eksctl-%s-addon-iamserviceaccount-%s-%s", clusterName, namespace, name)
}
