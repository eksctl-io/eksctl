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

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
)

const (
	desiredCapacityPath = resourcesRootPath + ".NodeGroup.Properties.DesiredCapacity"
	maxSizePath         = resourcesRootPath + ".NodeGroup.Properties.MaxSize"
	minSizePath         = resourcesRootPath + ".NodeGroup.Properties.MinSize"
	instanceTypePath    = resourcesRootPath + ".NodeLaunchConfig.Properties.InstanceType"
	imageIDPath         = resourcesRootPath + ".NodeLaunchConfig.Properties.ImageId"
)

// NodeGroupSummary represents a summary of a nodegroup stack
type NodeGroupSummary struct {
	StackName       string
	Cluster         string
	Name            string
	MaxSize         int
	MinSize         int
	DesiredCapacity int
	InstanceType    string
	ImageID         string
	CreationTime    *time.Time
}

// MakeNodeGroupStackName generates the name of the node group identified by its ID, isolated by the cluster this StackCollection operates on
func (c *StackCollection) MakeNodeGroupStackName(name string) string {
	return fmt.Sprintf("eksctl-%s-nodegroup-%s", c.spec.Metadata.Name, name)
}

// CreateNodeGroup creates the nodegroup
func (c *StackCollection) CreateNodeGroup(errs chan error, data interface{}) error {
	ng := data.(*api.NodeGroup)
	name := c.MakeNodeGroupStackName(ng.Name)
	logger.Info("creating nodegroup stack %q", name)
	stack := builder.NewNodeGroupResourceSet(c.provider, c.spec, c.makeClusterStackName(), ng)
	if err := stack.AddAllResources(); err != nil {
		return err
	}

	if ng.Tags == nil {
		ng.Tags = make(map[string]string)
	}
	ng.Tags[api.NodeGroupNameTag] = ng.Name

	return c.CreateStack(name, stack, ng.Tags, nil, errs)
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
		if getNodeGroupName(s) != "" {
			nodeGroupStacks = append(nodeGroupStacks, s)
		}
	}
	logger.Debug("nodegroups = %v", nodeGroupStacks)
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
		allResources[getNodeGroupName(s)] = StackInfo{
			Resources: resources.StackResources,
			Template:  &template,
			Stack:     s,
		}
	}

	return allResources, nil
}

// DeleteNodeGroup deletes a nodegroup stack
func (c *StackCollection) DeleteNodeGroup(name string) error {
	name = c.MakeNodeGroupStackName(name)
	_, err := c.DeleteStack(name, false)
	return err
}

// WaitDeleteNodeGroup waits until the nodegroup is deleted,
// it calls WaitDeleteStack without force
func (c *StackCollection) WaitDeleteNodeGroup(errs chan error, data interface{}) error {
	name := c.MakeNodeGroupStackName(data.(string))
	return c.WaitDeleteStack(name, false, errs)
}

// WaitForceDeleteNodeGroup waits until the nodegroup is deleted,
// it calls WaitDeleteStack with force
func (c *StackCollection) WaitForceDeleteNodeGroup(errs chan error, data interface{}) error {
	name := c.MakeNodeGroupStackName(data.(string))
	return c.WaitDeleteStack(name, true, errs)
}

// BlockingWaitDeleteNodeGroup waits until the nodegroup is deleted
func (c *StackCollection) BlockingWaitDeleteNodeGroup(name string, force bool) error {
	name = c.MakeNodeGroupStackName(name)
	return c.BlockingWaitDeleteStack(name, force)
}

// ScaleNodeGroup will scale an existing nodegroup
func (c *StackCollection) ScaleNodeGroup(ng *api.NodeGroup) error {
	clusterName := c.makeClusterStackName()
	c.spec.Status = &api.ClusterStatus{StackName: clusterName}
	name := c.MakeNodeGroupStackName(ng.Name)
	logger.Info("scaling nodegroup stack %q in cluster %s", name, clusterName)

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

	return c.UpdateStack(name, "scale-nodegroup", descriptionBuffer.String(), []byte(template), nil)
}

// GetNodeGroupSummaries returns a list of summaries for the nodegroups of a cluster
func (c *StackCollection) GetNodeGroupSummaries(name string) ([]*NodeGroupSummary, error) {
	stacks, err := c.DescribeNodeGroupStacks()
	if err != nil {
		return nil, errors.Wrap(err, "getting nodegroup stacks")
	}

	summaries := []*NodeGroupSummary{}
	for _, s := range stacks {
		summary, err := c.mapStackToNodeGroupSummary(s)
		if err != nil {
			return nil, errors.Wrap(err, "mapping stack to nodegorup summary")
		}

		if name == "" {
			summaries = append(summaries, summary)
		} else if summary.Name == name {
			summaries = append(summaries, summary)
		}
	}

	return summaries, nil
}

func (c *StackCollection) mapStackToNodeGroupSummary(stack *Stack) (*NodeGroupSummary, error) {
	template, err := c.GetStackTemplate(*stack.StackName)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting Cloudformation template for stack %s", *stack.StackName)
	}

	cluster := getClusterNameTag(stack)
	name := getNodeGroupName(stack)
	maxSize := gjson.Get(template, maxSizePath)
	minSize := gjson.Get(template, minSizePath)
	desired := gjson.Get(template, desiredCapacityPath)
	instanceType := gjson.Get(template, instanceTypePath)
	imageID := gjson.Get(template, imageIDPath)

	summary := &NodeGroupSummary{
		StackName:       *stack.StackName,
		Cluster:         cluster,
		Name:            name,
		MaxSize:         int(maxSize.Int()),
		MinSize:         int(minSize.Int()),
		DesiredCapacity: int(desired.Int()),
		InstanceType:    instanceType.String(),
		ImageID:         imageID.String(),
		CreationTime:    stack.CreationTime,
	}

	return summary, nil
}

func getNodeGroupName(s *Stack) string {
	for _, tag := range s.Tags {
		if *tag.Key == api.NodeGroupNameTag {
			return *tag.Value
		}
		if *tag.Key == api.OldNodeGroupIDTag {
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
