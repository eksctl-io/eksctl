package manager

import (
	"context"
	"fmt"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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

// NewTasksToCreateClusterWithNodeGroups defines all tasks required to create a cluster along
// with some nodegroups; see CreateAllNodeGroups for how onlyNodeGroupSubset works.
func (c *StackCollection) NewTasksToCreateClusterWithNodeGroups(ctx context.Context, nodeGroups []*api.NodeGroup,
	managedNodeGroups []*api.ManagedNodeGroup, postClusterCreationTasks ...tasks.Task) *tasks.TaskTree {

	taskTree := tasks.TaskTree{Parallel: false}

	taskTree.Append(
		&createClusterTask{
			info:                 fmt.Sprintf("create cluster control plane %q", c.spec.Metadata.Name),
			stackCollection:      c,
			supportsManagedNodes: true,
			ctx:                  ctx,
		},
	)

	if c.spec.IPv6Enabled() {
		taskTree.Append(
			&AssignIpv6AddressOnCreationTask{
				ClusterConfig: c.spec,
				EC2API:        c.ec2API,
				Context:       ctx,
			},
		)
	}

	appendNodeGroupTasksTo := func(taskTree *tasks.TaskTree) {
		vpcImporter := vpc.NewStackConfigImporter(c.MakeClusterStackName())
		nodeGroupTasks := c.NewUnmanagedNodeGroupTask(ctx, nodeGroups, false, vpcImporter)
		managedNodeGroupTasks := c.NewManagedNodeGroupTask(ctx, managedNodeGroups, false, vpcImporter)
		if managedNodeGroupTasks.Len() > 0 {
			nodeGroupTasks.Append(managedNodeGroupTasks.Tasks...)
		}

		if nodeGroupTasks.Len() > 0 {
			nodeGroupTasks.IsSubTask = true
			taskTree.Append(nodeGroupTasks)
		}
	}

	if len(postClusterCreationTasks) > 0 {
		postClusterCreationTaskTree := tasks.TaskTree{
			Parallel:  false,
			IsSubTask: true,
		}
		postClusterCreationTaskTree.Append(postClusterCreationTasks...)
		appendNodeGroupTasksTo(&postClusterCreationTaskTree)
		taskTree.Append(&postClusterCreationTaskTree)
	} else {
		appendNodeGroupTasksTo(&taskTree)
	}

	return &taskTree
}

// NewUnmanagedNodeGroupTask defines tasks required to create all of the nodegroups
func (c *StackCollection) NewUnmanagedNodeGroupTask(ctx context.Context, nodeGroups []*api.NodeGroup, forceAddCNIPolicy bool, vpcImporter vpc.Importer) *tasks.TaskTree {
	taskTree := &tasks.TaskTree{Parallel: true}

	for _, ng := range nodeGroups {
		taskTree.Append(&nodeGroupTask{
			info:              fmt.Sprintf("create nodegroup %q", ng.NameString()),
			ctx:               ctx,
			nodeGroup:         ng,
			stackCollection:   c,
			forceAddCNIPolicy: forceAddCNIPolicy,
			vpcImporter:       vpcImporter,
		})
		// TODO: move authconfigmap tasks here using kubernetesTask and kubernetes.CallbackClientSet
	}

	return taskTree
}

// NewManagedNodeGroupTask defines tasks required to create managed nodegroups
func (c *StackCollection) NewManagedNodeGroupTask(ctx context.Context, nodeGroups []*api.ManagedNodeGroup, forceAddCNIPolicy bool, vpcImporter vpc.Importer) *tasks.TaskTree {
	taskTree := &tasks.TaskTree{Parallel: true}
	for _, ng := range nodeGroups {
		taskTree.Append(&managedNodeGroupTask{
			stackCollection:   c,
			nodeGroup:         ng,
			forceAddCNIPolicy: forceAddCNIPolicy,
			vpcImporter:       vpcImporter,
			info:              fmt.Sprintf("create managed nodegroup %q", ng.Name),
			ctx:               ctx,
		})
		if api.IsEnabled(ng.PropagateASGTags) {
			// disable parallelisation if any tags propagation is done
			// since nodegroup must be created to propagate tags to its ASGs
			taskTree.Parallel = false
			taskTree.Append(&managedNodeGroupTagsToASGPropagationTask{
				stackCollection: c,
				nodeGroup:       ng,
				info:            fmt.Sprintf("propagate tags to ASG for managed nodegroup %q", ng.Name),
			})
		}
	}
	return taskTree
}

// NewClusterCompatTask creates a new task that checks for cluster compatibility with new features like
// Managed Nodegroups and Fargate, and updates the CloudFormation cluster stack if the required resources are missing
func (c *StackCollection) NewClusterCompatTask(ctx context.Context) tasks.Task {
	return &clusterCompatTask{
		stackCollection: c,
		info:            "fix cluster compatibility",
		ctx:             ctx,
	}
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
