package podidentityassociation

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclient "k8s.io/client-go/kubernetes"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_stack_creator.go . StackCreator
type StackCreator interface {
	CreateStack(ctx context.Context, name string, stack builder.ResourceSetReader, tags, parameters map[string]string, errs chan error) error
}

type Creator struct {
	clusterName string

	stackCreator StackCreator
	eksAPI       awsapi.EKS
	clientSet    kubeclient.Interface
}

func NewCreator(clusterName string, stackCreator StackCreator, eksAPI awsapi.EKS, clientSet kubeclient.Interface) *Creator {
	return &Creator{
		clusterName:  clusterName,
		stackCreator: stackCreator,
		eksAPI:       eksAPI,
		clientSet:    clientSet,
	}
}

func (c *Creator) CreatePodIdentityAssociations(ctx context.Context, podIdentityAssociations []api.PodIdentityAssociation) error {
	return runAllTasks(c.CreateTasks(ctx, podIdentityAssociations, false))
}

func (c *Creator) CreateTasks(ctx context.Context, podIdentityAssociations []api.PodIdentityAssociation, ignorePodIdentityExistsErr bool) *tasks.TaskTree {
	taskTree := &tasks.TaskTree{
		Parallel: true,
	}
	for _, pia := range podIdentityAssociations {
		pia := pia
		piaCreationTasks := &tasks.TaskTree{
			Parallel:  false,
			IsSubTask: true,
		}
		if pia.RoleARN == "" {
			piaCreationTasks.Append(&tasks.GenericTask{
				Description: fmt.Sprintf("create IAM role for pod identity association for service account %q", pia.NameString()),
				Doer: func() error {
					roleCreator := &IAMRoleCreator{
						ClusterName:  c.clusterName,
						StackCreator: c.stackCreator,
					}
					roleARN, err := roleCreator.Create(ctx, &pia, "")
					if err != nil {
						return err
					}
					pia.RoleARN = roleARN
					return nil
				},
			})
		}
		if pia.CreateServiceAccount {
			piaCreationTasks.Append(&tasks.GenericTask{
				Description: fmt.Sprintf("create service account %q, if it does not already exist", pia.NameString()),
				Doer: func() error {
					if err := kubernetes.MaybeCreateServiceAccountOrUpdateMetadata(c.clientSet, v1.ObjectMeta{
						Name:      pia.ServiceAccountName,
						Namespace: pia.Namespace,
					}); err != nil {
						return fmt.Errorf("failed to create service account %q: %w", pia.NameString(), err)
					}
					return nil
				},
			})
		}
		piaCreationTasks.Append(&createPodIdentityAssociationTask{
			ctx:                        ctx,
			info:                       fmt.Sprintf("create pod identity association for service account %q", pia.NameString()),
			clusterName:                c.clusterName,
			podIdentityAssociation:     &pia,
			eksAPI:                     c.eksAPI,
			ignorePodIdentityExistsErr: ignorePodIdentityExistsErr,
		})
		taskTree.Append(piaCreationTasks)
	}
	return taskTree
}
