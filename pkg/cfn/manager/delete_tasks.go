package manager

import (
	"context"
	"fmt"
	"time"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	"github.com/weaveworks/eksctl/pkg/actions/accessentry"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/cfn/waiter"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/utils/apierrors"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

// think Jake is deleting this soon
func deleteAll(_ string) bool { return true }

type NewOIDCManager func() (*iamoidc.OpenIDConnectManager, error)

// NewTasksToDeleteAddonIAM temporary type, to be removed after moving NewTasksToDeleteClusterWithNodeGroups to actions package
type NewTasksToDeleteAddonIAM func(ctx context.Context, wait bool) (*tasks.TaskTree, error)

// NewTasksToDeletePodIdentityRoles temporary type, to be removed after moving NewTasksToDeleteClusterWithNodeGroups to actions package
type NewTasksToDeletePodIdentityRole func() (*tasks.TaskTree, error)

// NewTasksToDeleteClusterWithNodeGroups defines tasks required to delete the given cluster along with all of its resources
func (c *StackCollection) NewTasksToDeleteClusterWithNodeGroups(
	ctx context.Context,
	clusterStack *Stack,
	nodeGroupStacks []NodeGroupStack,
	clusterOperable bool,
	newOIDCManager NewOIDCManager,
	newTasksToDeleteAddonIAM NewTasksToDeleteAddonIAM,
	newTasksToDeletePodIdentityRole NewTasksToDeletePodIdentityRole,
	cluster *ekstypes.Cluster,
	clientSetGetter kubernetes.ClientSetGetter,
	wait, force bool,
	cleanup func(chan error, string) error) (*tasks.TaskTree, error) {
	taskTree := &tasks.TaskTree{Parallel: false}

	nodeGroupTasks, err := c.NewTasksToDeleteNodeGroups(nodeGroupStacks, deleteAll, true, cleanup)

	if err != nil {
		return nil, err
	}
	if nodeGroupTasks.Len() > 0 {
		nodeGroupTasks.IsSubTask = true
		taskTree.Append(nodeGroupTasks)
	}

	if clusterOperable {
		serviceAccountAndOIDCTasks, err := c.NewTasksToDeleteOIDCProviderWithIAMServiceAccounts(ctx, newOIDCManager, cluster, clientSetGetter, force)
		if err != nil {
			return nil, err
		}

		if serviceAccountAndOIDCTasks.Len() > 0 {
			serviceAccountAndOIDCTasks.IsSubTask = true
			taskTree.Append(serviceAccountAndOIDCTasks)
		}
	}

	deleteAddonIAMTasks, err := newTasksToDeleteAddonIAM(ctx, wait)
	if err != nil {
		return nil, err
	}

	if deleteAddonIAMTasks.Len() > 0 {
		deleteAddonIAMTasks.IsSubTask = true
		taskTree.Append(deleteAddonIAMTasks)
	}

	deletePodIdentityRoleTasks, err := newTasksToDeletePodIdentityRole()
	if err != nil {
		return nil, err
	}
	if deletePodIdentityRoleTasks.Len() > 0 {
		deletePodIdentityRoleTasks.IsSubTask = true
		taskTree.Append(deletePodIdentityRoleTasks)
	}

	deleteAccessEntriesTasks, err := accessentry.
		NewRemover(c.spec.Metadata.Name, c, c.eksAPI).
		DeleteTasks(ctx, []api.AccessEntry{})
	if err != nil {
		return nil, err
	}
	if deleteAccessEntriesTasks.Len() > 0 {
		deleteAccessEntriesTasks.IsSubTask = true
		taskTree.Append(deleteAccessEntriesTasks)
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

func usesAccessEntry(stack *Stack) bool {
	for _, output := range stack.Outputs {
		if *output.OutputKey == outputs.NodeGroupUsesAccessEntry {
			return *output.OutputValue == "true"
		}
	}
	return false
}

type DeleteWaitCondition struct {
	Condition func() (bool, error)
	Timeout   time.Duration
	Interval  time.Duration
}

//counterfeiter:generate -o fakes/fake_nodegroup_deleter.go . NodeGroupDeleter
type NodeGroupDeleter interface {
	DeleteNodegroup(ctx context.Context, params *awseks.DeleteNodegroupInput, optFns ...func(*awseks.Options)) (*awseks.DeleteNodegroupOutput, error)
}

type DeleteUnownedNodegroupTask struct {
	cluster          string
	nodegroup        string
	wait             *DeleteWaitCondition
	info             string
	nodeGroupDeleter NodeGroupDeleter
	ctx              context.Context
}

func (d *DeleteUnownedNodegroupTask) Describe() string {
	return d.info
}

func (d *DeleteUnownedNodegroupTask) Do() error {
	out, err := d.nodeGroupDeleter.DeleteNodegroup(d.ctx, &awseks.DeleteNodegroupInput{
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

func (c *StackCollection) NewTaskToDeleteUnownedNodeGroup(ctx context.Context, clusterName, nodegroup string, nodeGroupDeleter NodeGroupDeleter, waitCondition *DeleteWaitCondition) tasks.Task {
	return tasks.SynchronousTask{
		SynchronousTaskIface: &DeleteUnownedNodegroupTask{
			cluster:          clusterName,
			nodegroup:        nodegroup,
			nodeGroupDeleter: nodeGroupDeleter,
			wait:             waitCondition,
			info:             fmt.Sprintf("delete unowned nodegroup %s", nodegroup),
			ctx:              ctx,
		}}
}

// NewTasksToDeleteOIDCProviderWithIAMServiceAccounts defines tasks required to delete all of the iamserviceaccounts
// along with associated IAM OIDC provider
func (c *StackCollection) NewTasksToDeleteOIDCProviderWithIAMServiceAccounts(ctx context.Context, newOIDCManager NewOIDCManager, cluster *ekstypes.Cluster, clientSetGetter kubernetes.ClientSetGetter, force bool) (*tasks.TaskTree, error) {
	taskTree := &tasks.TaskTree{Parallel: false}

	oidc, err := newOIDCManager()
	if err != nil {
		if _, ok := err.(*iamoidc.UnsupportedOIDCError); ok {
			logger.Debug("OIDC is not supported for this cluster")
			return taskTree, nil
		}
		return nil, fmt.Errorf("error creating OIDC manager: %w", err)
	}

	allServiceAccountsWithStacks, err := c.getAllServiceAccounts(ctx)
	if err != nil {
		return nil, err
	}

	if len(allServiceAccountsWithStacks) > 0 {
		saTasks, err := c.NewTasksToDeleteIAMServiceAccounts(ctx, allServiceAccountsWithStacks, clientSetGetter, true)
		if err != nil {
			return nil, err
		}

		if saTasks.Len() > 0 {
			saTasks.IsSubTask = true
			taskTree.Append(saTasks)
		}
	}

	providerExists, err := oidc.CheckProviderExists(ctx)
	if err != nil {
		if apierrors.IsAccessDeniedError(err) {
			clusterHasOIDC, foundTag := clusterHasOIDCProvider(cluster)
			errMsg := "IAM permissions are required to delete OIDC provider"
			switch {
			case foundTag:
				if clusterHasOIDC {
					return nil, fmt.Errorf("%s: %w", errMsg, err)
				}
				if len(allServiceAccountsWithStacks) > 0 {
					logger.Warning("expected an OIDC provider to be associated with the cluster; found %d service account(s)", len(allServiceAccountsWithStacks))
				}
			case len(allServiceAccountsWithStacks) > 0:
				if !force {
					return nil, fmt.Errorf("found %d IAM service account(s); %s: %w", len(allServiceAccountsWithStacks), errMsg, err)
				}
			default:
				logger.Info("could not determine if cluster has an OIDC provider because of missing IAM permissions; " +
					"if an OIDC provider was associated with the cluster (either by setting `iam.withOIDC: true` or by using `eksctl utils associate-iam-oidc-provider`), " +
					"run `aws iam delete-open-id-connect-provider` from an authorized IAM entity to delete it")
			}
			return taskTree, nil
		}
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

func clusterHasOIDCProvider(cluster *ekstypes.Cluster) (hasOIDC bool, found bool) {
	for k, v := range cluster.Tags {
		if k == api.ClusterOIDCEnabledTag {
			return v == "true", true
		}
	}
	return false, false
}
