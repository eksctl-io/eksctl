package iam

import (
	"fmt"

	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"

	kubewrapper "github.com/weaveworks/eksctl/pkg/kubernetes"
)

func (m *Manager) CreateIAMServiceAccount(iamServiceAccounts []*api.ClusterIAMServiceAccount, plan bool) error {
	taskTree := m.NewTasksToCreateIAMServiceAccounts(iamServiceAccounts, kubernetes.NewCachedClientSet(m.clientSet))
	taskTree.PlanMode = plan

	err := doTasks(taskTree)

	LogPlanModeWarning(plan && len(iamServiceAccounts) > 0)

	return err
}

func (m *Manager) NewTasksToCreateIAMServiceAccounts(serviceAccounts []*api.ClusterIAMServiceAccount, clientSetGetter kubernetes.ClientSetGetter) *tasks.TaskTree {
	taskTree := &tasks.TaskTree{Parallel: true}

	for i := range serviceAccounts {
		sa := serviceAccounts[i]
		saTasks := &tasks.TaskTree{
			Parallel:  false,
			IsSubTask: true,
		}

		saTasks.Append(&createIAMServiceAccountTask{
			info:           fmt.Sprintf("create IAM role for serviceaccount %q", sa.NameString()),
			serviceAccount: sa,
			oidc:           m.oidcManager,
			stackManager:   m.stackManager,
		})

		saTasks.Append(&kubernetesTask{
			info:       fmt.Sprintf("create serviceaccount %q", sa.NameString()),
			kubernetes: clientSetGetter,
			call: func(clientSet kubernetes.Interface) error {
				sa.SetAnnotations()
				if err := kubernetes.MaybeCreateServiceAccountOrUpdateMetadata(clientSet, sa.ClusterIAMMeta.AsObjectMeta()); err != nil {
					return errors.Wrapf(err, "failed to create service account %s", sa.NameString())
				}
				return nil
			},
		})

		taskTree.Append(saTasks)
	}
	return taskTree
}

type createIAMServiceAccountTask struct {
	info           string
	serviceAccount *api.ClusterIAMServiceAccount
	oidc           *iamoidc.OpenIDConnectManager
	stackManager   StackManager
}

func (t *createIAMServiceAccountTask) Describe() string { return t.info }
func (t *createIAMServiceAccountTask) Do(errs chan error) error {
	return t.stackManager.CreateIAMServiceAccount(errs, t.serviceAccount, t.oidc)
}

type kubernetesTask struct {
	info       string
	kubernetes kubewrapper.ClientSetGetter
	call       func(kubernetes.Interface) error
}

func (t *kubernetesTask) Describe() string { return t.info }
func (t *kubernetesTask) Do(errs chan error) error {
	if t.kubernetes == nil {
		return fmt.Errorf("cannot start task %q as Kubernetes client configurtaion wasn't provided", t.Describe())
	}
	clientSet, err := t.kubernetes.ClientSet()
	if err != nil {
		return err
	}
	err = t.call(clientSet)
	close(errs)
	return err
}
