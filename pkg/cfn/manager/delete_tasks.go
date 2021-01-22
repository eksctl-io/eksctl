package manager

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/cloudformation"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

func deleteAll(_ string) bool { return true }

// NewTasksToDeleteClusterWithNodeGroups defines tasks required to delete the given cluster along with all of its resources
func (c *StackCollection) NewTasksToDeleteClusterWithNodeGroups(deleteOIDCProvider bool, oidc *iamoidc.OpenIDConnectManager, clientSetGetter kubernetes.ClientSetGetter, wait bool, cleanup func(chan error, string) error) (*tasks.TaskTree, error) {
	taskTree := &tasks.TaskTree{Parallel: false}

	nodeGroupTasks, err := c.NewTasksToDeleteNodeGroups(deleteAll, true, cleanup)

	if err != nil {
		return nil, err
	}
	if nodeGroupTasks.Len() > 0 {
		nodeGroupTasks.IsSubTask = true
		taskTree.Append(nodeGroupTasks)
	}

	if deleteOIDCProvider {
		serviceAccountAndOIDCTasks, err := c.NewTasksToDeleteOIDCProviderWithIAMServiceAccounts(oidc, clientSetGetter)
		if err != nil {
			return nil, err
		}

		if serviceAccountAndOIDCTasks.Len() > 0 {
			serviceAccountAndOIDCTasks.IsSubTask = true
			taskTree.Append(serviceAccountAndOIDCTasks)
		}
	}

	deleteAddonIAMtasks, err := c.NewTaskToDeleteAddonIAM(wait)
	if err != nil {
		return nil, err
	}

	if deleteAddonIAMtasks.Len() > 0 {
		deleteAddonIAMtasks.IsSubTask = true
		taskTree.Append(deleteAddonIAMtasks)
	}

	clusterStack, err := c.DescribeClusterStack()
	if err != nil {
		return nil, err
	}
	if clusterStack == nil {
		return nil, c.ErrStackNotFound()
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
func (c *StackCollection) NewTasksToDeleteNodeGroups(shouldDelete func(string) bool, wait bool, cleanup func(chan error, string) error) (*tasks.TaskTree, error) {
	nodeGroupStacks, err := c.DescribeNodeGroupStacks()
	if err != nil {
		return nil, err
	}

	taskTree := &tasks.TaskTree{Parallel: true}

	for _, s := range nodeGroupStacks {
		name := c.GetNodeGroupName(s)

		if !shouldDelete(name) {
			continue
		}
		if *s.StackStatus == cloudformation.StackStatusDeleteFailed && cleanup != nil {
			taskTree.Append(&tasks.TaskWithNameParam{
				Info: fmt.Sprintf("cleanup for nodegroup %q", name),
				Call: cleanup,
			})
		}
		info := fmt.Sprintf("delete nodegroup %q", name)
		if wait {
			taskTree.Append(&taskWithStackSpec{
				info:  info,
				stack: s,
				call:  c.DeleteStackBySpecSync,
			})
		} else {
			taskTree.Append(&asyncTaskWithStackSpec{
				info:  info,
				stack: s,
				call:  c.DeleteStackBySpec,
			})
		}
	}

	return taskTree, nil
}

// NewTasksToDeleteOIDCProviderWithIAMServiceAccounts defines tasks required to delete all of the iamserviceaccounts
// along with associated IAM ODIC provider
func (c *StackCollection) NewTasksToDeleteOIDCProviderWithIAMServiceAccounts(oidc *iamoidc.OpenIDConnectManager, clientSetGetter kubernetes.ClientSetGetter) (*tasks.TaskTree, error) {
	taskTree := &tasks.TaskTree{Parallel: false}

	saTasks, err := c.NewTasksToDeleteIAMServiceAccounts(deleteAll, clientSetGetter, true)
	if err != nil {
		return nil, err
	}

	if saTasks.Len() > 0 {
		saTasks.IsSubTask = true
		taskTree.Append(saTasks)
	}

	providerExists, err := oidc.CheckProviderExists()
	if err != nil {
		return nil, err
	}

	if providerExists {
		taskTree.Append(&asyncTaskWithoutParams{
			info: "delete IAM OIDC provider",
			call: oidc.DeleteProvider,
		})
	}

	return taskTree, nil
}

// NewTasksToDeleteIAMServiceAccounts defines tasks required to delete all of the iamserviceaccounts
func (c *StackCollection) NewTasksToDeleteIAMServiceAccounts(shouldDelete func(string) bool, clientSetGetter kubernetes.ClientSetGetter, wait bool) (*tasks.TaskTree, error) {
	serviceAccountStacks, err := c.DescribeIAMServiceAccountStacks()
	if err != nil {
		return nil, err
	}

	taskTree := &tasks.TaskTree{Parallel: true}

	for _, s := range serviceAccountStacks {
		saTasks := &tasks.TaskTree{
			Parallel:  false,
			IsSubTask: true,
		}
		name := GetIAMServiceAccountName(s)

		if !shouldDelete(name) {
			continue
		}

		info := fmt.Sprintf("delete IAM role for serviceaccount %q", name)
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
		saTasks.Append(&kubernetesTask{
			info:       fmt.Sprintf("delete serviceaccount %q", name),
			kubernetes: clientSetGetter,
			call: func(clientSet kubernetes.Interface) error {
				meta, err := api.ClusterIAMServiceAccountNameStringToClusterIAMMeta(name)
				if err != nil {
					return err
				}
				return kubernetes.MaybeDeleteServiceAccount(clientSet, meta.AsObjectMeta())
			},
		})
		taskTree.Append(saTasks)
	}

	return taskTree, nil
}

// NewTasksToDeleteIAMServiceAccounts defines tasks required to delete all of the iamserviceaccounts
func (c *StackCollection) NewTaskToDeleteAddonIAM(wait bool) (*tasks.TaskTree, error) {
	stacks, err := c.GetIAMAddonsStacks()
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
