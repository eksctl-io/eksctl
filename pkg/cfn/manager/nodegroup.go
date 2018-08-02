package manager

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"

	"github.com/kubicorn/kubicorn/pkg/logger"
)

func (c *StackCollection) makeNodeGroupStackName(sequence int) string {
	return fmt.Sprintf("eksctl-%s-nodegroup-%d", c.spec.ClusterName, sequence)
}

func (c *StackCollection) makeNodeGroupParams(sequence int) map[string]string {
	return map[string]string{
		builder.ParamClusterName:      c.spec.ClusterName,
		builder.ParamClusterStackName: c.spec.ClusterStackName,
		builder.ParamNodeGroupID:      fmt.Sprintf("%d", sequence),
	}
}

func (c *StackCollection) CreateNodeGroup(errs chan error) error {
	seq := 0
	name := c.makeNodeGroupStackName(seq)
	logger.Info("creating nodegroup stack %q", name)

	stack := builder.NewNodeGroupResourceSet(c.spec)
	if err := stack.AddAllResources(); err != nil {
		return err
	}

	templateBody, err := stack.RenderJSON()
	if err != nil {
		return errors.Wrap(err, "rendering template for nodegroup stack")
	}

	logger.Debug("templateBody = %s", string(templateBody))

	stackChan := make(chan Stack)
	taskErrs := make(chan error)

	if err := c.CreateStack(name, templateBody, c.makeNodeGroupParams(seq), true, stackChan, taskErrs); err != nil {
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
			errs <- errors.Wrap(err, "getting nodegroup stack outputs")
		}

		logger.Debug("clusterConfig = %#v", c.spec)
		logger.Success("created nodegroup stack %q", name)

		errs <- nil
	}()
	return nil
}

func (c *StackCollection) DeleteNodeGroup() error {
	return c.DeleteStack(c.makeNodeGroupStackName(0))
}
