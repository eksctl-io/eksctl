package irsa

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

type DeleteTasksBuilder interface {
	DeleteIAMServiceAccountsTasks(ctx context.Context, serviceAccounts []string, wait bool) (*tasks.TaskTree, error)
}

type Remover struct {
	clientSetGetter kubernetes.ClientSetGetter
	stackManager    StackManager
}

func NewRemover(
	clientSetGetter kubernetes.ClientSetGetter,
	stackManager StackManager) *Remover {
	return &Remover{
		clientSetGetter: clientSetGetter,
		stackManager:    stackManager,
	}
}

func (r *Remover) Delete(ctx context.Context, serviceAccounts []string, plan, wait bool) error {
	taskTree, err := r.DeleteIAMServiceAccountsTasks(ctx, serviceAccounts, wait)
	if err != nil {
		return err
	}
	taskTree.PlanMode = plan

	err = doTasks(taskTree, actionDelete)

	logPlanModeWarning(plan && taskTree.Len() > 0)
	return err
}

func (r *Remover) DeleteIAMServiceAccountsTasks(ctx context.Context, serviceAccounts []string, wait bool) (*tasks.TaskTree, error) {
	serviceAccountStacks, err := r.stackManager.DescribeIAMServiceAccountStacks(ctx)
	if err != nil {
		return nil, err
	}

	stacksMap := stacksToServiceAccountMap(serviceAccountStacks)
	taskTree := &tasks.TaskTree{Parallel: true}

	for _, serviceAccount := range serviceAccounts {
		saTasks := &tasks.TaskTree{
			Parallel:  false,
			IsSubTask: true,
		}

		if s, ok := stacksMap[serviceAccount]; ok {
			info := fmt.Sprintf("delete IAM role for serviceaccount %q", serviceAccount)
			saTasks.Append(&deleteIAMServiceAccountTask{
				ctx:          ctx,
				info:         info,
				stack:        s,
				stackManager: r.stackManager,
				wait:         wait,
			})
		}

		meta, err := api.ClusterIAMServiceAccountNameStringToClusterIAMMeta(serviceAccount)
		if err != nil {
			return nil, err
		}
		saTasks.Append(&kubernetesTask{
			info:       fmt.Sprintf("delete serviceaccount %q", serviceAccount),
			kubernetes: r.clientSetGetter,
			objectMeta: meta.AsObjectMeta(),
			call:       kubernetes.MaybeDeleteServiceAccount,
		})
		taskTree.Append(saTasks)
	}

	return taskTree, nil
}

func stacksToServiceAccountMap(stacks []*types.Stack) map[string]*types.Stack {
	stackMap := make(map[string]*types.Stack)
	for _, stack := range stacks {
		stackMap[getIAMServiceAccountName(stack)] = stack
	}

	return stackMap
}

func getIAMServiceAccountName(s *types.Stack) string {
	for _, tag := range s.Tags {
		if *tag.Key == api.IAMServiceAccountNameTag {
			return *tag.Value
		}
	}
	return ""
}
