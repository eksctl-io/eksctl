package irsa

import (
	"fmt"

	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"

	"github.com/weaveworks/eksctl/pkg/cfn/builder"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

func NewUpdateIAMServiceAccountTaskWithStack(clusterName string, sa *api.ClusterIAMServiceAccount, stack *manager.Stack, stackManager manager.StackManager, oidcManager *iamoidc.OpenIDConnectManager) (*tasks.TaskTree, error) {
	taskTree := &tasks.TaskTree{Parallel: false}
	templateData, err := buildTemplate(sa, oidcManager)
	if err != nil {
		return nil, err
	}

	taskTree.Append(
		&updateIAMServiceAccountTask{
			info:         fmt.Sprintf("update IAMServiceAccount %s/%s", sa.Namespace, sa.Name),
			stackManager: stackManager,
			templateData: templateData,
			sa:           sa,
			clusterName:  clusterName,
			stack:        stack,
		},
	)
	return taskTree, nil

}
func NewUpdateIAMServiceAccountTask(clusterName string, sa *api.ClusterIAMServiceAccount, stackManager manager.StackManager, oidcManager *iamoidc.OpenIDConnectManager) (*tasks.TaskTree, error) {
	taskTree := &tasks.TaskTree{Parallel: false}
	templateData, err := buildTemplate(sa, oidcManager)
	if err != nil {
		return nil, err
	}

	taskTree.Append(
		&updateIAMServiceAccountTask{
			info:         fmt.Sprintf("update IAMServiceAccount %s/%s", sa.Namespace, sa.Name),
			stackManager: stackManager,
			templateData: templateData,
			sa:           sa,
			clusterName:  clusterName,
		},
	)
	return taskTree, nil
}

func buildTemplate(sa *api.ClusterIAMServiceAccount, oidcManager *iamoidc.OpenIDConnectManager) (manager.TemplateBody, error) {
	rs := builder.NewIAMRoleResourceSetForServiceAccount(sa, oidcManager)
	err := rs.AddAllResources()
	if err != nil {
		return nil, err
	}

	template, err := rs.RenderJSON()
	if err != nil {
		return nil, err
	}

	var templateData manager.TemplateBody = template
	return templateData, nil
}

type updateIAMServiceAccountTask struct {
	sa           *api.ClusterIAMServiceAccount
	stackManager manager.StackManager
	templateData manager.TemplateData
	clusterName  string
	info         string
	stack        *manager.Stack
}

func (t *updateIAMServiceAccountTask) Describe() string { return t.info }

func (t *updateIAMServiceAccountTask) Do(errorCh chan error) error {
	go func() {
		errorCh <- nil
	}()

	desc := fmt.Sprintf("updating policies for IAMServiceAccount %s/%s", t.sa.Namespace, t.sa.Name)
	if t.stack == nil {
		stackName := makeIAMServiceAccountStackName(t.clusterName, t.sa.Namespace, t.sa.Name)
		return t.stackManager.UpdateStack(stackName, "updating-policy", desc, t.templateData, nil)
	}
	return t.stackManager.UpdateCachedStack(t.stack, "updating-poicy", desc, t.templateData, nil)
}
