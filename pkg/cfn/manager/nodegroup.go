package manager

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const (
	desirecCapacityPath = "Resources.NodeGroup.Properties.DesiredCapacity"
	maxSizePath         = "Resources.NodeGroup.Properties.MaxSize"
	minSizePath         = "Resources.NodeGroup.Properties.MinSize"
)

func (c *StackCollection) makeNodeGroupStackName(sequence int) string {
	return fmt.Sprintf("eksctl-%s-nodegroup-%d", c.spec.ClusterName, sequence)
}

// CreateInitialNodeGroup creates the initial node group
func (c *StackCollection) CreateInitialNodeGroup(errs chan error) error {
	return c.CreateNodeGroup(0, errs)
}

// CreateNodeGroup creates the node group
func (c *StackCollection) CreateNodeGroup(seq int, errs chan error) error {
	name := c.makeNodeGroupStackName(seq)
	logger.Info("creating nodegroup stack %q", name)
	stack := builder.NewNodeGroupResourceSet(c.spec, c.makeClusterStackName(), seq)
	if err := stack.AddAllResources(); err != nil {
		return err
	}

	c.tags = append(c.tags, newTag(NodeGroupIDTag, fmt.Sprintf("%d", seq)))

	return c.CreateStack(name, stack, nil, errs)
}

// DeleteNodeGroup deletes the node group
func (c *StackCollection) DeleteNodeGroup() error {
	_, err := c.DeleteStack(c.makeNodeGroupStackName(0))
	return err
}

// WaitDeleteNodeGroup waits till the node group is deleted
func (c *StackCollection) WaitDeleteNodeGroup() error {
	return c.WaitDeleteStack(c.makeNodeGroupStackName(0))
}

// ScaleInitialNodeGroup will scale the first (sequence 0) nodegroup
func (c *StackCollection) ScaleInitialNodeGroup() error {
	return c.ScaleNodeGroup(0)
}

// ScaleNodeGroup will scale an existing node group
func (c *StackCollection) ScaleNodeGroup(sequence int) error {
	clusterName := c.makeClusterStackName()
	c.spec.ClusterStackName = clusterName
	name := c.makeNodeGroupStackName(sequence)
	logger.Info("scaling nodegroup stack %q in cluster %s", name, clusterName)

	if c.spec.MinNodes == 0 && c.spec.MaxNodes == 0 {
		c.spec.MinNodes = c.spec.Nodes
		c.spec.MaxNodes = c.spec.Nodes
	}

	// Get current stack
	template, err := c.getStackTemplate(name)
	if err != nil {
		return errors.Wrapf(err, "error getting stack template %s", name)
	}
	logger.Debug("stack template (pre-scale change): %s", template)

	// Get the current values and set new values
	//TODO: In the future we might want to use Goformation for stronly typed
	//manipulation of the temaple.
	currentCapacity := gjson.Get(template, desirecCapacityPath)
	newCapacity := fmt.Sprintf("%d", c.spec.Nodes)
	template, err = sjson.Set(template, desirecCapacityPath, newCapacity)
	if err != nil {
		return errors.Wrap(err, "setting desired capacity")
	}
	currentMaxSize := gjson.Get(template, maxSizePath)
	newMaxSize := fmt.Sprintf("%d", c.spec.MaxNodes)
	template, err = sjson.Set(template, maxSizePath, newMaxSize)
	if err != nil {
		return errors.Wrap(err, "setting max size")
	}
	currentMinSize := gjson.Get(template, minSizePath)
	newMinSize := fmt.Sprintf("%d", c.spec.MinNodes)
	template, err = sjson.Set(template, minSizePath, newMinSize)
	if err != nil {
		return errors.Wrap(err, "setting min size")
	}

	logger.Debug("stack template (post-scale change): %s", template)

	description := fmt.Sprintf("scaling nodegroup, desired capacity from %s to %d, min size from %s to %d and max size from %s to %d",
		currentCapacity.Str, c.spec.Nodes,
		currentMinSize.Str, c.spec.MinNodes,
		currentMaxSize.Str, c.spec.MaxNodes)

	return c.UpdateStack(name, "scale-nodegroup", description, []byte(template), nil)
}

func (c *StackCollection) getStackTemplate(stackName string) (string, error) {
	input := &cfn.GetTemplateInput{
		StackName: aws.String(stackName),
	}

	output, err := c.cfn.GetTemplate(input)
	if err != nil {
		return "", err
	}

	return *output.TemplateBody, nil
}
