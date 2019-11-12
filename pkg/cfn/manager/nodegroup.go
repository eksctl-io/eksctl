package manager

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
)

const (
	instanceTypePath = resourcesRootPath + ".NodeGroupLaunchTemplate.Properties.LaunchTemplateData.InstanceType"
	imageIDPath      = resourcesRootPath + ".NodeGroupLaunchTemplate.Properties.LaunchTemplateData.ImageId"
)

// NodeGroupSummary represents a summary of a nodegroup stack
type NodeGroupSummary struct {
	StackName           string
	Cluster             string
	Name                string
	MaxSize             int
	MinSize             int
	DesiredCapacity     int
	InstanceType        string
	ImageID             string
	CreationTime        *time.Time
	NodeInstanceRoleARN string
}

// makeNodeGroupStackName generates the name of the nodegroup stack identified by its name, isolated by the cluster this StackCollection operates on
func (c *StackCollection) makeNodeGroupStackName(name string) string {
	return fmt.Sprintf("eksctl-%s-nodegroup-%s", c.spec.Metadata.Name, name)
}

// createNodeGroupTask creates the nodegroup
func (c *StackCollection) createNodeGroupTask(errs chan error, ng *api.NodeGroup) error {
	name := c.makeNodeGroupStackName(ng.Name)
	logger.Info("building nodegroup stack %q", name)
	stack := builder.NewNodeGroupResourceSet(c.provider, c.spec, c.makeClusterStackName(), ng)
	if err := stack.AddAllResources(); err != nil {
		return err
	}

	if ng.Tags == nil {
		ng.Tags = make(map[string]string)
	}
	ng.Tags[api.NodeGroupNameTag] = ng.Name
	ng.Tags[api.OldNodeGroupNameTag] = ng.Name
	ng.Tags[api.NodeGroupTypeTag] = string(api.NodeGroupTypeUnmanaged)

	return c.CreateStack(name, stack, ng.Tags, nil, errs)
}

func (c *StackCollection) createManagedNodeGroupTask(errorCh chan error, ng *api.ManagedNodeGroup) error {
	name := c.makeNodeGroupStackName(ng.Name)
	stack := builder.NewManagedNodeGroup(c.spec, ng, c.makeClusterStackName())
	if err := stack.AddAllResources(); err != nil {
		return err
	}
	return c.CreateStack(name, stack, ng.Tags, nil, errorCh)
}

// DescribeNodeGroupStacks calls DescribeStacks and filters out nodegroups
func (c *StackCollection) DescribeNodeGroupStacks() ([]*Stack, error) {
	stacks, err := c.DescribeStacks()
	if err != nil {
		return nil, err
	}

	nodeGroupStacks := []*Stack{}
	for _, s := range stacks {
		if *s.StackStatus == cfn.StackStatusDeleteComplete {
			continue
		}
		if c.GetNodeGroupName(s) != "" {
			nodeGroupStacks = append(nodeGroupStacks, s)
		}
	}
	logger.Debug("nodegroups = %v", nodeGroupStacks)
	return nodeGroupStacks, nil
}

// ListNodeGroupStacks calls DescribeNodeGroupStacks and returns only nodegroup names
func (c *StackCollection) ListNodeGroupStacks() ([]string, error) {
	stacks, err := c.DescribeNodeGroupStacks()
	if err != nil {
		return nil, err
	}

	names := []string{}
	for _, s := range stacks {
		names = append(names, c.GetNodeGroupName(s))
	}
	return names, nil
}

// DescribeNodeGroupStacksAndResources calls DescribeNodeGroupStacks and fetches all resources,
// then returns it in a map by nodegroup name
func (c *StackCollection) DescribeNodeGroupStacksAndResources() (map[string]StackInfo, error) {
	stacks, err := c.DescribeNodeGroupStacks()
	if err != nil {
		return nil, err
	}

	allResources := make(map[string]StackInfo)

	for _, s := range stacks {
		input := &cfn.DescribeStackResourcesInput{
			StackName: s.StackName,
		}
		template, err := c.GetStackTemplate(*s.StackName)
		if err != nil {
			return nil, errors.Wrapf(err, "getting template for %q stack", *s.StackName)
		}
		resources, err := c.provider.CloudFormation().DescribeStackResources(input)
		if err != nil {
			return nil, errors.Wrapf(err, "getting all resources for %q stack", *s.StackName)
		}
		allResources[c.GetNodeGroupName(s)] = StackInfo{
			Resources: resources.StackResources,
			Template:  &template,
			Stack:     s,
		}
	}

	return allResources, nil
}

