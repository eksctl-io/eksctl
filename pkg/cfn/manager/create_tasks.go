package manager

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"

	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// NewTasksToCreateClusterWithNodeGroups defines all tasks required to create a cluster along
// with some nodegroups; see CreateAllNodeGroups for how onlyNodeGroupSubset works
func (c *StackCollection) NewTasksToCreateClusterWithNodeGroups(nodeGroups []*api.NodeGroup) *TaskTree {
	tasks := &TaskTree{Parallel: false}

	tasks.Append(
		&taskWithoutParams{
			info: fmt.Sprintf("create cluster control plane %q", c.spec.Metadata.Name),
			call: c.createClusterTask,
		},
	)

	nodeGroupTasks := c.NewTasksToCreateNodeGroups(nodeGroups)
	if nodeGroupTasks.Len() > 0 {
		nodeGroupTasks.IsSubTask = true
		tasks.Append(nodeGroupTasks)
	}

	return tasks
}

// NewTasksToCreateNodeGroups defines tasks required to create all of the nodegroups
func (c *StackCollection) NewTasksToCreateNodeGroups(nodeGroups []*api.NodeGroup) *TaskTree {
	tasks := &TaskTree{Parallel: true}

	for i := range c.spec.NodeGroups {
		ng := c.spec.NodeGroups[i]
		if onlySubset != nil && !onlySubset.Has(ng.NameString()) {
			continue
		}
		tasks.Append(&taskWithNodeGroupSpec{
			info:      fmt.Sprintf("create nodegroup %q", ng.NameString()),
			nodeGroup: ng,
			call:      c.createNodeGroupTask,
		})
		// TODO: move authconfigmap tasks here using kubernetesTask and kubernetes.CallbackClientSet
	}

	return tasks
}

// NewTasksToCreateIAMServiceAccounts defines tasks required to create all of the IAM ServiceAccounts if
// onlySubset is nil, otherwise just the tasks for IAM ServiceAccounts that are in onlySubset
// will be defined
func (c *StackCollection) NewTasksToCreateIAMServiceAccounts(onlySubset sets.String, oidc *iamoidc.OpenIDConnectManager, clientSetGetter kubernetes.ClientSetGetter) *TaskTree {
	tasks := &TaskTree{Parallel: true}

	for i := range c.spec.IAM.ServiceAccounts {
		sa := c.spec.IAM.ServiceAccounts[i]
		saTasks := &TaskTree{
			Parallel:  false,
			IsSubTask: true,
		}
		if onlySubset != nil && !onlySubset.Has(sa.NameString()) {
			continue
		}

		saTasks.Append(&taskWithClusterIAMServiceAccountSpec{
			info:           fmt.Sprintf("create IAM role for serviceaccount %q", sa.NameString()),
			serviceAccount: sa,
			oidc:           oidc,
			call:           c.createIAMServiceAccountTask,
		})
		saTasks.Append(&kubernetesTask{
			info:       fmt.Sprintf("create serviceaccount %q", sa.NameString()),
			kubernetes: clientSetGetter,
			call: func(clientSet kubernetes.Interface) error {
				sa.SetAnnotations()
				return kubernetes.MaybeCreateServiceAccountOrUpdateMetadata(clientSet, sa.ObjectMeta)
			},
		})

		tasks.Append(saTasks)
	}
	return tasks
}
