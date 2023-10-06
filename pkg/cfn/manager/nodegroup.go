package manager

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	asgtypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	cfn "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/eks"

	"github.com/blang/semver"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
	"github.com/weaveworks/eksctl/pkg/version"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

// NodeGroupStack represents a nodegroup and its type
type NodeGroupStack struct {
	NodeGroupName string
	Type          api.NodeGroupType
	Stack         *Stack
}

// makeNodeGroupStackName generates the name of the nodegroup stack identified by its name, isolated by the cluster this StackCollection operates on
func (c *StackCollection) makeNodeGroupStackName(name string) string {
	return fmt.Sprintf("eksctl-%s-nodegroup-%s", c.spec.Metadata.Name, name)
}

// createNodeGroupTask creates the nodegroup
func (c *StackCollection) createNodeGroupTask(ctx context.Context, errs chan error, ng *api.NodeGroup, forceAddCNIPolicy, skipEgressRules bool, vpcImporter vpc.Importer) error {
	name := c.makeNodeGroupStackName(ng.Name)

	logger.Info("building nodegroup stack %q", name)
	bootstrapper, err := nodebootstrap.NewBootstrapper(c.spec, ng)
	if err != nil {
		return errors.Wrap(err, "error creating bootstrapper")
	}
	stack := builder.NewNodeGroupResourceSet(c.ec2API, c.iamAPI, builder.NodeGroupOptions{
		ClusterConfig:     c.spec,
		NodeGroup:         ng,
		Bootstrapper:      bootstrapper,
		ForceAddCNIPolicy: forceAddCNIPolicy,
		VPCImporter:       vpcImporter,
		SkipEgressRules:   skipEgressRules,
	})
	if err := stack.AddAllResources(ctx); err != nil {
		return err
	}

	if ng.Tags == nil {
		ng.Tags = make(map[string]string)
	}
	ng.Tags[api.NodeGroupNameTag] = ng.Name
	ng.Tags[api.OldNodeGroupNameTag] = ng.Name
	ng.Tags[api.NodeGroupTypeTag] = string(api.NodeGroupTypeUnmanaged)

	return c.CreateStack(ctx, name, stack, ng.Tags, nil, errs)
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
		return errors.Wrapf(err, "couldn't get managed nodegroup details for nodegroup %q", ng.Name)
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

// ListNodeGroupStacks returns a list of NodeGroupStacks
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
			NodeGroupName: c.GetNodeGroupName(stack),
			Type:          nodeGroupType,
			Stack:         stack,
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
			return nil, errors.Wrapf(err, "getting all resources for %q stack", *s.StackName)
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
				return v, false, errors.Wrapf(err, "unexpected error parsing eksctl version %q", *tag.Value)
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
