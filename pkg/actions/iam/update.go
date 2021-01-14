package iam

import (
	"fmt"

	"github.com/kris-nova/logger"

	"github.com/weaveworks/eksctl/pkg/utils/tasks"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

func (m *Manager) UpdateIAMServiceAccounts(iamServiceAccounts []*api.ClusterIAMServiceAccount, plan bool) error {
	var missingIAMServiceAccounts []string
	updateTasks := &tasks.TaskTree{Parallel: true}

	for _, iamServiceAccount := range iamServiceAccounts {
		stackName := makeIAMServiceAccountStackName(m.clusterName, iamServiceAccount.Namespace, iamServiceAccount.Name)
		stacks, err := m.stackManager.ListStacksMatching(stackName)
		if err != nil {
			return err
		}

		if len(stacks) == 0 {
			logger.Info("Cannot update IAMServiceAccount %s/%s as it does not exist", iamServiceAccount.Namespace, iamServiceAccount.Name)
			missingIAMServiceAccounts = append(missingIAMServiceAccounts, fmt.Sprintf("%s/%s", iamServiceAccount.Namespace, iamServiceAccount.Name))
			continue
		}

		rs := builder.NewIAMServiceAccountResourceSet(iamServiceAccount, m.oidcManager)
		err = rs.AddAllResources()
		if err != nil {
			return err
		}

		template, err := rs.RenderJSON()
		if err != nil {
			return err
		}

		var templateBody manager.TemplateBody = template
		taskTree := UpdateIAMServiceAccountTask(m.clusterName, iamServiceAccount, m.stackManager, templateBody)
		taskTree.PlanMode = plan
		updateTasks.Append(taskTree)
	}
	if len(missingIAMServiceAccounts) > 0 {
		logger.Info("the following IAMServiceAccounts will not be updated as they do not exist: %v", missingIAMServiceAccounts)
	}

	err := doTasks(updateTasks)
	LogPlanModeWarning(plan && len(iamServiceAccounts) > 0)
	return err

}

func UpdateIAMServiceAccountTask(clusterName string, sa *api.ClusterIAMServiceAccount, stackManager StackManager, templateData manager.TemplateData) *tasks.TaskTree {
	taskTree := &tasks.TaskTree{Parallel: false}

	taskTree.Append(
		&updateIAMServiceAccountTask{
			info:         fmt.Sprintf("update IAMServiceAccount %s/%s", sa.Namespace, sa.Name),
			stackManager: stackManager,
			templateData: templateData,
			sa:           sa,
			clusterName:  clusterName,
		},
	)
	return taskTree
}

type updateIAMServiceAccountTask struct {
	sa           *api.ClusterIAMServiceAccount
	stackManager StackManager
	templateData manager.TemplateData
	clusterName  string
	info         string
}

func (t *updateIAMServiceAccountTask) Describe() string { return t.info }

func (t *updateIAMServiceAccountTask) Do(errorCh chan error) error {
	stackName := makeIAMServiceAccountStackName(t.clusterName, t.sa.Namespace, t.sa.Name)
	go func() {
		errorCh <- nil
	}()
	return t.stackManager.UpdateStack(stackName, "updating-policy", "updating policies", t.templateData, nil)

}

func makeIAMServiceAccountStackName(clusterName, namespace, name string) string {
	return fmt.Sprintf("eksctl-%s-addon-iamserviceaccount-%s-%s", clusterName, namespace, name)
}
