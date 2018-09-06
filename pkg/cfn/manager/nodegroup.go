package manager

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
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
	stack := builder.NewNodeGroupResourceSet(c.spec)

	// Get current stack
	template, err := c.getStackTemplate(name)
	if err != nil {
		return errors.Wrapf(err, "error getting stack template %s", name)
	}

	if err := stack.AddResourcesForScaling(template); err != nil {
		return nil
	}

	return c.UpdateStack(name, stack, c.makeNodeGroupParams(0), errs)
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

	logger.Debug("retrieved template: %s", *output.TemplateBody)
	return *output.TemplateBody, nil
}
