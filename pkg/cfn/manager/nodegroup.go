package manager

import (
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
	logger.Info("scaling nodegroup stack %q in cluster", name, clusterName)
	//stack := builder.NewNodeGroupResourceSet(c.spec, clusterName, 0)

	// Get current stack
	template, err := c.getStackTemplate(name)
	if err != nil {
		return errors.Wrapf(err, "error getting stack template %s", name)
	}
	logger.Debug("stack template (pre-scale change): %s", template)

	// Get the node group ASG
	stackTemplate, err := goformation.ParseJSON([]byte(template))
	if err != nil {
		return errors.Wrapf(err, "error parsing stack template")
	}
	asg, err := stackTemplate.GetAWSAutoScalingAutoScalingGroupWithName("NodeGroup")
	if err != nil {
		return errors.Wrapf(err, "error finding ASG")
	}

	if c.spec.MinNodes == 0 && c.spec.MaxNodes == 0 {
		c.spec.MinNodes = c.spec.Nodes
		c.spec.MaxNodes = c.spec.Nodes
	}

	asg.DesiredCapacity = gfn.NewString(fmt.Sprintf("%d", c.spec.Nodes))
	asg.MinSize = gfn.NewString(fmt.Sprintf("%d", c.spec.MinNodes))
	asg.MaxSize = gfn.NewString(fmt.Sprintf("%d", c.spec.MaxNodes))

	updatedTemplate, err := stackTemplate.JSON()
	if err != nil {
		return errors.Wrapf(err, "rendering template for %q stack", name)
	}
	logger.Debug("stack template (post-scale change): %s", updatedTemplate)

	return c.UpdateStack(name, updatedTemplate, nil, errs)
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
