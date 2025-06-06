package manager

import (
	"context"
	"errors"
	"fmt"
	"strings"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	asgtypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	cfn "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/eks"

	"github.com/blang/semver/v4"
	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
	"github.com/weaveworks/eksctl/pkg/version"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

// NodeGroupStack represents a nodegroup and its type
type NodeGroupStack struct {
	NodeGroupName   string
	Type            api.NodeGroupType
	UsesAccessEntry bool
	Stack           *Stack
}

// makeNodeGroupStackName generates the name of the nodegroup stack identified by its name.
func makeNodeGroupStackName(clusterName, ngName string) string {
	return fmt.Sprintf("eksctl-%s-nodegroup-%s", clusterName, ngName)
}

// CreateNodeGroupOptions holds options for creating nodegroup tasks.
type CreateNodeGroupOptions struct {
	ForceAddCNIPolicy          bool
	SkipEgressRules            bool
	DisableAccessEntryCreation bool
	VPCImporter                vpc.Importer
	Parallelism                int
}

// A NodeGroupStackManager describes and creates nodegroup stacks.
type NodeGroupStackManager interface {
	// CreateStack creates a CloudFormation stack.
	CreateStack(ctx context.Context, stackName string, resourceSet builder.ResourceSetReader, tags, parameters map[string]string, errs chan error) error
}

// A NodeGroupResourceSet creates resources for a nodegroup.
//
//counterfeiter:generate -o fakes/fake_nodegroup_resource_set.go . NodeGroupResourceSet
type NodeGroupResourceSet interface {
	// AddAllResources adds all nodegroup resources.
	AddAllResources(ctx context.Context) error
	builder.ResourceSetReader
}

// CreateNodeGroupResourceSetFunc creates a new NodeGroupResourceSet.
type CreateNodeGroupResourceSetFunc func(options builder.NodeGroupOptions) NodeGroupResourceSet

// NewBootstrapperFunc creates a new Bootstrapper for ng.
type NewBootstrapperFunc func(clusterConfig *api.ClusterConfig, ng *api.NodeGroup) (nodebootstrap.Bootstrapper, error)

// UnmanagedNodeGroupTask creates tasks for creating self-managed nodegroups.
type UnmanagedNodeGroupTask struct {
	ClusterConfig              *api.ClusterConfig
	NodeGroups                 []*api.NodeGroup
	CreateNodeGroupResourceSet CreateNodeGroupResourceSetFunc
	NewBootstrapper            NewBootstrapperFunc
	EKSAPI                     awsapi.EKS
	StackManager               NodeGroupStackManager
}

// Create creates a TaskTree for creating nodegroups.
func (t *UnmanagedNodeGroupTask) Create(ctx context.Context, options CreateNodeGroupOptions) *tasks.TaskTree {
	taskTree := &tasks.TaskTree{Parallel: true, Limit: options.Parallelism}

	for _, ng := range t.NodeGroups {
		ng := ng
		createAccessEntryInStack := ng.IAM.InstanceRoleARN == ""
		createNodeGroupTask := &tasks.GenericTask{
			Description: fmt.Sprintf("create nodegroup %q", ng.NameString()),
			Doer: func() error {
				return t.createNodeGroup(ctx, ng, options, createAccessEntryInStack)
			},
		}

		if options.DisableAccessEntryCreation || createAccessEntryInStack {
			taskTree.Append(createNodeGroupTask)
		} else {
			var ngTask tasks.TaskTree
			ngTask.Append(createNodeGroupTask)
			ngTask.Append(&tasks.GenericTask{
				Description: fmt.Sprintf("create access entry for nodegroup %q", ng.NameString()),
				Doer: func() error {
					return t.maybeCreateAccessEntry(ctx, ng)
				},
			})
			taskTree.Append(&ngTask)
		}
	}

	return taskTree
}

func (t *UnmanagedNodeGroupTask) createNodeGroup(ctx context.Context, ng *api.NodeGroup, options CreateNodeGroupOptions, createAccessEntryInStack bool) error {
	name := makeNodeGroupStackName(t.ClusterConfig.Metadata.Name, ng.Name)

	logger.Info("building nodegroup stack %q", name)
	bootstrapper, err := t.NewBootstrapper(t.ClusterConfig, ng)
	if err != nil {
		return fmt.Errorf("error creating bootstrapper: %w", err)
	}

	resourceSet := t.CreateNodeGroupResourceSet(builder.NodeGroupOptions{
		ClusterConfig:              t.ClusterConfig,
		NodeGroup:                  ng,
		Bootstrapper:               bootstrapper,
		ForceAddCNIPolicy:          options.ForceAddCNIPolicy,
		VPCImporter:                options.VPCImporter,
		SkipEgressRules:            options.SkipEgressRules,
		DisableAccessEntry:         options.DisableAccessEntryCreation,
		DisableAccessEntryResource: !createAccessEntryInStack,
	})
	if err := resourceSet.AddAllResources(ctx); err != nil {
		return err
	}

	if ng.Tags == nil {
		ng.Tags = make(map[string]string)
	}
	ng.Tags[api.NodeGroupNameTag] = ng.Name
	ng.Tags[api.OldNodeGroupNameTag] = ng.Name
	ng.Tags[api.NodeGroupTypeTag] = string(api.NodeGroupTypeUnmanaged)

	errCh := make(chan error)
	if err := t.StackManager.CreateStack(ctx, name, resourceSet, ng.Tags, nil, errCh); err != nil {
		return err
	}
	return <-errCh
}

func (t *UnmanagedNodeGroupTask) maybeCreateAccessEntry(ctx context.Context, ng *api.NodeGroup) error {
	roleARN := ng.IAM.InstanceRoleARN
	_, err := t.EKSAPI.CreateAccessEntry(ctx, &eks.CreateAccessEntryInput{
		ClusterName:  aws.String(t.ClusterConfig.Metadata.Name),
		PrincipalArn: aws.String(roleARN),
		Type:         aws.String(string(api.GetAccessEntryType(ng))),
		Tags: map[string]string{
			api.ClusterNameLabel: t.ClusterConfig.Metadata.Name,
		},
	})
	if err != nil {
		var resourceInUse *ekstypes.ResourceInUseException
		if errors.As(err, &resourceInUse) {
			logger.Info("nodegroup %s: access entry for principal ARN %q already exists", ng.Name, roleARN)
			return nil
		}
		return fmt.Errorf("creating access entry for nodegroup %s: %w", ng.Name, err)
	}
	logger.Info("nodegroup %s: created access entry for principal ARN %q", ng.Name, roleARN)
	return nil
}

// makeNodeGroupStackName generates the name of the nodegroup stack identified by its name and this StackCollection's cluster.
func (c *StackCollection) makeNodeGroupStackName(ngName string) string {
	return makeNodeGroupStackName(c.spec.Metadata.Name, ngName)
}

func (c *StackCollection) createManagedNodeGroupTask(ctx context.Context, errorCh chan error, ng *api.ManagedNodeGroup, forceAddCNIPolicy bool, vpcImporter vpc.Importer) error {
	name := c.makeNodeGroupStackName(ng.Name)
	cluster, err := c.DescribeClusterStackIfExists(ctx)
	if err != nil {
		return err
	}
	if cluster == nil && c.spec.IPv6Enabled() {
		return errors.New("managed nodegroups cannot be created on IPv6 unowned clusters")
	}
	logger.Info("building managed nodegroup stack %q", name)
	bootstrapper, err := nodebootstrap.NewManagedBootstrapper(c.spec, ng)
	if err != nil {
		return err
	}
	stack := builder.NewManagedNodeGroup(c.ec2API, c.spec, ng, builder.NewLaunchTemplateFetcher(c.ec2API), bootstrapper, forceAddCNIPolicy, vpcImporter)
	if err := stack.AddAllResources(ctx); err != nil {
		return err
	}

	return c.CreateStack(ctx, name, stack, ng.Tags, nil, errorCh)
}

func (c *StackCollection) propagateManagedNodeGroupTagsToASGTask(ctx context.Context, errorCh chan error, ng *api.ManagedNodeGroup,
	propagateFunc func(string, map[string]string, []string, chan error) error) error {
	// describe node group to retrieve ASG names
	input := &eks.DescribeNodegroupInput{
		ClusterName:   aws.String(c.spec.Metadata.Name),
		NodegroupName: aws.String(ng.Name),
	}
	res, err := c.eksAPI.DescribeNodegroup(ctx, input)
	if err != nil {
		return fmt.Errorf("couldn't get managed nodegroup details for nodegroup %q: %w", ng.Name, err)
	}

	if res.Nodegroup.Resources == nil {
		return nil
	}

	asgNames := []string{}
	for _, asg := range res.Nodegroup.Resources.AutoScalingGroups {
		if asg.Name != nil && *asg.Name != "" {
			asgNames = append(asgNames, *asg.Name)
		}
	}

	// add labels and taints
	tags := map[string]string{}
	builder.GenerateClusterAutoscalerTags(ng, func(key, value string) {
		tags[key] = value
	})

	// add nodegroup tags
	for k, v := range ng.Tags {
		tags[k] = v
	}

	return propagateFunc(ng.Name, tags, asgNames, errorCh)
}

// ListNodeGroupStacks calls ListStacks and filters out nodegroups
func (c *StackCollection) ListNodeGroupStacks(ctx context.Context) ([]*Stack, error) {
	stacks, err := c.ListStacks(ctx)
	if err != nil {
		return nil, err
	}

	if len(stacks) == 0 {
		return nil, nil
	}

	nodeGroupStacks := []*Stack{}
	for _, s := range stacks {
		switch s.StackStatus {
		case types.StackStatusDeleteComplete:
			continue
		case types.StackStatusDeleteFailed:
			logger.Warning("stack's status of nodegroup named %s is %s", *s.StackName, s.StackStatus)
			continue
		}
		if c.GetNodeGroupName(s) != "" {
			nodeGroupStacks = append(nodeGroupStacks, s)
		}
	}
	logger.Debug("nodegroups = %v", nodeGroupStacks)
	return nodeGroupStacks, nil
}

// ListNodeGroupStacksWithStatuses returns a list of NodeGroupStacks.
func (c *StackCollection) ListNodeGroupStacksWithStatuses(ctx context.Context) ([]NodeGroupStack, error) {
	stacks, err := c.ListNodeGroupStacks(ctx)
	if err != nil {
		return nil, err
	}
	var nodeGroupStacks []NodeGroupStack
	for _, stack := range stacks {
		nodeGroupType, err := GetNodeGroupType(stack.Tags)
		if err != nil {
			return nil, err
		}
		nodeGroupStacks = append(nodeGroupStacks, NodeGroupStack{
			NodeGroupName:   c.GetNodeGroupName(stack),
			Type:            nodeGroupType,
			UsesAccessEntry: nodeGroupType == api.NodeGroupTypeUnmanaged && usesAccessEntry(stack),
			Stack:           stack,
		})
	}
	return nodeGroupStacks, nil
}

// DescribeNodeGroupStacksAndResources calls DescribeNodeGroupStackList and fetches all resources,
// then returns it in a map by nodegroup name
func (c *StackCollection) DescribeNodeGroupStacksAndResources(ctx context.Context) (map[string]StackInfo, error) {
	stacks, err := c.ListNodeGroupStacks(ctx)
	if err != nil {
		return nil, err
	}

	allResources := make(map[string]StackInfo)

	for _, s := range stacks {
		input := &cfn.DescribeStackResourcesInput{
			StackName: s.StackName,
		}
		resources, err := c.cloudformationAPI.DescribeStackResources(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("getting all resources for %q stack: %w", *s.StackName, err)
		}
		allResources[c.GetNodeGroupName(s)] = StackInfo{
			Resources: resources.StackResources,
			Stack:     s,
		}
	}

	return allResources, nil
}

func (c *StackCollection) GetAutoScalingGroupName(ctx context.Context, s *Stack) (string, error) {
	nodeGroupType, err := GetNodeGroupType(s.Tags)
	if err != nil {
		return "", err
	}

	switch nodeGroupType {
	case api.NodeGroupTypeManaged:
		res, err := c.getManagedNodeGroupAutoScalingGroupName(ctx, s)
		if err != nil {
			return "", err
		}
		return res, nil
	case api.NodeGroupTypeUnmanaged, "":
		res, err := c.GetUnmanagedNodeGroupAutoScalingGroupName(ctx, s)
		if err != nil {
			return "", err
		}
		return res, nil

	default:
		return "", fmt.Errorf("cant get autoscaling group name, because unexpected nodegroup type : %q", nodeGroupType)
	}
}

// GetUnmanagedNodeGroupAutoScalingGroupName returns the unmanaged nodegroup's AutoScalingGroupName.
func (c *StackCollection) GetUnmanagedNodeGroupAutoScalingGroupName(ctx context.Context, s *Stack) (string, error) {
	input := &cfn.DescribeStackResourceInput{
		StackName:         s.StackName,
		LogicalResourceId: aws.String("NodeGroup"),
	}

	res, err := c.cloudformationAPI.DescribeStackResource(ctx, input)
	if err != nil {
		return "", err
	}
	if res.StackResourceDetail.PhysicalResourceId == nil {
		return "", fmt.Errorf("%q resource of stack %q has no physical resource id", *input.LogicalResourceId, *res.StackResourceDetail.LogicalResourceId)
	}
	return *res.StackResourceDetail.PhysicalResourceId, nil
}

// GetManagedNodeGroupAutoScalingGroupName returns the managed nodegroup's AutoScalingGroup names
func (c *StackCollection) getManagedNodeGroupAutoScalingGroupName(ctx context.Context, s *Stack) (string, error) {
	input := &eks.DescribeNodegroupInput{
		ClusterName:   aws.String(getClusterNameTag(s)),
		NodegroupName: aws.String(c.GetNodeGroupName(s)),
	}

	res, err := c.eksAPI.DescribeNodegroup(ctx, input)
	if err != nil {
		logger.Warning("couldn't get managed nodegroup details for stack %q", *s.StackName)
		return "", nil
	}

	var asgs []string

	if res.Nodegroup.Resources != nil {
		for _, v := range res.Nodegroup.Resources.AutoScalingGroups {
			asgs = append(asgs, aws.ToString(v.Name))
		}
	}
	return strings.Join(asgs, ","), nil
}

func (c *StackCollection) GetAutoScalingGroupDesiredCapacity(ctx context.Context, name string) (asgtypes.AutoScalingGroup, error) {
	asg, err := c.asgAPI.DescribeAutoScalingGroups(ctx, &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{
			name,
		},
	})

	if err != nil {
		return asgtypes.AutoScalingGroup{}, fmt.Errorf("couldn't describe ASG: %s", name)
	}
	if len(asg.AutoScalingGroups) != 1 {
		logger.Warning("couldn't find ASG %s", name)
		return asgtypes.AutoScalingGroup{}, fmt.Errorf("couldn't find ASG: %s", name)
	}

	return asg.AutoScalingGroups[0], nil
}

