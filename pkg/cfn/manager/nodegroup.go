package manager

import (
	"fmt"

	"github.com/weaveworks/eksctl/pkg/cfn/builder"

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

func (c *StackCollection) DeleteNodeGroup() error {
	return c.DeleteStack(c.makeNodeGroupStackName(0))
}

func (c *StackCollection) WaitDeleteNodeGroup() error {
	return c.WaitDeleteStack(c.makeNodeGroupStackName(0))
}
