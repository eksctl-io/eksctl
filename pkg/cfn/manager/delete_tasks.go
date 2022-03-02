package manager

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/cfn/waiter"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/spot"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

// think Jake is deleting this soon
func deleteAll(_ string) bool { return true }

type NewOIDCManager func() (*iamoidc.OpenIDConnectManager, error)

// NewTasksToDeleteClusterWithNodeGroups defines tasks required to delete the given cluster along with all of its resources
func (c *StackCollection) NewTasksToDeleteClusterWithNodeGroups(ctx context.Context, clusterStack *Stack, nodeGroupStacks []NodeGroupStack, clusterOperable bool, newOIDCManager NewOIDCManager, cluster *ekstypes.Cluster, clientSetGetter kubernetes.ClientSetGetter, wait, force bool, cleanup func(chan error, string) error) (*tasks.TaskTree, error) {
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

	deleteAddonIAMTasks, err := c.NewTaskToDeleteAddonIAM(ctx, wait)
	if err != nil {
		return nil, err
	}

	if deleteAddonIAMTasks.Len() > 0 {
		deleteAddonIAMTasks.IsSubTask = true
		taskTree.Append(deleteAddonIAMTasks)
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

		if !shouldDelete(s.NodeGroupName) || s.NodeGroupName == api.SpotOceanClusterNodeGroupName {
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

	// Spot Ocean.
	{
		oceanTaskTree, err := c.NewTasksToDeleteSpotOceanNodeGroup(context.TODO(), shouldDelete)
		if err != nil {
			return nil, err
		}
		if oceanTaskTree.Len() > 0 {
			oceanTaskTree.IsSubTask = true
			taskTree.Append(oceanTaskTree)
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
	out, err := d.eksAPI.DeleteNodegroup(d.ctx, &awseks.DeleteNodegroupInput{
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
		if iamoidc.IsAccessDeniedError(err) {
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

// NewTasksToDeleteSpotOceanNodeGroup defines tasks required to delete Ocean nodegroup.
func (c *StackCollection) NewTasksToDeleteSpotOceanNodeGroup(ctx context.Context, shouldDelete func(string) bool) (*tasks.TaskTree, error) {
	taskTree := &tasks.TaskTree{Parallel: true}

	// Check whether the Ocean Cluster's nodegroup should be deleted.
	stacks, err := c.ListNodeGroupStacks(ctx)
	if err != nil {
		return nil, err
	}
	stack, err := spot.ShouldDeleteOceanCluster(stacks, shouldDelete)
	if err != nil {
		return nil, err
	}
	if stack == nil { // nothing to do
		return taskTree, nil
	}

	// ignoreListImportsError ignores errors that may occur while listing imports.
	ignoreListImportsError := func(errMsg string) bool {
		errMsgs := []string{
			"not imported by any stack",
			"does not exist",
		}
		for _, msg := range errMsgs {
			if strings.Contains(strings.ToLower(errMsg), msg) {
				return true
			}
		}
		return false
	}

	// All nodegroups are marked for deletion. We need to wait for their deletion
	// to complete before deleting the Ocean Cluster.
	deleter := func(ctx context.Context, s *Stack, errs chan error) error {
		maxAttempts := 360 // 1 hour
		delay := 10 * time.Second

		for attempt := 1; ; attempt++ {
			logger.Debug("ocean: attempting to delete cluster (attempt: %d)", attempt)

			input := &cloudformation.ListImportsInput{
				ExportName: aws.String(fmt.Sprintf("%s::%s",
					aws.ToString(s.StackName), outputs.NodeGroupSpotOceanClusterID)),
			}

			output, err := c.cloudformationAPI.ListImports(ctx, input)
			if err != nil {
				if !ignoreListImportsError(err.Error()) {
					return err
				}
			}

			if output != nil && len(output.Imports) > 0 {
				if attempt+1 > maxAttempts {
					return fmt.Errorf("ocean: max attempts reached: " +
						"giving up waiting for importers to become deleted")
				}

				logger.Debug("ocean: waiting for %d importers "+
					"to become deleted", len(output.Imports))
				time.Sleep(delay)
				continue
			}

			logger.Debug("ocean: no more importers, deleting cluster...")
			return c.DeleteStackBySpecSync(ctx, s, errs)
		}
	}

	// Add a new deletion task.
	taskTree.Append(&taskWithStackSpec{
		info:  "delete ocean cluster",
		stack: stack,
		call:  deleter,
	})

	return taskTree, nil
}

func clusterHasOIDCProvider(cluster *ekstypes.Cluster) (hasOIDC bool, found bool) {
	for k, v := range cluster.Tags {
		if k == api.ClusterOIDCEnabledTag {
			return v == "true", true
		}
	}
	return false, false
}
