package manager

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
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
)

func (c *StackCollection) makeNodeGroupStackName(id int) string {
	return fmt.Sprintf("eksctl-%s-nodegroup-%d", c.spec.Metadata.Name, id)
}

// CreateNodeGroup creates the nodegroup
func (c *StackCollection) CreateNodeGroup(errs chan error, data interface{}) error {
	ng := data.(*api.NodeGroup)
	name := c.makeNodeGroupStackName(ng.ID)
	logger.Info("creating nodegroup stack %q", name)
	stack := builder.NewNodeGroupResourceSet(c.spec, c.makeClusterStackName(), ng.ID)
	if err := stack.AddAllResources(); err != nil {
		return err
	}

	c.tags = append(c.tags, newTag(NodeGroupIDTag, fmt.Sprintf("%d", ng.ID)))

	for k, v := range ng.Tags {
		c.tags = append(c.tags, newTag(k, v))
	}

	return c.CreateStack(name, stack, nil, errs)
}

func (c *StackCollection) listAllNodeGroups() ([]string, error) {
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
	defer close(errs)
	name := data.(string)
	_, err := c.DeleteStack(name)
	return err
}

// WaitDeleteNodeGroup waits until the nodegroup is deleted
func (c *StackCollection) WaitDeleteNodeGroup(errs chan error, data interface{}) error {
	defer close(errs)
	name := data.(string)
	return c.WaitDeleteStack(name)
}

// ScaleInitialNodeGroup will scale the first nodegroup (ID: 0)
func (c *StackCollection) ScaleInitialNodeGroup() error {
	return c.ScaleNodeGroup(0)
}

// ScaleNodeGroup will scale an existing nodegroup
func (c *StackCollection) ScaleNodeGroup(id int) error {
	ng := c.spec.NodeGroups[id]
	clusterName := c.makeClusterStackName()
	c.spec.ClusterStackName = clusterName
	name := c.makeNodeGroupStackName(id)
	logger.Info("scaling nodegroup stack %q in cluster %s", name, clusterName)

	// Get current stack
	template, err := c.getStackTemplate(name)
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

func (c *StackCollection) getStackTemplate(stackName string) (string, error) {
	input := &cfn.GetTemplateInput{
		StackName: aws.String(stackName),
	}

	output, err := c.provider.CloudFormation().GetTemplate(input)
	if err != nil {
		return "", err
	}

	return *output.TemplateBody, nil
}
