package manager

import (
	"context"
	"fmt"

	"github.com/kris-nova/logger"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/weaveworks/eksctl/pkg/actions/accessentry"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

// NewTasksToCreateCluster defines all tasks required to create a cluster along
// with some nodegroups; see CreateAllNodeGroups for how onlyNodeGroupSubset works.
func (c *StackCollection) NewTasksToCreateCluster(ctx context.Context, nodeGroups []*api.NodeGroup,
	managedNodeGroups []*api.ManagedNodeGroup, accessConfig *api.AccessConfig, accessEntryCreator accessentry.CreatorInterface, nodeGroupParallelism int, postClusterCreationTasks ...tasks.Task) *tasks.TaskTree {
	taskTree := tasks.TaskTree{Parallel: false}

	taskTree.Append(&createClusterTask{
		info:                 fmt.Sprintf("create cluster control plane %q", c.spec.Metadata.Name),
		stackCollection:      c,
		supportsManagedNodes: true,
		ctx:                  ctx,
	})

	if len(accessConfig.AccessEntries) > 0 {
		taskTree.Append(accessEntryCreator.CreateTasks(ctx, accessConfig.AccessEntries))
	}

	appendNodeGroupTasksTo := func(taskTree *tasks.TaskTree) {
		vpcImporter := vpc.NewStackConfigImporter(c.MakeClusterStackName())
		nodeGroupTasks := &tasks.TaskTree{
			Parallel:  true,
			IsSubTask: true,
		}
		disableAccessEntryCreation := accessConfig.AuthenticationMode == ekstypes.AuthenticationModeConfigMap
		if unmanagedNodeGroupTasks := c.NewUnmanagedNodeGroupTask(ctx, nodeGroups, false, false, disableAccessEntryCreation, vpcImporter, nodeGroupParallelism); unmanagedNodeGroupTasks.Len() > 0 {
			unmanagedNodeGroupTasks.IsSubTask = true
			nodeGroupTasks.Append(unmanagedNodeGroupTasks)
		}
		if managedNodeGroupTasks := c.NewManagedNodeGroupTask(ctx, managedNodeGroups, false, vpcImporter, nodeGroupParallelism); managedNodeGroupTasks.Len() > 0 {
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

// NewUnmanagedNodeGroupTask returns tasks for creating self-managed nodegroups.
func (c *StackCollection) NewUnmanagedNodeGroupTask(ctx context.Context, nodeGroups []*api.NodeGroup, forceAddCNIPolicy, skipEgressRules, disableAccessEntryCreation bool, vpcImporter vpc.Importer, parallelism int) *tasks.TaskTree {
	task := &UnmanagedNodeGroupTask{
		ClusterConfig: c.spec,
		NodeGroups:    nodeGroups,
		CreateNodeGroupResourceSet: func(options builder.NodeGroupOptions) NodeGroupResourceSet {
			return builder.NewNodeGroupResourceSet(c.ec2API, c.iamAPI, options)
		},
		NewBootstrapper: func(clusterConfig *api.ClusterConfig, ng *api.NodeGroup) (nodebootstrap.Bootstrapper, error) {
			return nodebootstrap.NewBootstrapper(clusterConfig, ng)
		},
		EKSAPI:       c.eksAPI,
		StackManager: c,
	}
	return task.Create(ctx, CreateNodeGroupOptions{
		ForceAddCNIPolicy:          forceAddCNIPolicy,
		SkipEgressRules:            skipEgressRules,
		DisableAccessEntryCreation: disableAccessEntryCreation,
		VPCImporter:                vpcImporter,
		Parallelism:                parallelism,
	})
}

// NewManagedNodeGroupTask defines tasks required to create managed nodegroups
func (c *StackCollection) NewManagedNodeGroupTask(ctx context.Context, nodeGroups []*api.ManagedNodeGroup, forceAddCNIPolicy bool, vpcImporter vpc.Importer, nodeGroupParallelism int) *tasks.TaskTree {
	taskTree := &tasks.TaskTree{Parallel: true, Limit: nodeGroupParallelism}
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
						return fmt.Errorf("failed to create service account %s/%s: %w", objectMeta.GetNamespace(), objectMeta.GetName(), err)
					}
					return nil
				},
			})
		}

		taskTree.Append(saTasks)
	}
	return taskTree
}
