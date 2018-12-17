package manager

import (
	"github.com/kris-nova/logger"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
)

func (c *StackCollection) makeClusterStackName() string {
	return "eksctl-" + c.spec.Metadata.Name + "-cluster"
}

// CreateCluster creates the cluster
func (c *StackCollection) CreateCluster(errs chan error, _ interface{}) error {
	name := c.makeClusterStackName()
	logger.Info("creating cluster stack %q", name)
	stack := builder.NewClusterResourceSet(c.provider, c.spec)
	if err := stack.AddAllResources(); err != nil {
		return err
	}

	return c.CreateStack(name, stack, nil, errs)
}

// DeleteCluster deletes the cluster
func (c *StackCollection) DeleteCluster() error {
	_, err := c.DeleteStack(c.makeClusterStackName())
	return err
}

// WaitDeleteCluster waits till the cluster is deleted
func (c *StackCollection) WaitDeleteCluster() error {
	return c.WaitDeleteStack(c.makeClusterStackName())
}
