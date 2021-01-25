package irsa

import (
	"fmt"

	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"

	"github.com/weaveworks/eksctl/pkg/cfn/builder"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

func NewUpdateIAMServiceAccountTask(clusterName string, sa *api.ClusterIAMServiceAccount, stackManager StackManager, iamServiceAccount *api.ClusterIAMServiceAccount, oidcManager *iamoidc.OpenIDConnectManager) (*tasks.TaskTree, error) {

	rs := builder.NewIAMServiceAccountResourceSet(iamServiceAccount, oidcManager)
	err := rs.AddAllResources()
	if err != nil {
		return nil, err
	}

	template, err := rs.RenderJSON()
	if err != nil {
		return nil, err
	}

	var templateData manager.TemplateBody = template

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
	return taskTree, nil
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

	desc := fmt.Sprintf("updating policies for IAMServiceAccount %s/%s", t.sa.Namespace, t.sa.Name)
	return t.stackManager.UpdateStack(stackName, "updating-policy", desc, t.templateData, nil)
}
