package manager

import (
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"

	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
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

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_stack_manager.go . StackManager
type StackManager interface {
	AppendNewClusterStackResource(plan, supportsManagedNodes bool) (bool, error)
	CreateStack(name string, stack builder.ResourceSet, tags, parameters map[string]string, errs chan error) error
	DeleteStackBySpec(s *Stack) (*Stack, error)
	DeleteStackBySpecSync(s *Stack, errs chan error) error
	DeleteStackSync(s *Stack) error
	DeleteTasksForDeprecatedStacks() (*tasks.TaskTree, error)
	DescribeClusterStack() (*Stack, error)
	DescribeIAMServiceAccountStacks() ([]*Stack, error)
	DescribeNodeGroupStack(nodeGroupName string) (*Stack, error)
	DescribeNodeGroupStacks() ([]*Stack, error)
	DescribeNodeGroupStacksAndResources() (map[string]StackInfo, error)
	DescribeStack(i *Stack) (*Stack, error)
	DescribeStackChangeSet(i *Stack, changeSetName string) (*ChangeSet, error)
	DescribeStackEvents(i *Stack) ([]*cloudformation.StackEvent, error)
	DescribeStacks() ([]*Stack, error)
	DoCreateStackRequest(i *Stack, templateData TemplateData, tags, parameters map[string]string, withIAM bool, withNamedIAM bool) error
	DoWaitUntilStackIsCreated(i *Stack) error
	EnsureMapPublicIPOnLaunchEnabled() error
	FixClusterCompatibility() error
	GetAutoScalingGroupName(s *Stack) (string, error)
	GetClusterStackIfExists() (*Stack, error)
	GetFargateStack() (*Stack, error)
	GetIAMAddonName(s *Stack) string
	GetIAMAddonsStacks() ([]*Stack, error)
	GetIAMServiceAccounts() ([]*v1alpha5.ClusterIAMServiceAccount, error)
	GetKarpenterStack() (*Stack, error)
	GetManagedNodeGroupTemplate(options GetNodegroupOption) (string, error)
	GetNodeGroupName(s *Stack) string
	GetNodeGroupStackType(options GetNodegroupOption) (v1alpha5.NodeGroupType, error)
	GetStackTemplate(stackName string) (string, error)
	GetUnmanagedNodeGroupSummaries(name string) ([]*NodeGroupSummary, error)
	HasClusterStackUsingCachedList(clusterStackNames []string, clusterName string) (bool, error)
	ListClusterStackNames() ([]string, error)
	ListIAMServiceAccountStacks() ([]string, error)
	ListNodeGroupStacks() ([]NodeGroupStack, error)
	ListStacks(statusFilters ...string) ([]*Stack, error)
	ListStacksMatching(nameRegex string, statusFilters ...string) ([]*Stack, error)
	LookupCloudTrailEvents(i *Stack) ([]*cloudtrail.Event, error)
	MakeChangeSetName(action string) string
	MakeClusterStackName() string
	NewClusterCompatTask() tasks.Task
	NewManagedNodeGroupTask(nodeGroups []*v1alpha5.ManagedNodeGroup, forceAddCNIPolicy bool, importer vpc.Importer) *tasks.TaskTree
	NewTaskToDeleteAddonIAM(wait bool) (*tasks.TaskTree, error)
	NewTaskToDeleteUnownedNodeGroup(clusterName, nodegroup string, eksAPI eksiface.EKSAPI, waitCondition *DeleteWaitCondition) tasks.Task
	NewTasksToCreateClusterWithNodeGroups(nodeGroups []*v1alpha5.NodeGroup, managedNodeGroups []*v1alpha5.ManagedNodeGroup, supportsManagedNodes bool, postClusterCreationTasks ...tasks.Task) *tasks.TaskTree
	NewTasksToCreateIAMServiceAccounts(serviceAccounts []*v1alpha5.ClusterIAMServiceAccount, oidc *iamoidc.OpenIDConnectManager, clientSetGetter kubernetes.ClientSetGetter) *tasks.TaskTree
	NewTasksToDeleteClusterWithNodeGroups(stack *Stack, stacks []NodeGroupStack, deleteOIDCProvider bool, oidc *iamoidc.OpenIDConnectManager, clientSetGetter kubernetes.ClientSetGetter, wait bool, cleanup func(chan error, string) error) (*tasks.TaskTree, error)
	NewTasksToDeleteIAMServiceAccounts(serviceAccounts []string, clientSetGetter kubernetes.ClientSetGetter, wait bool) (*tasks.TaskTree, error)
	NewTasksToDeleteNodeGroups(stacks []NodeGroupStack, shouldDelete func(_ string) bool, wait bool, cleanup func(chan error, string) error) (*tasks.TaskTree, error)
	NewTasksToDeleteOIDCProviderWithIAMServiceAccounts(oidc *iamoidc.OpenIDConnectManager, clientSetGetter kubernetes.ClientSetGetter) (*tasks.TaskTree, error)
	NewUnmanagedNodeGroupTask(nodeGroups []*v1alpha5.NodeGroup, forceAddCNIPolicy bool, importer vpc.Importer) *tasks.TaskTree
	RefreshFargatePodExecutionRoleARN() error
	StackStatusIsNotReady(s *Stack) bool
	StackStatusIsNotTransitional(s *Stack) bool
	UpdateNodeGroupStack(nodeGroupName, template string, wait bool) error
	UpdateStack(options UpdateStackOptions) error
}
