package manager

import (
	"context"
	"fmt"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/weaveworks/eksctl/pkg/actions/accessentry"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

const (
	managedByKubernetesLabelKey   = "app.kubernetes.io/managed-by"
	managedByKubernetesLabelValue = "eksctl"
)

// NewTasksToCreateCluster defines all tasks required to create a cluster along
// with some nodegroups; see CreateAllNodeGroups for how onlyNodeGroupSubset works.
func (c *StackCollection) NewTasksToCreateCluster(ctx context.Context, nodeGroups []*api.NodeGroup,
	managedNodeGroups []*api.ManagedNodeGroup, accessEntries []api.AccessEntry, accessEntryCreator accessentry.CreatorInterface, postClusterCreationTasks ...tasks.Task) *tasks.TaskTree {
	taskTree := tasks.TaskTree{Parallel: false}

	taskTree.Append(&createClusterTask{
		info:                 fmt.Sprintf("create cluster control plane %q", c.spec.Metadata.Name),
		stackCollection:      c,
		supportsManagedNodes: true,
		ctx:                  ctx,
	})

	if len(accessEntries) > 0 {
		taskTree.Append(accessEntryCreator.CreateTasks(ctx, accessEntries))
	}

	appendNodeGroupTasksTo := func(taskTree *tasks.TaskTree) {
		vpcImporter := vpc.NewStackConfigImporter(c.MakeClusterStackName())
		nodeGroupTasks := &tasks.TaskTree{
			Parallel:  true,
			IsSubTask: true,
		}
		if unmanagedNodeGroupTasks := c.NewUnmanagedNodeGroupTask(ctx, nodeGroups, false, false, false, vpcImporter); unmanagedNodeGroupTasks.Len() > 0 {
			unmanagedNodeGroupTasks.IsSubTask = true
			nodeGroupTasks.Append(unmanagedNodeGroupTasks)
		}
		if managedNodeGroupTasks := c.NewManagedNodeGroupTask(ctx, managedNodeGroups, false, vpcImporter); managedNodeGroupTasks.Len() > 0 {
			managedNodeGroupTasks.IsSubTask = true
			nodeGroupTasks.Append(managedNodeGroupTasks)
		}

		if nodeGroupTasks.Len() > 0 {
			taskTree.Append(nodeGroupTasks)
		}
	}

	if len(postClusterCreationTasks) > 0 {
		postClusterCreationTaskTree := &tasks.TaskTree{
			Parallel:  false,
			IsSubTask: true,
		}
		postClusterCreationTaskTree.Append(postClusterCreationTasks...)
		appendNodeGroupTasksTo(postClusterCreationTaskTree)
		taskTree.Append(postClusterCreationTaskTree)
	} else {
		appendNodeGroupTasksTo(&taskTree)
	}
	return &taskTree
}

// NewUnmanagedNodeGroupTask defines tasks required to create all nodegroups.
func (c *StackCollection) NewUnmanagedNodeGroupTask(ctx context.Context, nodeGroups []*api.NodeGroup, forceAddCNIPolicy, skipEgressRules, disableAccessEntryCreation bool, vpcImporter vpc.Importer) *tasks.TaskTree {
	taskTree := &tasks.TaskTree{Parallel: true}

	for _, ng := range nodeGroups {
		taskTree.Append(&nodeGroupTask{
			info:                       fmt.Sprintf("create nodegroup %q", ng.NameString()),
			ctx:                        ctx,
			nodeGroup:                  ng,
			stackCollection:            c,
			forceAddCNIPolicy:          forceAddCNIPolicy,
			vpcImporter:                vpcImporter,
			skipEgressRules:            skipEgressRules,
			disableAccessEntryCreation: disableAccessEntryCreation,
		})
		// TODO: move authconfigmap tasks here using kubernetesTask and kubernetes.CallbackClientSet
	}

	return taskTree
}

// NewManagedNodeGroupTask defines tasks required to create managed nodegroups
func (c *StackCollection) NewManagedNodeGroupTask(ctx context.Context, nodeGroups []*api.ManagedNodeGroup, forceAddCNIPolicy bool, vpcImporter vpc.Importer) *tasks.TaskTree {
	taskTree := &tasks.TaskTree{Parallel: true}
	for _, ng := range nodeGroups {
		// Disable parallelisation if any tags propagation is done
		// since nodegroup must be created to propagate tags to its ASGs.
		subTask := &tasks.TaskTree{
			Parallel:  false,
			IsSubTask: true,
		}
		subTask.Append(&managedNodeGroupTask{
			stackCollection:   c,
			nodeGroup:         ng,
			forceAddCNIPolicy: forceAddCNIPolicy,
			vpcImporter:       vpcImporter,
			info:              fmt.Sprintf("create managed nodegroup %q", ng.Name),
			ctx:               ctx,
		})
		if api.IsEnabled(ng.PropagateASGTags) {
			subTask.Append(&managedNodeGroupTagsToASGPropagationTask{
				stackCollection: c,
				nodeGroup:       ng,
				info:            fmt.Sprintf("propagate tags to ASG for managed nodegroup %q", ng.Name),
				ctx:             ctx,
			})
		}
		taskTree.Append(subTask)
	}
	return taskTree
}

// NewTasksToCreateIAMServiceAccounts defines tasks required to create all of the IAM ServiceAccounts
func (c *StackCollection) NewTasksToCreateIAMServiceAccounts(serviceAccounts []*api.ClusterIAMServiceAccount, oidc *iamoidc.OpenIDConnectManager, clientSetGetter kubernetes.ClientSetGetter) *tasks.TaskTree {
	taskTree := &tasks.TaskTree{Parallel: true}

	for i := range serviceAccounts {
		sa := serviceAccounts[i]
		saTasks := &tasks.TaskTree{
			Parallel:  false,
			IsSubTask: true,
		}

		if sa.AttachRoleARN == "" {
			saTasks.Append(&taskWithClusterIAMServiceAccountSpec{
				info:            fmt.Sprintf("create IAM role for serviceaccount %q", sa.NameString()),
				stackCollection: c,
				serviceAccount:  sa,
				oidc:            oidc,
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
				kubernetes: clientSetGetter,
				objectMeta: sa.ClusterIAMMeta.AsObjectMeta(),
				call: func(clientSet kubernetes.Interface, objectMeta v1.ObjectMeta) error {
					sa.SetAnnotations()
					objectMeta.SetAnnotations(sa.AsObjectMeta().Annotations)
					objectMeta.SetLabels(sa.AsObjectMeta().Labels)
					if err := kubernetes.MaybeCreateServiceAccountOrUpdateMetadata(clientSet, objectMeta); err != nil {
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
