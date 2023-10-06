package irsa

import (
	"context"
	"fmt"

	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/google/uuid"
	"github.com/kris-nova/logger"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

func NewUpdateIAMServiceAccountTask(clusterName string, sa *api.ClusterIAMServiceAccount, stackManager StackManager, oidcManager *iamoidc.OpenIDConnectManager) (*tasks.TaskTree, error) {
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

type createIAMRoleForServiceAccountTask struct {
	ctx          context.Context
	info         string
	clusterName  string
	region       string
	stackManager StackManager
	sa           *api.ClusterIAMServiceAccount
	oidc         *iamoidc.OpenIDConnectManager
}

func (t *createIAMRoleForServiceAccountTask) Describe() string { return t.info }
func (t *createIAMRoleForServiceAccountTask) Do(errs chan error) error {
	defer close(errs)

	name := makeIAMServiceAccountStackName(t.clusterName, t.sa.Namespace, t.sa.Name)
	logger.Info("building iamserviceaccount stack %q", name)
	stack := builder.NewIAMRoleResourceSetForServiceAccount(t.sa, t.oidc)
	if err := stack.AddAllResources(); err != nil {
		return err
	}

	if t.sa.Tags == nil {
		t.sa.Tags = make(map[string]string)
	}
	t.sa.Tags[api.IAMServiceAccountNameTag] = t.sa.NameString()

	if err := t.stackManager.CreateStack(t.ctx, name, stack, t.sa.Tags, nil, errs); err != nil {
		logger.Info("an error occurred creating the stack, to cleanup resources, run 'eksctl delete iamserviceaccount --region=%s --name=%s --namespace=%s'", t.region, t.sa.Name, t.sa.Namespace)
		return err
	}
	return nil
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
	return t.stackManager.UpdateStack(context.TODO(), manager.UpdateStackOptions{
		StackName:     stackName,
		ChangeSetName: fmt.Sprintf("updating-policy-%s", uuid.NewString()),
		Description:   desc,
		TemplateData:  t.templateData,
		Wait:          true,
	})
}

type deleteIAMRoleForServiceAccountTask struct {
	ctx          context.Context
	info         string
	stack        *cfntypes.Stack
	stackManager StackManager
	wait         bool
}

func (t *deleteIAMRoleForServiceAccountTask) Describe() string { return t.info }

func (t *deleteIAMRoleForServiceAccountTask) Do(errorCh chan error) error {
	defer close(errorCh)

	errMsg := fmt.Sprintf("deleting IAM role for serviceaccount %q", *t.stack.StackName)
	if t.wait {
		if err := t.stackManager.DeleteStackBySpecSync(t.ctx, t.stack, errorCh); err != nil {
			return fmt.Errorf("%s: %w", errMsg, err)
		}
		return nil
	}
	if _, err := t.stackManager.DeleteStackBySpec(t.ctx, t.stack); err != nil {
		return fmt.Errorf("%s: %w", errMsg, err)
	}
	return nil
}

type kubernetesTask struct {
	info       string
	kubernetes kubernetes.ClientSetGetter
	objectMeta v1.ObjectMeta
	call       func(kubernetes.Interface, v1.ObjectMeta) error
}

func (t *kubernetesTask) Describe() string { return t.info }
func (t *kubernetesTask) Do(errs chan error) error {
	defer close(errs)

	if t.kubernetes == nil {
		return fmt.Errorf("cannot start task %q as Kubernetes client configurtaion wasn't provided", t.Describe())
	}
	clientSet, err := t.kubernetes.ClientSet()
	if err != nil {
		return err
	}
	err = t.call(clientSet, t.objectMeta)
	return err
}
