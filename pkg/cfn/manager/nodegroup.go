package manager

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/awslabs/goformation"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	gfn "github.com/awslabs/goformation/cloudformation"
	"github.com/kubicorn/kubicorn/pkg/logger"
)

func (c *StackCollection) makeNodeGroupStackName(sequence int) string {
	return fmt.Sprintf("eksctl-%s-nodegroup-%d", c.spec.ClusterName, sequence)
}
func (c *StackCollection) CreateInitialNodeGroup(errs chan error) error {
	return c.CreateNodeGroup(0, errs)
}

func (c *StackCollection) CreateNodeGroup(seq int, errs chan error) error {
	name := c.makeNodeGroupStackName(seq)
	logger.Info("creating nodegroup stack %q", name)
	stack := builder.NewNodeGroupResourceSet(c.spec, c.makeClusterStackName(), seq)
	if err := stack.AddAllResources(); err != nil {
		return err
	}

	c.tags = append(c.tags, newTag(NodeGroupTagID, fmt.Sprintf("%d", seq)))

	return c.CreateStack(name, stack, nil, errs)
}

func (c *StackCollection) ScaleNodeGroup(errs chan error) error {
	clusterName := c.makeClusterStackName()
	c.spec.ClusterStackName = clusterName
	name := c.makeNodeGroupStackName(0)
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

	// TODO: In the future this needs to use Goformation but at present we manipulate the
	// JSON directly as the version of goformation we are using isn't handling Refs well
	var f interface{}
	err = json.Unmarshal([]byte(template), &f)
	if err != nil {
		return errors.Wrap(err, "error pasring JSON")
	}

	m := f.(map[string]interface{})
	res := m["Resources"].(map[string]interface{})
	ng := res["NodeGroup"].(map[string]interface{})
	props := ng["Properties"].(map[string]interface{})

	currentCapacity := props["DesiredCapacity"]
	props["DesiredCapacity"] = fmt.Sprintf("%d", c.spec.Nodes)
	currentMaxSize := props["MaxSize"]
	props["MaxSize"] = fmt.Sprintf("%d", c.spec.MaxNodes)
	currentMinSize := props["MinSize"]
	props["MinSize"] = fmt.Sprintf("%d", c.spec.MinNodes)

	updatedTemplate, err := json.Marshal(m)
	if err != nil {
		return errors.Wrapf(err, "rendering template for %q stack", name)
	}
	logger.Debug("stack template (post-scale change): %s", updatedTemplate)

	description := fmt.Sprintf("scaling nodegroup, desired from %s to %d, min size from %s to %d and max size from %s to %d",
		currentCapacity, c.spec.Nodes,
		currentMinSize, c.spec.MinNodes,
		currentMaxSize, c.spec.MaxNodes)

	return c.UpdateStack(name, "scale-nodegroup", description, updatedTemplate, nil, errs)
}

func (c *StackCollection) DeleteNodeGroup() error {
	_, err := c.DeleteStack(c.makeNodeGroupStackName(0))
	return err
}

func (c *StackCollection) WaitDeleteNodeGroup() error {
	return c.WaitDeleteStack(c.makeNodeGroupStackName(0))
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

func (c *StackCollection) getCurrentNodeGroup(templateBody string) (*gfn.AWSAutoScalingAutoScalingGroup, error) {
	template, err := goformation.ParseYAML([]byte(templateBody))
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse CloudFormation template")
	}

	asg, err := template.GetAWSAutoScalingAutoScalingGroupWithName("NodeGroup")

	if err != nil {
		return nil, fmt.Errorf("Unable to find NodeGroup in existing template")
	}

	return &asg, nil
}
