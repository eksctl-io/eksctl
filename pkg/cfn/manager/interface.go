package manager

import (
	"context"

	asgtypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	cttypes "github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/weaveworks/eksctl/pkg/actions/accessentry"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

// UpdateStackOptions contains options for updating a stack.
type UpdateStackOptions struct {
	Stack         *Stack
	StackName     string
	ChangeSetName string
	Description   string
	TemplateData  TemplateData
	Parameters    map[string]string
	Wait          bool
}

// GetNodegroupOption nodegroup options.
type GetNodegroupOption struct {
	Stack         *NodeGroupStack
	NodeGroupName string
}

var _ StackManager = &StackCollection{}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_stack_manager.go . StackManager
type StackManager interface {
	AppendNewClusterStackResource(ctx context.Context, extendForOutposts, plan bool) (bool, error)
	CreateStack(ctx context.Context, name string, stack builder.ResourceSetReader, tags, parameters map[string]string, errs chan error) error
	DeleteStackBySpec(ctx context.Context, s *Stack) (*Stack, error)
	DeleteStackBySpecSync(ctx context.Context, s *Stack, errs chan error) error
	DeleteStackSync(ctx context.Context, s *Stack) error
	DeleteTasksForDeprecatedStacks(ctx context.Context) (*tasks.TaskTree, error)
	DescribeClusterStackIfExists(ctx context.Context) (*Stack, error)
	DescribeClusterStack(ctx context.Context) (*Stack, error)
	DescribeIAMServiceAccountStacks(ctx context.Context) ([]*Stack, error)
	DescribeNodeGroupStack(ctx context.Context, nodeGroupName string) (*Stack, error)
	DescribeNodeGroupStacksAndResources(ctx context.Context) (map[string]StackInfo, error)
	DescribeStack(ctx context.Context, i *Stack) (*Stack, error)
	DescribeStackChangeSet(ctx context.Context, i *Stack, changeSetName string) (*ChangeSet, error)
	DescribeStackEvents(ctx context.Context, i *Stack) ([]cfntypes.StackEvent, error)
	DoCreateStackRequest(ctx context.Context, i *Stack, templateData TemplateData, tags, parameters map[string]string, withIAM bool, withNamedIAM bool) error
	DoWaitUntilStackIsCreated(ctx context.Context, i *Stack) error
	EnsureMapPublicIPOnLaunchEnabled(ctx context.Context) error
	FixClusterCompatibility(ctx context.Context) error
	ClusterHasDedicatedVPC(ctx context.Context) (bool, error)
	GetAutoScalingGroupDesiredCapacity(ctx context.Context, name string) (asgtypes.AutoScalingGroup, error)
	GetAutoScalingGroupName(ctx context.Context, s *Stack) (string, error)
	GetClusterStackIfExists(ctx context.Context) (*Stack, error)
	GetFargateStack(ctx context.Context) (*Stack, error)
	GetIAMAddonName(s *Stack) string
	GetIAMAddonsStacks(ctx context.Context) ([]*Stack, error)
	GetIAMServiceAccounts(ctx context.Context) ([]*api.ClusterIAMServiceAccount, error)
	GetKarpenterStack(ctx context.Context) (*Stack, error)
	GetManagedNodeGroupTemplate(ctx context.Context, options GetNodegroupOption) (string, error)
	GetNodeGroupName(s *Stack) string
	GetNodeGroupStackType(ctx context.Context, options GetNodegroupOption) (api.NodeGroupType, error)
	GetStackTemplate(ctx context.Context, stackName string) (string, error)
	GetUnmanagedNodeGroupAutoScalingGroupName(ctx context.Context, s *Stack) (string, error)
	HasClusterStackFromList(ctx context.Context, clusterStackNames []string, clusterName string) (bool, error)
	ListClusterStackNames(ctx context.Context) ([]string, error)
	ListIAMServiceAccountStacks(ctx context.Context) ([]string, error)
	ListNodeGroupStacks(ctx context.Context) ([]*Stack, error)
	ListNodeGroupStacksWithStatuses(ctx context.Context) ([]NodeGroupStack, error)
	ListPodIdentityStackNames(ctx context.Context) ([]string, error)
	ListStacks(ctx context.Context) ([]*Stack, error)
	ListStacksWithStatuses(ctx context.Context, statusFilters ...cfntypes.StackStatus) ([]*Stack, error)
	ListStacksMatching(ctx context.Context, nameRegex string, statusFilters ...cfntypes.StackStatus) ([]*Stack, error)
	ListStackNames(ctx context.Context, regExp string) ([]string, error)
	ListAccessEntryStackNames(ctx context.Context, clusterName string) ([]string, error)
	LookupCloudTrailEvents(ctx context.Context, i *Stack) ([]cttypes.Event, error)
	MakeChangeSetName(action string) string
	MakeClusterStackName() string
	NewManagedNodeGroupTask(ctx context.Context, nodeGroups []*api.ManagedNodeGroup, forceAddCNIPolicy bool, importer vpc.Importer) *tasks.TaskTree
	NewTasksToDeleteClusterWithNodeGroups(ctx context.Context, clusterStack *Stack, nodeGroupStacks []NodeGroupStack, clusterOperable bool, newOIDCManager NewOIDCManager, newTasksToDeleteAddonIAM NewTasksToDeleteAddonIAM, newTasksToDeletePodIdentityRole NewTasksToDeletePodIdentityRole, cluster *ekstypes.Cluster, clientSetGetter kubernetes.ClientSetGetter, wait, force bool, cleanup func(chan error, string) error) (*tasks.TaskTree, error)
	NewTasksToCreateIAMServiceAccounts(serviceAccounts []*api.ClusterIAMServiceAccount, oidc *iamoidc.OpenIDConnectManager, clientSetGetter kubernetes.ClientSetGetter) *tasks.TaskTree
	NewTaskToDeleteUnownedNodeGroup(ctx context.Context, clusterName, nodegroup string, nodeGroupDeleter NodeGroupDeleter, waitCondition *DeleteWaitCondition) tasks.Task
	NewTasksToCreateCluster(ctx context.Context, nodeGroups []*api.NodeGroup, managedNodeGroups []*api.ManagedNodeGroup, accessConfig *api.AccessConfig, accessEntryCreator accessentry.CreatorInterface, postClusterCreationTasks ...tasks.Task) *tasks.TaskTree
	NewTasksToDeleteIAMServiceAccounts(ctx context.Context, serviceAccounts []string, clientSetGetter kubernetes.ClientSetGetter, wait bool) (*tasks.TaskTree, error)
	NewTasksToDeleteNodeGroups(stacks []NodeGroupStack, shouldDelete func(_ string) bool, wait bool, cleanup func(chan error, string) error) (*tasks.TaskTree, error)
	NewTasksToDeleteOIDCProviderWithIAMServiceAccounts(ctx context.Context, newOIDCManager NewOIDCManager, cluster *ekstypes.Cluster, clientSetGetter kubernetes.ClientSetGetter, force bool) (*tasks.TaskTree, error)
	NewUnmanagedNodeGroupTask(ctx context.Context, nodeGroups []*api.NodeGroup, forceAddCNIPolicy, skipEgressRules, disableAccessEntryCreation bool, importer vpc.Importer) *tasks.TaskTree
	PropagateManagedNodeGroupTagsToASG(ngName string, ngTags map[string]string, asgNames []string, errCh chan error) error
	RefreshFargatePodExecutionRoleARN(ctx context.Context) error
	StackStatusIsNotTransitional(s *Stack) bool
	TroubleshootStackFailureCause(ctx context.Context, s *cfntypes.Stack, desiredStatus cfntypes.StackStatus)
	UpdateNodeGroupStack(ctx context.Context, nodeGroupName, template string, wait bool) error
	UpdateStack(ctx context.Context, options UpdateStackOptions) error
	MustUpdateStack(ctx context.Context, options UpdateStackOptions) error
}
