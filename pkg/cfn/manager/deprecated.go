package manager

import (
	"context"
	"fmt"
	"strings"

	"github.com/kris-nova/logger"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/pkg/errors"

	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

func deprecatedStackSuffices() []string {
	return []string{
		"DefaultNodeGroup",
		"ControlPlane",
		"ServiceRole",
		"VPC",
	}
}
func fmtDeprecatedStacksRegexForCluster(name string) string {
	const ourStackRegexFmt = "^EKS-%s-(%s)$"
	return fmt.Sprintf(ourStackRegexFmt, name, strings.Join(deprecatedStackSuffices(), "|"))
}

// DeleteTasksForDeprecatedStacks deletes all deprecated stacks.
func (c *StackCollection) DeleteTasksForDeprecatedStacks(ctx context.Context) (*tasks.TaskTree, error) {
	stacks, err := c.ListStacksMatching(ctx, fmtDeprecatedStacksRegexForCluster(c.spec.Metadata.Name))
	if err != nil {
		return nil, errors.Wrapf(err, "describing deprecated CloudFormation stacks for %q", c.spec.Metadata.Name)
	}
	if len(stacks) == 0 {
		return nil, nil
	}

	deleteControlPlaneTask := &tasks.TaskWithoutParams{
		Info: fmt.Sprintf("delete control plane %q", c.spec.Metadata.Name),
		Call: func(errs chan error) error {
			describeClusterInput := &eks.DescribeClusterInput{
				Name: &c.spec.Metadata.Name,
			}
			_, err := c.eksAPI.DescribeCluster(ctx, describeClusterInput)
			if err != nil {
				return err
			}

			_, err = c.eksAPI.DeleteCluster(ctx, &eks.DeleteClusterInput{
				Name: &c.spec.Metadata.Name,
			})
			if err != nil {
				return err
			}

			logger.Info("waiting for control plane %q to be deleted", c.spec.Metadata.Name)
			waiter := eks.NewClusterDeletedWaiter(c.eksAPI)
			return waiter.Wait(ctx, describeClusterInput, c.waitTimeout)
		},
	}

	cpStackFound := false
	for _, s := range stacks {
		if strings.HasSuffix(*s.StackName, "-ControlPlane") {
			cpStackFound = true
		}
	}
	taskTree := &tasks.TaskTree{}

	for _, suffix := range deprecatedStackSuffices() {
		for _, s := range stacks {
			if strings.HasSuffix(*s.StackName, "-"+suffix) {
				if suffix == "-ControlPlane" && !cpStackFound {
					taskTree.Append(deleteControlPlaneTask)
				} else {
					taskTree.Append(&taskWithStackSpec{
						stack: s,
						call:  c.DeleteStackBySpecSync,
					})
				}
			}
		}
	}
	return taskTree, nil
}
