package iam

import (
	"fmt"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

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
