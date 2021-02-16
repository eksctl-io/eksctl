package manager

import (
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

//go:generate counterfeiter -o fakes/fake_stack_manager.go . StackManager
type StackManager interface {
	ListNodeGroupStacks() ([]NodeGroupStack, error)
	DescribeNodeGroupStacksAndResources() (map[string]StackInfo, error)
	ScaleNodeGroup(ng *v1alpha5.NodeGroup) error
	GetNodeGroupSummaries(name string) ([]*NodeGroupSummary, error)
	GetNodeGroupAutoScalingGroupName(s *Stack) (string, error)
	GetManagedNodeGroupAutoScalingGroupName(s *Stack) (string, error)
	DescribeNodeGroupStack(nodeGroupName string) (*Stack, error)
	DescribeNodeGroupStacks() ([]*Stack, error)
	GetNodeGroupStackType(name string) (v1alpha5.NodeGroupType, error)
	GetNodeGroupName(s *Stack) string
	DoWaitUntilStackIsCreated(i *Stack) error
	DoCreateStackRequest(i *Stack, templateData TemplateData, tags, parameters map[string]string, withIAM bool, withNamedIAM bool) error
	CreateStack(name string, stack builder.ResourceSet, tags, parameters map[string]string, errs chan error) error
	UpdateStack(stackName, changeSetName, description string, templateData TemplateData, parameters map[string]string) error
	DescribeStack(i *Stack) (*Stack, error)
	GetManagedNodeGroupTemplate(nodeGroupName string) (string, error)
	UpdateNodeGroupStack(nodeGroupName, template string) error
	ListStacksMatching(nameRegex string, statusFilters ...string) ([]*Stack, error)
	ListClusterStackNames() ([]string, error)
	ListStacks(statusFilters ...string) ([]*Stack, error)
	StackStatusIsNotTransitional(s *Stack) bool
	StackStatusIsNotReady(s *Stack) bool
	DeleteStackByName(name string) (*Stack, error)
	DeleteStackByNameSync(name string) error
	DeleteStackBySpec(s *Stack) (*Stack, error)
	DeleteStackBySpecSync(s *Stack, errs chan error) error
	DescribeStacks() ([]*Stack, error)
	HasClusterStack() (bool, error)
	HasClusterStackUsingCachedList(clusterStackNames []string) (bool, error)
	DescribeStackEvents(i *Stack) ([]*cloudformation.StackEvent, error)
	LookupCloudTrailEvents(i *Stack) ([]*cloudtrail.Event, error)
	DescribeStackChangeSet(i *Stack, changeSetName string) (*ChangeSet, error)
	MakeChangeSetName(action string) string
	DescribeClusterStack() (*Stack, error)
	RefreshFargatePodExecutionRoleARN() error
	AppendNewClusterStackResource(plan, supportsManagedNodes bool) (bool, error)
	GetFargateStack() (*Stack, error)
	GetStackTemplate(stackName string) (string, error)
	NewTasksToCreateClusterWithNodeGroups(nodeGroups []*v1alpha5.NodeGroup,
		managedNodeGroups []*v1alpha5.ManagedNodeGroup, supportsManagedNodes bool, postClusterCreationTasks ...tasks.Task) *tasks.TaskTree
	NewUnmanagedNodeGroupTask(nodeGroups []*v1alpha5.NodeGroup, supportsManagedNodes bool, forceAddCNIPolicy bool) *tasks.TaskTree
	NewManagedNodeGroupTask(nodeGroups []*v1alpha5.ManagedNodeGroup, forceAddCNIPolicy bool) *tasks.TaskTree
	NewClusterCompatTask() tasks.Task
	NewTasksToCreateIAMServiceAccounts(serviceAccounts []*v1alpha5.ClusterIAMServiceAccount, oidc *iamoidc.OpenIDConnectManager, clientSetGetter kubernetes.ClientSetGetter) *tasks.TaskTree
	DeleteTasksForDeprecatedStacks() (*tasks.TaskTree, error)
	NewTasksToDeleteClusterWithNodeGroups(deleteOIDCProvider bool, oidc *iamoidc.OpenIDConnectManager, clientSetGetter kubernetes.ClientSetGetter, wait bool, cleanup func(chan error, string) error) (*tasks.TaskTree, error)
	NewTasksToDeleteNodeGroups(shouldDelete func(string) bool, wait bool, cleanup func(chan error, string) error) (*tasks.TaskTree, error)
	NewTasksToDeleteOIDCProviderWithIAMServiceAccounts(oidc *iamoidc.OpenIDConnectManager, clientSetGetter kubernetes.ClientSetGetter) (*tasks.TaskTree, error)
	NewTasksToDeleteIAMServiceAccounts(shouldDelete func(string) bool, clientSetGetter kubernetes.ClientSetGetter, wait bool) (*tasks.TaskTree, error)
	NewTaskToDeleteAddonIAM(wait bool) (*tasks.TaskTree, error)
	FixClusterCompatibility() error
	DescribeIAMServiceAccountStacks() ([]*Stack, error)
	ListIAMServiceAccountStacks() ([]string, error)
	GetIAMServiceAccounts() ([]*v1alpha5.ClusterIAMServiceAccount, error)
	GetIAMAddonsStacks() ([]*Stack, error)
	GetIAMAddonName(s *Stack) string
	EnsureMapPublicIPOnLaunchEnabled() error
	GetAutoScalingGroupName(s *Stack) (string, error)
}
