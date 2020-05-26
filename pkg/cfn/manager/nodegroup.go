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
	imageIDPath = resourcesRootPath + ".NodeGroupLaunchTemplate.Properties.LaunchTemplateData.ImageId"
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

// NodeGroupStack represents a nodegroup and its type
type NodeGroupStack struct {
	NodeGroupName string
	Type          api.NodeGroupType
}

// makeNodeGroupStackName generates the name of the nodegroup stack identified by its name, isolated by the cluster this StackCollection operates on
func (c *StackCollection) makeNodeGroupStackName(name string) string {
	return fmt.Sprintf("eksctl-%s-nodegroup-%s", c.spec.Metadata.Name, name)
}

// createNodeGroupTask creates the nodegroup
func (c *StackCollection) createNodeGroupTask(errs chan error, ng *api.NodeGroup, supportsManagedNodes bool) error {
	name := c.makeNodeGroupStackName(ng.Name)
	logger.Info("building nodegroup stack %q", name)
	stack := builder.NewNodeGroupResourceSet(c.provider, c.spec, c.makeClusterStackName(), ng, supportsManagedNodes)
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
	logger.Info("building managed nodegroup stack %q", name)
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

// ListNodeGroupStacks returns a list of NodeGroupStacks
func (c *StackCollection) ListNodeGroupStacks() ([]NodeGroupStack, error) {
	stacks, err := c.DescribeNodeGroupStacks()
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
		})
	}
	return nodeGroupStacks, nil
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
	descriptionBuffer.WriteString("scaling nodegroup")

	ngPaths, err := getNodeGroupPaths(stack.Tags)
	if err != nil {
		return err
	}
	var (
		desiredCapacityPath = ngPaths.DesiredCapacity
		maxSizePath         = ngPaths.MaxSize
		minSizePath         = ngPaths.MinSize
	)

	// TODO rewrite this using types
	// Get the current values
	currentCapacity := gjson.Get(template, desiredCapacityPath)
	currentMaxSize := gjson.Get(template, maxSizePath)
	currentMinSize := gjson.Get(template, minSizePath)

	hasChanged := func(desiredVal *int, currentVal gjson.Result) bool {
		return desiredVal != nil && int64(*desiredVal) != currentVal.Int()
	}
	changed := hasChanged(ng.DesiredCapacity, currentCapacity) || hasChanged(ng.MaxSize, currentMaxSize) || hasChanged(ng.MinSize, currentMinSize)

	if !changed {
		logger.Info("no change for nodegroup %q in cluster %q: nodes-min %d, desired %d, nodes-max %d", ng.Name,
			clusterName, currentMinSize.Int(), *ng.DesiredCapacity, currentMaxSize.Int())
		return nil
	}

	if ng.MinSize == nil && int64(*ng.DesiredCapacity) < currentMinSize.Int() {
		logger.Warning("the desired nodes %d is less than current nodes-min/minSize %d", *ng.DesiredCapacity, currentMinSize.Int())
		return errors.Errorf("the desired nodes %d is less than current nodes-min/minSize %d", *ng.DesiredCapacity, currentMinSize.Int())
	}

	if ng.MaxSize == nil && int64(*ng.DesiredCapacity) > currentMaxSize.Int() {
		logger.Warning("the desired nodes %d is greater than current nodes-max/maxSize %d", *ng.DesiredCapacity, currentMaxSize.Int())
		return errors.Errorf("the desired nodes %d is greater than current nodes-max/maxSize %d", *ng.DesiredCapacity, currentMaxSize.Int())
	}

	// Set the new values
	updateField := func(path, fieldName string, newVal *int, oldVal gjson.Result) error {
		if !hasChanged(newVal, oldVal) {
			return nil
		}
		template, err = sjson.Set(template, path, fmt.Sprintf("%d", *newVal))
		if err != nil {
			return errors.Wrapf(err, "error setting %s", fieldName)
		}
		descriptionBuffer.WriteString(fmt.Sprintf(", %s from %d to %d", fieldName, oldVal.Int(), *newVal))
		return nil
	}

	if err := updateField(desiredCapacityPath, "desired capacity", ng.DesiredCapacity, currentCapacity); err != nil {
		return err
	}

	if err := updateField(minSizePath, "min size", ng.MinSize, currentMinSize); err != nil {
		return err
	}

	if err := updateField(maxSizePath, "max size", ng.MaxSize, currentMaxSize); err != nil {
		return err
	}
	logger.Debug("stack template (post-scale change): %s", template)

	return c.UpdateStack(name, c.MakeChangeSetName("scale-nodegroup"), descriptionBuffer.String(), TemplateBody(template), nil)
}

