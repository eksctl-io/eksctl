package manager

import (
	"bytes"
	"fmt"
	"time"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/eks/api"
)

const (
	desiredCapacityPath = "Resources.NodeGroup.Properties.DesiredCapacity"
	maxSizePath         = "Resources.NodeGroup.Properties.MaxSize"
	minSizePath         = "Resources.NodeGroup.Properties.MinSize"
	instanceTypePath    = "Resources.NodeLaunchConfig.Properties.InstanceType"
	imageIDPath         = "Resources.NodeLaunchConfig.Properties.ImageId"
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
	stack := builder.NewNodeGroupResourceSet(c.spec, c.makeClusterStackName(), ng)
	if err := stack.AddAllResources(); err != nil {
		return err
	}

	c.tags = append(c.tags, newTag(NodeGroupNameTag, fmt.Sprintf("%s", ng.Name)))

	for k, v := range ng.Tags {
		c.tags = append(c.tags, newTag(k, v))
	}

	return c.CreateStack(name, stack, nil, errs)
}

func (c *StackCollection) listAllNodeGroups() ([]string, error) {
	stacks, err := c.ListStacks(fmt.Sprintf("^eksctl-%s-nodegroup-.+$", c.spec.Metadata.Name))
	if err != nil {
		return nil, err
	}
	stackNames := []string{}
	for _, s := range stacks {
		if *s.StackStatus == cfn.StackStatusDeleteComplete {
			continue
		}
		stackNames = append(stackNames, getNodeGroupName(s.Tags))
	}
	logger.Debug("nodegroups = %v", stackNames)
	return stackNames, nil
}

// DeleteNodeGroup deletes a nodegroup stack
func (c *StackCollection) DeleteNodeGroup(errs chan error, data interface{}) error {
	defer close(errs)
	name := data.(string)
	stack := c.MakeNodeGroupStackName(name)
	_, err := c.DeleteStack(stack)
	errs <- err
	return nil
}

// WaitDeleteNodeGroup waits until the nodegroup is deleted
func (c *StackCollection) WaitDeleteNodeGroup(errs chan error, data interface{}) error {
	name := data.(string)
	stack := c.MakeNodeGroupStackName(name)
	return c.WaitDeleteStackTask(stack, errs)
}

// ScaleInitialNodeGroup will scale the first nodegroup (ID: 0)
func (c *StackCollection) ScaleInitialNodeGroup() error {
	return c.ScaleNodeGroup(c.spec.NodeGroups[0])
}

// ScaleNodeGroup will scale an existing nodegroup
func (c *StackCollection) ScaleNodeGroup(ng *api.NodeGroup) error {
	clusterName := c.makeClusterStackName()
	c.spec.ClusterStackName = clusterName
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

	// Set the new values
	newCapacity := fmt.Sprintf("%d", ng.DesiredCapacity)
	template, err = sjson.Set(template, desiredCapacityPath, newCapacity)
	if err != nil {
		return errors.Wrap(err, "setting desired capacity")
	}
	descriptionBuffer.WriteString(fmt.Sprintf("desired capacity from %s to %d", currentCapacity.Str, ng.DesiredCapacity))

	// If the desired number of nodes is less than the min then update the min
	if int64(ng.DesiredCapacity) < currentMinSize.Int() {
		newMinSize := fmt.Sprintf("%d", ng.DesiredCapacity)
		template, err = sjson.Set(template, minSizePath, newMinSize)
		if err != nil {
			return errors.Wrap(err, "setting min size")
		}
		descriptionBuffer.WriteString(fmt.Sprintf(", min size from %s to %d", currentMinSize.Str, ng.DesiredCapacity))
	}
	// If the desired number of nodes is greater than the max then update the max
	if int64(ng.DesiredCapacity) > currentMaxSize.Int() {
		newMaxSize := fmt.Sprintf("%d", ng.DesiredCapacity)
		template, err = sjson.Set(template, maxSizePath, newMaxSize)
		if err != nil {
			return errors.Wrap(err, "setting max size")
		}
		descriptionBuffer.WriteString(fmt.Sprintf(", max size from %s to %d", currentMaxSize.Str, ng.DesiredCapacity))
	}
	logger.Debug("stack template (post-scale change): %s", template)

	return c.UpdateStack(name, "scale-nodegroup", descriptionBuffer.String(), []byte(template), nil)
}

// GetNodeGroupSummaries returns a list of summaries for the nodegroups of a cluster
func (c *StackCollection) GetNodeGroupSummaries() ([]*NodeGroupSummary, error) {
	stacks, err := c.ListStacks(fmt.Sprintf("^(eksctl|EKS)-%s-nodegroup-.+$", c.spec.Metadata.Name))
	if err != nil {
		return nil, errors.Wrap(err, "getting nodegroup stacks")
	}

	summaries := []*NodeGroupSummary{}
	for _, stack := range stacks {
		logger.Info("stack %s\n", *stack.StackName)
		logger.Debug("stack = %#v", stack)

		summary, err := c.mapStackToNodeGroupSummary(stack)
		if err != nil {
			return nil, errors.New("error mapping stack to node gorup summary")
		}

		summaries = append(summaries, summary)
	}

	return summaries, nil
}

func (c *StackCollection) mapStackToNodeGroupSummary(stack *Stack) (*NodeGroupSummary, error) {
	template, err := c.GetStackTemplate(*stack.StackName)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting Cloudformation template for stack %s", *stack.StackName)
	}

	cluster := getClusterName(stack.Tags)
	name := getNodeGroupName(stack.Tags)
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

func getNodeGroupName(tags []*cfn.Tag) string {
	for _, tag := range tags {
		if *tag.Key == NodeGroupNameTag {
			return *tag.Value
		}
	}
	return ""
}

func getClusterName(tags []*cfn.Tag) string {
	for _, tag := range tags {
		if *tag.Key == ClusterNameTag {
			return *tag.Value
		}
	}
	return ""
}
