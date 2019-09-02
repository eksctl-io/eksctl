package manager

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/cloudformation"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

// NewTasksToDeleteClusterWithNodeGroups defines tasks required to delete the given cluster along with all of its resources
func (c *StackCollection) NewTasksToDeleteClusterWithNodeGroups(oidc *iamoidc.OpenIDConnectManager, clientSetGetter kubernetes.ClientSetGetter, wait bool, cleanup func(chan error, string) error) (*TaskTree, error) {
	tasks := &TaskTree{Parallel: false}

	deleteAll := func(_ string) bool { return true }
	nodeGroupTasks, err := c.NewTasksToDeleteNodeGroups(deleteAll, true, cleanup)

	if err != nil {
		return nil, err
	}
	if nodeGroupTasks.Len() > 0 {
		nodeGroupTasks.IsSubTask = true
		tasks.Append(nodeGroupTasks)
	}

	serviceAccountAndOIDCTasks, err := c.NewTasksToDeleteOIDCProviderWithIAMServiceAccounts(oidc, clientSetGetter, true)
	if err != nil {
		return nil, err
	}

	if serviceAccountAndOIDCTasks.Len() > 0 {
		serviceAccountAndOIDCTasks.IsSubTask = true
		tasks.Append(serviceAccountAndOIDCTasks)
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

// NewTasksToDeleteNodeGroups defines tasks required to delete all of the nodegroups
func (c *StackCollection) NewTasksToDeleteNodeGroups(shouldDelete func(string) bool, wait bool, cleanup func(chan error, string) error) (*TaskTree, error) {
	nodeGroupStacks, err := c.DescribeNodeGroupStacks()
	if err != nil {
		return nil, err
	}

	tasks := &TaskTree{Parallel: true}

	for _, s := range nodeGroupStacks {
		name := c.GetNodeGroupName(s)

		if !shouldDelete(name) {
			continue
		}
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
func (c *StackCollection) NewTasksToDeleteOIDCProviderWithIAMServiceAccounts(oidc *iamoidc.OpenIDConnectManager, clientSetGetter kubernetes.ClientSetGetter, wait bool) (*TaskTree, error) {
	providerExists, err := oidc.CheckProviderExists()
	if err != nil {
		return nil, err
	}

	if !providerExists {
		return nil, fmt.Errorf("unable to delete iamserviceaccount(s) without IAM OIDC provider enabled")
	}

	tasks := &TaskTree{Parallel: false}

	saTasks, err := c.NewTasksToDeleteIAMServiceAccounts(c.spec.IAM.ServiceAccounts, oidc, clientSetGetter, true)
	if err != nil {
		return nil, err
	}

	if saTasks.Len() > 0 {
		saTasks.IsSubTask = true
		tasks.Append(saTasks)
	}
	tasks.Append(&asyncTaskWithoutParams{
		info: "delete IAM OIDC provider",
		call: oidc.DeleteProvider,
	})
	return tasks, nil
}

// NewTasksToDeleteIAMServiceAccounts defines tasks required to delete all of the iamserviceaccounts
func (c *StackCollection) NewTasksToDeleteIAMServiceAccounts(serviceAccounts []*api.ClusterIAMServiceAccount, oidc *iamoidc.OpenIDConnectManager, clientSetGetter kubernetes.ClientSetGetter, wait bool) (*TaskTree, error) {
	serviceAccountStacks, err := c.DescribeIAMServiceAccountStacks()
	if err != nil {
		return nil, err
	}

	tasks := &TaskTree{Parallel: true}

	hasServiceAccount := func(saName string) bool {
		for _, sa := range serviceAccounts {
			if saName == sa.Name {
				return true
			}
		}
		return false
	}

	for _, s := range serviceAccountStacks {
		saTasks := &TaskTree{
			Parallel:  false,
			IsSubTask: true,
		}
		name := c.GetIAMServiceAccountName(s)

		if !hasServiceAccount(name) {
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
