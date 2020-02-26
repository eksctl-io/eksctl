package manager

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

func deleteAll(_ string) bool { return true }

// NewTasksToDeleteClusterWithNodeGroups defines tasks required to delete the given cluster along with all of its resources
func (c *StackCollection) NewTasksToDeleteClusterWithNodeGroups(deleteOIDCProvider bool, oidc *iamoidc.OpenIDConnectManager, clientSetGetter kubernetes.ClientSetGetter, wait bool, cleanup func(chan error, string) error) (*TaskTree, error) {
	tasks := &TaskTree{Parallel: false}

	nodeGroupStacks, err := c.DescribeNodeGroupStacks()
	if err != nil {
		return nil, err
	}

	nodeGroupTasks, err := c.MakeTasksToDeleteNodeGroupsFromStacks(nodeGroupStacks, true, cleanup)

	if err != nil {
		return nil, err
	}
	if nodeGroupTasks.Len() > 0 {
		nodeGroupTasks.IsSubTask = true
		tasks.Append(nodeGroupTasks)
	}

	if deleteOIDCProvider {
		serviceAccountAndOIDCTasks, err := c.NewTasksToDeleteOIDCProviderWithIAMServiceAccounts(oidc, clientSetGetter)
		if err != nil {
			return nil, err
		}

		if serviceAccountAndOIDCTasks.Len() > 0 {
			serviceAccountAndOIDCTasks.IsSubTask = true
			tasks.Append(serviceAccountAndOIDCTasks)
		}
	}

	clusterStack, err := c.DescribeClusterStack()
	if err != nil {
		return nil, err
	}

	info := fmt.Sprintf("delete cluster control plane %q", c.spec.Metadata.Name)
	if wait {
		tasks.Append(&taskWithStackSpec{
			info:  info,
			stack: clusterStack,
			call:  c.DeleteStackBySpecSync,
		})
	} else {
		tasks.Append(&asyncTaskWithStackSpec{
			info:  info,
			stack: clusterStack,
			call:  c.DeleteStackBySpec,
		})
	}

	return tasks, nil
}

// MakeTasksToDeleteNodeGroupsFromStacks creates tasks to delete the nodegroups given in a list of cloudformation node group stack descriptions
func (c *StackCollection) MakeTasksToDeleteNodeGroupsFromStacks(stacks []*Stack, wait bool, cleanup func(chan error, string) error) (*TaskTree, error) {
	tasks := &TaskTree{Parallel: true}

	for _, s := range stacks {
		name := c.GetNodeGroupName(s)

		if *s.StackStatus == cloudformation.StackStatusDeleteFailed && cleanup != nil {
			tasks.Append(&taskWithNameParam{
				info: fmt.Sprintf("cleanup for nodegroup %q", name),
				call: cleanup,
			})
		}
		info := fmt.Sprintf("delete nodegroup %q", name)
		if wait {
			tasks.Append(&taskWithStackSpec{
				info:  info,
				stack: s,
				call:  c.DeleteStackBySpecSync,
			})
		} else {
			tasks.Append(&asyncTaskWithStackSpec{
				info:  info,
				stack: s,
				call:  c.DeleteStackBySpec,
			})
		}
	}

	return tasks, nil
}

// NewTasksToDeleteOIDCProviderWithIAMServiceAccounts defines tasks required to delete all of the iamserviceaccounts
// along with associated IAM ODIC provider
func (c *StackCollection) NewTasksToDeleteOIDCProviderWithIAMServiceAccounts(oidc *iamoidc.OpenIDConnectManager, clientSetGetter kubernetes.ClientSetGetter) (*TaskTree, error) {
	tasks := &TaskTree{Parallel: false}

	saTasks, err := c.NewTasksToDeleteIAMServiceAccounts(deleteAll, oidc, clientSetGetter, true)
	if err != nil {
		return nil, err
	}

	if saTasks.Len() > 0 {
		saTasks.IsSubTask = true
		tasks.Append(saTasks)
	}

	providerExists, err := oidc.CheckProviderExists()
	if err != nil {
		return nil, err
	}

	if providerExists {
		tasks.Append(&asyncTaskWithoutParams{
			info: "delete IAM OIDC provider",
			call: oidc.DeleteProvider,
		})
	}

	return tasks, nil
}

// NewTasksToDeleteIAMServiceAccounts defines tasks required to delete all of the iamserviceaccounts
func (c *StackCollection) NewTasksToDeleteIAMServiceAccounts(shouldDelete func(string) bool, oidc *iamoidc.OpenIDConnectManager, clientSetGetter kubernetes.ClientSetGetter, wait bool) (*TaskTree, error) {
	serviceAccountStacks, err := c.DescribeIAMServiceAccountStacks()
	if err != nil {
		return nil, err
	}

	tasks := &TaskTree{Parallel: true}

	for _, s := range serviceAccountStacks {
		saTasks := &TaskTree{
			Parallel:  false,
			IsSubTask: true,
		}
		name := c.GetIAMServiceAccountName(s)

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
				meta, err := api.ClusterIAMServiceAccountNameStringToObjectMeta(name)
				if err != nil {
					return err
				}
				return kubernetes.MaybeDeleteServiceAccount(clientSet, *meta)
			},
		})
		tasks.Append(saTasks)
	}

	return tasks, nil
}

// NewTasksToDeleteIdentityRoleARNFromAuthConfigMap defines tasks required to delete identity role ARNs from the auth configmap
// node groups in the cluster config must already be populated with IAM role ARNs
func (c *StackCollection) NewTasksToDeleteIdentityRoleARNFromAuthConfigMap(cfgMarkedForDeletion *api.ClusterConfig, clientSet kubernetes.Interface) (*TaskTree, error) {
	tasks := &TaskTree{Parallel: false}

	for _, n := range cfgMarkedForDeletion.NodeGroups {
		tasks.Append(&deleteFromAuthConfigMapTask{
			info:            fmt.Sprintf("remove identity role arn %q from auth configmap", n.IAM.InstanceRoleARN),
			clientSet:       clientSet,
			instanceRoleARN: n.IAM.InstanceRoleARN,
		})
	}

	// deletion of the identity role arn in the auth configmap for managed node groups are
	// currently handled by AWS when deleting the managed node group

	return tasks, nil
}