// ScaleNodeGroup will scale an existing nodegroup
func (c *StackCollection) ScaleNodeGroup(ng *api.NodeGroup) error {
	clusterName := c.makeClusterStackName()
	c.spec.Status = &api.ClusterStatus{StackName: clusterName}
	name := c.makeNodeGroupStackName(ng.Name)
	logger.Info("scaling nodegroup stack %q in cluster %s", name, clusterName)

	stack, err := c.DescribeStack(&Stack{StackName: &name})
	if err != nil {
		return errors.Wrapf(err, "error describing nodegroup stack %s", name)
	}

	// Get current stack
	template, err := c.GetStackTemplate(name)
	if err != nil {
		return errors.Wrapf(err, "error getting stack template %s", name)
	}
	logger.Debug("stack template (pre-scale change): %s", template)

	//TODO: In the future we might want to use Goformation for strongly typed
	//manipulation of the template.

	var descriptionBuffer bytes.Buffer
	descriptionBuffer.WriteString("scaling nodegroup, ")

	scalingPaths, err := getScalingConfigPaths(stack.Tags)
	if err != nil {
		return err
	}
	var (
		desiredCapacityPath = scalingPaths.DesiredCapacity
		maxSizePath         = scalingPaths.MaxSize
		minSizePath         = scalingPaths.MinSize
	)

	// TODO rewrite this using types
	// Get the current values
	currentCapacity := gjson.Get(template, desiredCapacityPath)
	currentMaxSize := gjson.Get(template, maxSizePath)
	currentMinSize := gjson.Get(template, minSizePath)

	if ng.DesiredCapacity != nil && int64(*ng.DesiredCapacity) == currentCapacity.Int() {
		logger.Info("desired capacity of nodegroup %q in cluster %q is already %d", ng.Name, clusterName, *ng.DesiredCapacity)
		return nil
	}

	// Set the new values
	newCapacity := fmt.Sprintf("%d", *ng.DesiredCapacity)
	template, err = sjson.Set(template, desiredCapacityPath, newCapacity)
	if err != nil {
		return errors.Wrap(err, "setting desired capacity")
	}
	descriptionBuffer.WriteString(fmt.Sprintf("desired capacity from %s to %d", currentCapacity.Str, *ng.DesiredCapacity))

	// If the desired number of nodes is less than the min then update the min
	if int64(*ng.DesiredCapacity) < currentMinSize.Int() {
		newMinSize := fmt.Sprintf("%d", *ng.DesiredCapacity)
		template, err = sjson.Set(template, minSizePath, newMinSize)
		if err != nil {
			return errors.Wrap(err, "setting min size")
		}
		descriptionBuffer.WriteString(fmt.Sprintf(", min size from %s to %d", currentMinSize.Str, *ng.DesiredCapacity))
	}
	// If the desired number of nodes is greater than the max then update the max
	if int64(*ng.DesiredCapacity) > currentMaxSize.Int() {
		newMaxSize := fmt.Sprintf("%d", *ng.DesiredCapacity)
		template, err = sjson.Set(template, maxSizePath, newMaxSize)
		if err != nil {
			return errors.Wrap(err, "setting max size")
		}
		descriptionBuffer.WriteString(fmt.Sprintf(", max size from %s to %d", currentMaxSize.Str, *ng.DesiredCapacity))
	}
	logger.Debug("stack template (post-scale change): %s", template)

	return c.UpdateStack(name, c.MakeChangeSetName("scale-nodegroup"), descriptionBuffer.String(), []byte(template), nil)
}

// GetNodeGroupSummaries returns a list of summaries for the nodegroups of a cluster
func (c *StackCollection) GetNodeGroupSummaries(name string) ([]*NodeGroupSummary, error) {
	stacks, err := c.DescribeNodeGroupStacks()
	if err != nil {
		return nil, errors.Wrap(err, "getting nodegroup stacks")
	}

	summaries := []*NodeGroupSummary{}
	for _, s := range stacks {
		scalingPaths, err := getScalingConfigPaths(s.Tags)
		if err != nil {
			return nil, err
		}

		summary, err := c.mapStackToNodeGroupSummary(s, scalingPaths)
		if err != nil {
			return nil, errors.Wrap(err, "mapping stack to nodegroup summary")
		}

		if name == "" {
			summaries = append(summaries, summary)
		} else if summary.Name == name {
			summaries = append(summaries, summary)
		}
	}

	return summaries, nil
}

