package manager

import (
	"github.com/weaveworks/eksctl/pkg/cfn/builder"

	"github.com/kubicorn/kubicorn/pkg/logger"
)

func (c *StackCollection) makeClusterStackName() string {
	return "eksctl-" + c.spec.ClusterName + "-cluster"
}

func (c *StackCollection) CreateCluster(errs chan error) error {
	name := c.makeClusterStackName()
	logger.Info("creating cluster stack %q", name)
	stack := builder.NewClusterResourceSet(c.spec)
	if err := stack.AddAllResources(); err != nil {
		return err
	}

	return c.CreateStack(name, stack, nil, errs)
}

func (c *StackCollection) DeleteCluster() error {
	return c.DeleteStack(c.makeClusterStackName())
}