// DescribeNodeGroupStack gets the specified nodegroup stack
func (c *StackCollection) DescribeNodeGroupStack(ctx context.Context, nodeGroupName string) (*Stack, error) {
	stackName := c.makeNodeGroupStackName(nodeGroupName)
	return c.DescribeStack(ctx, &Stack{StackName: &stackName})
}

// GetNodeGroupStackType returns the nodegroup stack type
func (c *StackCollection) GetNodeGroupStackType(ctx context.Context, options GetNodegroupOption) (api.NodeGroupType, error) {
	var (
		err   error
		stack *Stack
	)
	if options.Stack != nil && options.Stack.Stack != nil {
		stack = options.Stack.Stack
	}
	if stack == nil {
		stack, err = c.DescribeNodeGroupStack(ctx, options.NodeGroupName)
		if err != nil {
			return "", err
		}
	}
	return GetNodeGroupType(stack.Tags)
}

// GetNodeGroupType returns the nodegroup type
func GetNodeGroupType(tags []types.Tag) (api.NodeGroupType, error) {
	var nodeGroupType api.NodeGroupType

	if ngNameTagValue := GetNodegroupTagName(tags); ngNameTagValue == "" {
		return "", errors.New("failed to find the nodegroup name tag")
	}

	for _, tag := range tags {
		switch *tag.Key {
		case api.NodeGroupTypeTag:
			nodeGroupType = api.NodeGroupType(*tag.Value)
		}
	}

	if nodeGroupType == "" {
		nodeGroupType = api.NodeGroupTypeUnmanaged
	}

	return nodeGroupType, nil
}

