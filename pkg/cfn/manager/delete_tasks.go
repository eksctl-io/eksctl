package manager

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/eks"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/cfn/waiter"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

// think Jake is deleting this soon
func deleteAll(_ string) bool { return true }

// NewTasksToDeleteClusterWithNodeGroups defines tasks required to delete the given cluster along with all of its resources
func (c *StackCollection) NewTasksToDeleteClusterWithNodeGroups(ctx context.Context, clusterStack *Stack, nodeGroupStacks []NodeGroupStack, deleteOIDCProvider bool, oidc *iamoidc.OpenIDConnectManager, clientSetGetter kubernetes.ClientSetGetter, wait bool, cleanup func(chan error, string) error) (*tasks.TaskTree, error) {
	taskTree := &tasks.TaskTree{Parallel: false}

	nodeGroupTasks, err := c.NewTasksToDeleteNodeGroups(nodeGroupStacks, deleteAll, true, cleanup)

	if err != nil {
		return nil, err
	}
	if nodeGroupTasks.Len() > 0 {
		nodeGroupTasks.IsSubTask = true
		taskTree.Append(nodeGroupTasks)
	}

	if deleteOIDCProvider {
		serviceAccountAndOIDCTasks, err := c.NewTasksToDeleteOIDCProviderWithIAMServiceAccounts(ctx, oidc, clientSetGetter)
		if err != nil {
			return nil, err
		}

		if serviceAccountAndOIDCTasks.Len() > 0 {
			serviceAccountAndOIDCTasks.IsSubTask = true
			taskTree.Append(serviceAccountAndOIDCTasks)
		}
	}

	deleteAddonIAMtasks, err := c.NewTaskToDeleteAddonIAM(ctx, wait)
	if err != nil {
		return nil, err
	}

	if deleteAddonIAMtasks.Len() > 0 {
		deleteAddonIAMtasks.IsSubTask = true
		taskTree.Append(deleteAddonIAMtasks)
	}

	if clusterStack == nil {
		return nil, &StackNotFoundErr{ClusterName: c.spec.Metadata.Name}
	}

	info := fmt.Sprintf("delete cluster control plane %q", c.spec.Metadata.Name)
	if wait {
		taskTree.Append(&taskWithStackSpec{
			info:  info,
			stack: clusterStack,
			call:  c.DeleteStackBySpecSync,
		})
	} else {
		taskTree.Append(&asyncTaskWithStackSpec{
			info:  info,
			stack: clusterStack,
			call:  c.DeleteStackBySpec,
		})
	}

	return taskTree, nil
}

// NewTasksToDeleteNodeGroups defines tasks required to delete all of the nodegroups
func (c *StackCollection) NewTasksToDeleteNodeGroups(nodeGroupStacks []NodeGroupStack, shouldDelete func(string) bool, wait bool, cleanup func(chan error, string) error) (*tasks.TaskTree, error) {
	taskTree := &tasks.TaskTree{Parallel: true}

	for _, s := range nodeGroupStacks {

		if !shouldDelete(s.NodeGroupName) {
			continue
		}

		if s.Stack.StackStatus == types.StackStatusDeleteFailed && cleanup != nil {
			taskTree.Append(&tasks.TaskWithNameParam{
				Info: fmt.Sprintf("cleanup for nodegroup %q", s.NodeGroupName),
				Call: cleanup,
			})
		}
		info := fmt.Sprintf("delete nodegroup %q", s.NodeGroupName)
		if wait {
			taskTree.Append(&taskWithStackSpec{
				info:  info,
				stack: s.Stack,
				call:  c.DeleteStackBySpecSync,
			})
		} else {
			taskTree.Append(&asyncTaskWithStackSpec{
				info:  info,
				stack: s.Stack,
				call:  c.DeleteStackBySpec,
			})
		}
	}

	return taskTree, nil
}

type DeleteWaitCondition struct {
	Condition func() (bool, error)
	Timeout   time.Duration
	Interval  time.Duration
}

type DeleteUnownedNodegroupTask struct {
	cluster   string
	nodegroup string
	wait      *DeleteWaitCondition
	info      string
	eksAPI    awsapi.EKS
	ctx       context.Context
}

func (d *DeleteUnownedNodegroupTask) Describe() string {
	return d.info
}

func (d *DeleteUnownedNodegroupTask) Do() error {
	out, err := d.eksAPI.DeleteNodegroup(d.ctx, &eks.DeleteNodegroupInput{
		ClusterName:   &d.cluster,
		NodegroupName: &d.nodegroup,
	})
	if err != nil {
		return err
	}

	if d.wait != nil {
		w := waiter.Waiter{
			NextDelay: func(_ int) time.Duration {
				return d.wait.Interval
			},
			Operation: d.wait.Condition,
		}

		if err := w.WaitWithTimeout(d.wait.Timeout); err != nil {
			if err == context.DeadlineExceeded {
				return errors.Errorf("timed out waiting for nodegroup deletion after %s", d.wait.Timeout)
			}
			return err
		}
	}

	if out != nil {
		logger.Debug("delete nodegroup %q output: %+v", d.nodegroup, out.Nodegroup)
	}
	return nil
}