// GetNodeGroupSummaries returns a list of summaries for the nodegroups of a cluster
func (c *StackCollection) GetNodeGroupSummaries(name string) ([]*NodeGroupSummary, error) {
	stacks, err := c.DescribeNodeGroupStacks()
	if err != nil {
		return nil, errors.Wrap(err, "getting nodegroup stacks")
	}

	var summaries []*NodeGroupSummary
	for _, s := range stacks {
		ngPaths, err := getNodeGroupPaths(s.Tags)
		if err != nil {
			return nil, err
		}

		summary, err := c.mapStackToNodeGroupSummary(s, ngPaths)
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

// GetNodeGroupStackType returns the nodegroup stack type
func (c *StackCollection) GetNodeGroupStackType(name string) (api.NodeGroupType, error) {
	stackName := c.makeNodeGroupStackName(name)
	stack, err := c.DescribeStack(&Stack{StackName: &stackName})
	if err != nil {
		return "", err
	}
	return GetNodeGroupType(stack.Tags)
}

// GetNodeGroupType returns the nodegroup type
func GetNodeGroupType(tags []*cfn.Tag) (api.NodeGroupType, error) {
	var (
		nodeGroupType api.NodeGroupType
	)
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

type nodeGroupPaths struct {
	InstanceType    string
	DesiredCapacity string
	MinSize         string
	MaxSize         string
}

func getNodeGroupPaths(tags []*cfn.Tag) (*nodeGroupPaths, error) {
	nodeGroupType, err := GetNodeGroupType(tags)
	if err != nil {
		return nil, err
	}

	switch nodeGroupType {
	case api.NodeGroupTypeManaged:
		makePath := func(fieldPath string) string {
			return fmt.Sprintf("%s.ManagedNodeGroup.Properties.%s", resourcesRootPath, fieldPath)
		}
		makeScalingPath := func(field string) string {
			return makePath(fmt.Sprintf("ScalingConfig.%s", field))

		}
		return &nodeGroupPaths{
			InstanceType:    makePath("InstanceTypes.0"),
			DesiredCapacity: makeScalingPath("DesiredSize"),
			MinSize:         makeScalingPath("MinSize"),
			MaxSize:         makeScalingPath("MaxSize"),
		}, nil

		// Tag may not exist for existing nodegroups
	case api.NodeGroupTypeUnmanaged, "":
		makePath := func(field string) string {
			return fmt.Sprintf("%s.NodeGroup.Properties.%s", resourcesRootPath, field)
		}
		return &nodeGroupPaths{
			InstanceType:    resourcesRootPath + ".NodeGroupLaunchTemplate.Properties.LaunchTemplateData.InstanceType",
			DesiredCapacity: makePath("DesiredCapacity"),
			MinSize:         makePath("MinSize"),
			MaxSize:         makePath("MaxSize"),
		}, nil

	default:
		return nil, fmt.Errorf("unexpected nodegroup type tag: %q", nodeGroupType)
	}

}

func (c *StackCollection) mapStackToNodeGroupSummary(stack *Stack, ngPaths *nodeGroupPaths) (*NodeGroupSummary, error) {
	template, err := c.GetStackTemplate(*stack.StackName)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting CloudFormation template for stack %s", *stack.StackName)
	}

	cluster := getClusterNameTag(stack)
	name := c.GetNodeGroupName(stack)
	maxSize := gjson.Get(template, ngPaths.MaxSize)
	minSize := gjson.Get(template, ngPaths.MinSize)
	desired := gjson.Get(template, ngPaths.DesiredCapacity)
	instanceType := gjson.Get(template, ngPaths.InstanceType)
	imageID := gjson.Get(template, imageIDPath)

	nodeGroupType, err := GetNodeGroupType(stack.Tags)
	if err != nil {
		return nil, err
	}

	var nodeInstanceRoleARN string
	if nodeGroupType == api.NodeGroupTypeUnmanaged {
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
func GetNodegroupTagName(tags []*cfn.Tag) string {
	for _, tag := range tags {
		switch *tag.Key {
		case api.NodeGroupNameTag, api.OldNodeGroupNameTag, api.OldNodeGroupIDTag:
			return *tag.Value
		}
	}
	return ""
}
