package manager

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
	"github.com/weaveworks/eksctl/pkg/utils/waiters"
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

// DeleteTasksForDeprecatedStacks all deprecated stacks
func (c *StackCollection) DeleteTasksForDeprecatedStacks() (*tasks.TaskTree, error) {
	stacks, err := c.ListStacksMatching(fmtDeprecatedStacksRegexForCluster(c.spec.Metadata.Name))
	if err != nil {
		return nil, errors.Wrapf(err, "describing deprecated CloudFormation stacks for %q", c.spec.Metadata.Name)
	}
	if len(stacks) == 0 {
		return nil, nil
	}

	deleteControlPlaneTask := &tasks.TaskWithoutParams{
		Info: fmt.Sprintf("delete control plane %q", c.spec.Metadata.Name),
		Call: func(errs chan error) error {
			_, err := c.eksAPI.DescribeCluster(&eks.DescribeClusterInput{
				Name: &c.spec.Metadata.Name,
			})
			if err != nil {
				return err
			}

			_, err = c.eksAPI.DeleteCluster(&eks.DeleteClusterInput{
				Name: &c.spec.Metadata.Name,
			})
			if err != nil {
				return err
			}

			newRequest := func() *request.Request {
				input := &eks.DescribeClusterInput{
					Name: &c.spec.Metadata.Name,
				}
				req, _ := c.eksAPI.DescribeClusterRequest(input)
				return req
			}

			msg := fmt.Sprintf("waiting for control plane %q to be deleted", c.spec.Metadata.Name)

			acceptors := waiters.MakeAcceptors(
				"Cluster.Status",
				eks.ClusterStatusDeleting,
				[]string{
					eks.ClusterStatusFailed,
				},
			)

			return waiters.Wait(c.spec.Metadata.Name, msg, acceptors, newRequest, c.waitTimeout, nil)
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
