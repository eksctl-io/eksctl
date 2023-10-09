package irsa

import (
	"context"
	"fmt"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

var (
	managedByKubernetesLabelKey               = "app.kubernetes.io/managed-by"
	managedByKubernetesLabelValue             = "eksctl"
	maybeCreateServiceAccountOrUpdateMetadata = kubernetes.MaybeCreateServiceAccountOrUpdateMetadata
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_create_tasks_builder.go . CreateTasksBuilder
type CreateTasksBuilder interface {
	CreateIAMServiceAccountsTasks(ctx context.Context, serviceAccounts []*api.ClusterIAMServiceAccount) *tasks.TaskTree
}

type Creator struct {
	clusterName string
	region      string

	clientSetGetter kubernetes.ClientSetGetter
	oidcManager     *iamoidc.OpenIDConnectManager
	stackManager    StackManager
}

func NewCreator(
	clusterName string,
	region string,
	clientSetGetter kubernetes.ClientSetGetter,
	oidcManager *iamoidc.OpenIDConnectManager,
	stackManager StackManager) *Creator {
	return &Creator{
		clusterName:     clusterName,
		region:          region,
		clientSetGetter: clientSetGetter,
		oidcManager:     oidcManager,
		stackManager:    stackManager,
	}
}

func (c *Creator) CreateIAMServiceAccounts(ctx context.Context, serviceAccounts []*api.ClusterIAMServiceAccount, plan bool) error {
	taskTree := c.CreateIAMServiceAccountsTasks(ctx, serviceAccounts)
	taskTree.PlanMode = plan

	err := doTasks(taskTree, actionCreate)

	logPlanModeWarning(plan && len(serviceAccounts) > 0)

	return err
}

func (c *Creator) CreateIAMServiceAccountsTasks(ctx context.Context, serviceAccounts []*api.ClusterIAMServiceAccount) *tasks.TaskTree {
	taskTree := &tasks.TaskTree{Parallel: true}

	for i := range serviceAccounts {
		sa := serviceAccounts[i]
		saTasks := &tasks.TaskTree{
			Parallel:  false,
			IsSubTask: true,
		}

		if sa.AttachRoleARN == "" {
			saTasks.Append(&createIAMRoleForServiceAccountTask{
				ctx:          ctx,
				info:         fmt.Sprintf("create IAM role for serviceaccount %q", sa.NameString()),
				clusterName:  c.clusterName,
				region:       c.region,
				stackManager: c.stackManager,
				sa:           sa,
				oidc:         c.oidcManager,
			})
		} else {
			logger.Debug("attachRoleARN was provided, skipping role creation")
			sa.Status = &api.ClusterIAMServiceAccountStatus{
				RoleARN: &sa.AttachRoleARN,
			}
		}

		if sa.Labels == nil {
			sa.Labels = make(map[string]string)
		}
		sa.Labels[managedByKubernetesLabelKey] = managedByKubernetesLabelValue
		if !api.IsEnabled(sa.RoleOnly) {
			saTasks.Append(&kubernetesTask{
				info:       fmt.Sprintf("create serviceaccount %q", sa.NameString()),
				kubernetes: c.clientSetGetter,
				objectMeta: sa.ClusterIAMMeta.AsObjectMeta(),
				call: func(clientSet kubernetes.Interface, objectMeta v1.ObjectMeta) error {
					sa.SetAnnotations()
					objectMeta.SetAnnotations(sa.AsObjectMeta().Annotations)
					objectMeta.SetLabels(sa.AsObjectMeta().Labels)
					if err := maybeCreateServiceAccountOrUpdateMetadata(clientSet, objectMeta); err != nil {
						return errors.Wrapf(err, "failed to create service account %s/%s", objectMeta.GetNamespace(), objectMeta.GetName())
					}
					return nil
				},
			})
		}

		taskTree.Append(saTasks)
	}
	return taskTree
}