type scalingConfigPaths struct {
	DesiredCapacity string
	MinSize         string
	MaxSize         string
}

func getScalingConfigPaths(tags []*cfn.Tag) (*scalingConfigPaths, error) {
	var (
		foundNodeGroupTag bool
		nodeGroupType     string
	)
	for _, tag := range tags {
		switch *tag.Key {
		case api.NodeGroupNameTag:
			foundNodeGroupTag = true
		case api.NodeGroupTypeTag:
			nodeGroupType = *tag.Value
		}
	}

	if !foundNodeGroupTag {
		return nil, fmt.Errorf("failed to find a nodegroup tag (%s)", api.NodeGroupNameTag)
	}

	switch api.NodeGroupType(nodeGroupType) {
	case api.NodeGroupTypeManaged:
		makePath := func(field string) string {
			return fmt.Sprintf("%s.ManagedNodeGroup.Properties.ScalingConfig.%s", resourcesRootPath, field)
		}
		return &scalingConfigPaths{
			DesiredCapacity: makePath("DesiredSize"),
			MinSize:         makePath("MinSize"),
			MaxSize:         makePath("MaxSize"),
		}, nil

		// Tag may not exist for existing nodegroups
	case api.NodeGroupTypeUnmanaged, "":
		makePath := func(field string) string {
			return fmt.Sprintf("%s.NodeGroup.Properties.%s", resourcesRootPath, field)
		}
		return &scalingConfigPaths{
			DesiredCapacity: makePath("DesiredCapacity"),
			MinSize:         makePath("MaxSize"),
			MaxSize:         makePath("MinSize"),
		}, nil

	default:
		return nil, fmt.Errorf("unexpected nodegroup type tag: %q", nodeGroupType)
	}

}

func (c *StackCollection) mapStackToNodeGroupSummary(stack *Stack, scalingPaths *scalingConfigPaths) (*NodeGroupSummary, error) {
	template, err := c.GetStackTemplate(*stack.StackName)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting CloudFormation template for stack %s", *stack.StackName)
	}

	cluster := getClusterNameTag(stack)
	name := c.GetNodeGroupName(stack)
	maxSize := gjson.Get(template, scalingPaths.MaxSize)
	minSize := gjson.Get(template, scalingPaths.MinSize)
	desired := gjson.Get(template, scalingPaths.DesiredCapacity)
	// FIXME for managed nodes
	instanceType := gjson.Get(template, instanceTypePath)
	imageID := gjson.Get(template, imageIDPath)

	var nodeInstanceRoleARN string
	nodeInstanceRoleARNCollector := func(s string) error {
		nodeInstanceRoleARN = s
		return nil
	}
	collectors := map[string]outputs.Collector{
		outputs.NodeGroupInstanceRoleARN: nodeInstanceRoleARNCollector,
	}
	collectorSet := outputs.NewCollectorSet(collectors)
	if err := collectorSet.MustCollect(*stack); err != nil {
		return nil, errors.Wrapf(err, "error collecting Cloudformation outputs for stack %s", *stack.StackName)
	}

	summary := &NodeGroupSummary{
		StackName:           *stack.StackName,
		Cluster:             cluster,
		Name:                name,
		MaxSize:             int(maxSize.Int()),
		MinSize:             int(minSize.Int()),
		DesiredCapacity:     int(desired.Int()),
		InstanceType:        instanceType.String(),
		ImageID:             imageID.String(),
		CreationTime:        stack.CreationTime,
		NodeInstanceRoleARN: nodeInstanceRoleARN,
	}

	return summary, nil
}

// GetNodeGroupName will return nodegroup name based on tags
func (*StackCollection) GetNodeGroupName(s *Stack) string {
	for _, tag := range s.Tags {
		switch *tag.Key {
		case api.NodeGroupNameTag, api.OldNodeGroupNameTag, api.OldNodeGroupIDTag:
			return *tag.Value
		}
	}
	if strings.HasSuffix(*s.StackName, "-nodegroup-0") {
		return "legacy-nodegroup-0"
	}
	if strings.HasSuffix(*s.StackName, "-DefaultNodeGroup") {
		return "legacy-default"
	}
	return ""
}