// GetEksctlVersionFromTags returns the eksctl version used to create or update the stack
func GetEksctlVersionFromTags(tags []types.Tag) (semver.Version, bool, error) {
	for _, tag := range tags {
		if *tag.Key == api.EksctlVersionTag {
			v, err := version.ParseEksctlVersion(*tag.Value)
			if err != nil {
				return v, false, fmt.Errorf("unexpected error parsing eksctl version %q: %w", *tag.Value, err)
			}
			return v, true, nil
		}
	}
	return semver.Version{}, false, nil
}

// GetNodeGroupName will return nodegroup name based on tags
func (*StackCollection) GetNodeGroupName(s *Stack) string {
	if tagName := GetNodegroupTagName(s.Tags); tagName != "" {
		return tagName
	}
	if strings.HasSuffix(*s.StackName, "-nodegroup-0") {
		return "legacy-nodegroup-0"
	}
	if strings.HasSuffix(*s.StackName, "-DefaultNodeGroup") {
		return "legacy-default"
	}
	return ""
}

// GetNodegroupTagName returns the nodegroup name of a stack based on its tags. Taking into account legacy tags.
func GetNodegroupTagName(tags []types.Tag) string {
	for _, tag := range tags {
		switch *tag.Key {
		case api.NodeGroupNameTag, api.OldNodeGroupNameTag, api.OldNodeGroupIDTag:
			return *tag.Value
		}
	}
	return ""
}
