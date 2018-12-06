package manager

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/eks/api"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
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
	Seq             int
	StackName       string
	MaxSize         int
	MinSize         int
	DesiredCapacity int
	InstanceType    string
	ImageID         string
	CreationTime    *time.Time
}

// MakeNodeGroupStackName generates the name of the node group identified by its ID, isolated by the cluster this StackCollection operates on
func (c *StackCollection) MakeNodeGroupStackName(id int) string {
	return fmt.Sprintf("eksctl-%s-nodegroup-%d", c.spec.Metadata.Name, id)
}

// CreateEmbeddedNodeGroup creates the nodegroup embedded in the cluster spec
func (c *StackCollection) CreateEmbeddedNodeGroup(errs chan error, data interface{}) error {
	ng := data.(*api.NodeGroup)
	name := c.MakeNodeGroupStackName(ng.ID)
	logger.Info("creating nodegroup stack %q", name)
	stack := builder.NewEmbeddedNodeGroupResourceSet(c.spec, c.makeClusterStackName(), ng.ID)
	if err := stack.AddAllResources(); err != nil {
		return err
	}

	c.tags = append(c.tags, newTag(NodeGroupIDTag, fmt.Sprintf("%d", ng.ID)))

	for k, v := range ng.Tags {
		c.tags = append(c.tags, newTag(k, v))
	}

	return c.CreateStack(name, stack, nil, errs)
}

// CreateNodeGroup creates the nodegroup
func (c *StackCollection) CreateNodeGroup(errs chan error, data interface{}) error {
	ng := data.(*api.NodeGroup)
	name := c.MakeNodeGroupStackName(ng.ID)
	logger.Info("creating nodegroup stack %q", name)
	stack := builder.NewNodeGroupResourceSet(c.spec, c.makeClusterStackName(), ng)
	if err := stack.AddAllResources(); err != nil {
		return err
	}

	c.tags = append(c.tags, newTag(NodeGroupIDTag, fmt.Sprintf("%d", ng.ID)))

	for k, v := range ng.Tags {
		c.tags = append(c.tags, newTag(k, v))
	}

	return c.CreateStack(name, stack, nil, errs)
}

func (c *StackCollection) listAllNodeGroupStacks() ([]string, error) {
	stacks, err := c.ListStacks(fmt.Sprintf("^eksctl-%s-nodegroup-\\d$", c.spec.Metadata.Name))
	if err != nil {
		return nil, err
	}
	stackNames := []string{}
	for _, s := range stacks {
		if *s.StackStatus == cfn.StackStatusDeleteComplete {
			continue
		}
		stackNames = append(stackNames, *s.StackName)
	}
	logger.Debug("nodegroups = %v", stackNames)
	return stackNames, nil
}

// DeleteNodeGroup deletes a nodegroup stack
func (c *StackCollection) DeleteNodeGroup(errs chan error, data interface{}) error {
	name := data.(string)
	_, err := c.DeleteStack(name)
	return err
}

// WaitDeleteNodeGroup waits until the nodegroup is deleted
func (c *StackCollection) WaitDeleteNodeGroup(errs chan error, data interface{}) error {
	name := data.(string)
	return c.WaitDeleteStack(name)
}

// ScaleInitialNodeGroup will scale the first nodegroup (ID: 0)
func (c *StackCollection) ScaleInitialNodeGroup() error {
	return c.ScaleNodeGroup(c.spec.NodeGroups[0])
}

// ScaleNodeGroup will scale an existing nodegroup
func (c *StackCollection) ScaleNodeGroup(ng *api.NodeGroup) error {
	clusterName := c.makeClusterStackName()
	c.spec.ClusterStackName = clusterName
	name := c.MakeNodeGroupStackName(ng.ID)
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
	stacks, err := c.ListStacks(fmt.Sprintf("^(eksctl|EKS)-%s-nodegroup-\\d+$", c.spec.Metadata.Name), cfn.StackStatusCreateComplete)
	if err != nil {
		return nil, errors.Wrap(err, "getting nodegroup stacks")
	}

	summaries := []*NodeGroupSummary{}
	for _, stack := range stacks {
		logger.Info("stack %s\n", *stack.StackName)
		logger.Debug("stack = %#v", stack)

		err, summary := c.mapStackToNodeGroupSummary(stack)
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

	seq := getNodeGroupID(stack.Tags)
	maxSize := gjson.Get(template, maxSizePath)
	minSize := gjson.Get(template, minSizePath)
	desired := gjson.Get(template, desiredCapacityPath)
	instanceType := gjson.Get(template, instanceTypePath)
	imageID := gjson.Get(template, imageIDPath)

	summary := &NodeGroupSummary{
		Seq:             seq,
		StackName:       *stack.StackName,
		MaxSize:         int(maxSize.Int()),
		MinSize:         int(minSize.Int()),
		DesiredCapacity: int(desired.Int()),
		InstanceType:    instanceType.String(),
		ImageID:         imageID.String(),
		CreationTime:    stack.CreationTime,
	}

	return summary, nil
}

// GetMaxNodeGroupSeq returns the sequence number og the highest node group
func (c *StackCollection) GetMaxNodeGroupSeq() (int, error) {
	stacks, err := c.ListStacks(fmt.Sprintf("^(eksctl|EKS)-%s-nodegroup-\\d+$", c.spec.Metadata.Name))
	if err != nil {
		return -1, errors.Wrap(err, "getting nodegroup stacks")
	}
	seq := -1
	for _, stack := range stacks {
		stackSeq := getNodeGroupID(stack.Tags)
		if stackSeq > seq {
			seq = stackSeq
		}
	}
	logger.Debug("stacks = %v", stacks)
	return seq, nil
}

func getNodeGroupID(tags []*cfn.Tag) int {
	for _, tag := range tags {
		if *tag.Key == NodeGroupIDTag {
			i, _ := strconv.Atoi(*tag.Value)
			return i
		}
	}
	return -1
}
