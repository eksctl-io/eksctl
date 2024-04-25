package podidentityassociation

import (
	"context"
	"fmt"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
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
}

func NewCreator(clusterName string, stackCreator StackCreator, eksAPI awsapi.EKS) *Creator {
	return &Creator{
		clusterName:  clusterName,
		stackCreator: stackCreator,
		eksAPI:       eksAPI,
	}
}

func (c *Creator) CreatePodIdentityAssociations(ctx context.Context, podIdentityAssociations []api.PodIdentityAssociation) error {
	return runAllTasks(c.CreateTasks(ctx, podIdentityAssociations))
}

func (c *Creator) CreateTasks(ctx context.Context, podIdentityAssociations []api.PodIdentityAssociation) *tasks.TaskTree {
	taskTree := &tasks.TaskTree{
		Parallel: true,
	}
	for i, pia := range podIdentityAssociations {
		piaCreationTasks := &tasks.TaskTree{
			Parallel:  false,
			IsSubTask: true,
		}
		if pia.RoleARN == "" {
			piaCreationTasks.Append(&createIAMRoleTask{
				ctx:                    ctx,
				info:                   fmt.Sprintf("create IAM role for pod identity association for service account %q", pia.NameString()),
				clusterName:            c.clusterName,
				podIdentityAssociation: &podIdentityAssociations[i],
				stackCreator:           c.stackCreator,
			})
		}
		piaCreationTasks.Append(&createPodIdentityAssociationTask{
			ctx:                    ctx,
			info:                   fmt.Sprintf("create pod identity association for service account %q", pia.NameString()),
			clusterName:            c.clusterName,
			podIdentityAssociation: &podIdentityAssociations[i],
			eksAPI:                 c.eksAPI,
		})
		taskTree.Append(piaCreationTasks)
	}
	return taskTree
}
