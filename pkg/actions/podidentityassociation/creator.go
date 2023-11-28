package podidentityassociation

import (
	"context"
	"fmt"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

type Creator struct {
	clusterName string

	stackManager StackManager
	eksAPI       awsapi.EKS
}

func NewCreator(clusterName string, stackManager StackManager, eksAPI awsapi.EKS) *Creator {
	return &Creator{
		clusterName:  clusterName,
		stackManager: stackManager,
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
				info:                   fmt.Sprintf("create IAM role for pod identity association for service account %s in namespace %s", pia.ServiceAccountName, pia.Namespace),
				clusterName:            c.clusterName,
				podIdentityAssociation: &podIdentityAssociations[i],
				stackManager:           c.stackManager,
			})
		}
		piaCreationTasks.Append(&createPodIdentityAssociationTask{
			ctx:                    ctx,
			info:                   fmt.Sprintf("create pod identity association for service account %s in namespace %s", pia.ServiceAccountName, pia.Namespace),
			clusterName:            c.clusterName,
			podIdentityAssociation: &podIdentityAssociations[i],
			eksAPI:                 c.eksAPI,
		})
		taskTree.Append(piaCreationTasks)
	}
	return taskTree
}