func (c *StackCollection) NewTaskToDeleteUnownedNodeGroup(ctx context.Context, clusterName, nodegroup string, eksAPI awsapi.EKS, waitCondition *DeleteWaitCondition) tasks.Task {
	return tasks.SynchronousTask{
		SynchronousTaskIface: &DeleteUnownedNodegroupTask{
			cluster:   clusterName,
			nodegroup: nodegroup,
			eksAPI:    eksAPI,
			wait:      waitCondition,
			info:      fmt.Sprintf("delete unowned nodegroup %s", nodegroup),
			ctx:       ctx,
		}}
}

// NewTasksToDeleteOIDCProviderWithIAMServiceAccounts defines tasks required to delete all of the iamserviceaccounts
// along with associated IAM ODIC provider
func (c *StackCollection) NewTasksToDeleteOIDCProviderWithIAMServiceAccounts(ctx context.Context, oidc *iamoidc.OpenIDConnectManager, clientSetGetter kubernetes.ClientSetGetter) (*tasks.TaskTree, error) {
	taskTree := &tasks.TaskTree{Parallel: false}

	allServiceAccountsWithStacks, err := c.getAllServiceAccounts(ctx)
	if err != nil {
		return nil, err
	}
	saTasks, err := c.NewTasksToDeleteIAMServiceAccounts(ctx, allServiceAccountsWithStacks, clientSetGetter, true)
	if err != nil {
		return nil, err
	}

	if saTasks.Len() > 0 {
		saTasks.IsSubTask = true
		taskTree.Append(saTasks)
	}

	providerExists, err := oidc.CheckProviderExists(ctx)
	if err != nil {
		return nil, err
	}

	if providerExists {
		taskTree.Append(&asyncTaskWithoutParams{
			info: "delete IAM OIDC provider",
			call: func() error {
				return oidc.DeleteProvider(ctx)
			},
		})
	}

	return taskTree, nil
}

func (c *StackCollection) getAllServiceAccounts(ctx context.Context) ([]string, error) {
	serviceAccountStacks, err := c.DescribeIAMServiceAccountStacks(ctx)
	if err != nil {
		return nil, err
	}

	var serviceAccounts []string
	for _, stack := range serviceAccountStacks {
		serviceAccounts = append(serviceAccounts, GetIAMServiceAccountName(stack))
	}

	return serviceAccounts, nil
}

// NewTasksToDeleteIAMServiceAccounts defines tasks required to delete all of the iamserviceaccounts
func (c *StackCollection) NewTasksToDeleteIAMServiceAccounts(ctx context.Context, serviceAccounts []string, clientSetGetter kubernetes.ClientSetGetter, wait bool) (*tasks.TaskTree, error) {
	serviceAccountStacks, err := c.DescribeIAMServiceAccountStacks(ctx)
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
			if wait {
				saTasks.Append(&taskWithStackSpec{
					info:  info,
					stack: s,
					call:  c.DeleteStackBySpecSync,
				})
			} else {
				saTasks.Append(&asyncTaskWithStackSpec{
					info:  info,
					stack: s,
					call:  c.DeleteStackBySpec,
				})
			}
		}

		meta, err := api.ClusterIAMServiceAccountNameStringToClusterIAMMeta(serviceAccount)
		if err != nil {
			return nil, err
		}
		saTasks.Append(&kubernetesTask{
			info:       fmt.Sprintf("delete serviceaccount %q", serviceAccount),
			kubernetes: clientSetGetter,
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
		stackMap[GetIAMServiceAccountName(stack)] = stack
	}

	return stackMap
}

// NewTaskToDeleteAddonIAM defines tasks required to delete all of the addons
func (c *StackCollection) NewTaskToDeleteAddonIAM(ctx context.Context, wait bool) (*tasks.TaskTree, error) {
	stacks, err := c.GetIAMAddonsStacks(ctx)
	if err != nil {
		return nil, err
	}
	taskTree := &tasks.TaskTree{Parallel: true}
	for _, s := range stacks {
		info := fmt.Sprintf("delete addon IAM %q", *s.StackName)

		deleteStackTasks := &tasks.TaskTree{
			Parallel:  false,
			IsSubTask: true,
		}
		if wait {
			deleteStackTasks.Append(&taskWithStackSpec{
				info:  info,
				stack: s,
				call:  c.DeleteStackBySpecSync,
			})
		} else {
			deleteStackTasks.Append(&asyncTaskWithStackSpec{
				info:  info,
				stack: s,
				call:  c.DeleteStackBySpec,
			})
		}
		taskTree.Append(deleteStackTasks)
	}
	return taskTree, nil

}
