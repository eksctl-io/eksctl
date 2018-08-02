package manager

import (
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"

	"github.com/kubicorn/kubicorn/pkg/logger"
)

func (c *StackCollection) makeClusterStackName() string {
	return "eksctl-" + c.spec.ClusterName + "-cluster"
}

func (c *StackCollection) makeClusterStackParams() map[string]string {
	return map[string]string{
		builder.ParamClusterName: c.spec.ClusterName,
	}
}

func (c *StackCollection) CreateCluster(errs chan error) error {
	name := c.makeClusterStackName()
	logger.Info("creating cluster stack %q", name)

	stack := builder.NewClusterResourceSet(c.spec)
	stack.AddAllResources()

	templateBody, err := stack.RenderJSON()
	if err != nil {
		return errors.Wrap(err, "rendering template for cluster stack")
	}

	logger.Debug("templateBody = %s", string(templateBody))

	stackChan := make(chan Stack)
	taskErrs := make(chan error)

	if err := c.CreateStack(name, templateBody, c.makeClusterStackParams(), true, stackChan, taskErrs); err != nil {
		return err
	}

	go func() {
		defer close(errs)
		defer close(stackChan)

		if err := <-taskErrs; err != nil {
			errs <- err
			return
		}

		if err := stack.GetAllOutputs(<-stackChan); err != nil {
			errs <- errors.Wrap(err, "getting cluster stack outputs")
		}

		logger.Debug("clusterConfig = %#v", c.spec)
		logger.Success("created cluster stack %q", name)

		errs <- nil
	}()
	return nil
}

func (c *StackCollection) DeleteCluster() error {
	return c.DeleteStack(c.makeClusterStackName())
}
