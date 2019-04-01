package manager

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/pkg/errors"
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
func (c *StackCollection) DeleteTasksForDeprecatedStacks() (*TaskTree, error) {
	stacks, err := c.ListStacks(fmtDeprecatedStacksRegexForCluster(c.spec.Metadata.Name))
	if err != nil {
		return nil, errors.Wrapf(err, "describing deprecated CloudFormation stacks for %q", c.spec.Metadata.Name)
	}
	if len(stacks) == 0 {
		return nil, nil
	}

	deleteControlPlaneTask := &taskWithoutParams{
		info: fmt.Sprintf("delete control plane %q", c.spec.Metadata.Name),
		call: func(errs chan error) error {
			_, err := c.provider.EKS().DescribeCluster(&eks.DescribeClusterInput{
				Name: &c.spec.Metadata.Name,
			})
			if err != nil {
				return err
			}

			_, err = c.provider.EKS().DeleteCluster(&eks.DeleteClusterInput{
				Name: &c.spec.Metadata.Name,
			})
			if err != nil {
				return err
			}

			newRequest := func() *request.Request {
				input := &eks.DescribeClusterInput{
					Name: &c.spec.Metadata.Name,
				}
				req, _ := c.provider.EKS().DescribeClusterRequest(input)
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

			return waiters.Wait(c.spec.Metadata.Name, msg, acceptors, newRequest, c.provider.WaitTimeout(), nil)
		},
	}

	cpStackFound := false
	for _, s := range stacks {
		if strings.HasSuffix(*s.StackName, "-ControlPlane") {
			cpStackFound = true
		}
	}
	tasks := &TaskTree{}

	for _, suffix := range deprecatedStackSuffices() {
		for _, s := range stacks {
			if strings.HasSuffix(*s.StackName, "-"+suffix) {
				if suffix == "-ControlPlane" && !cpStackFound {
					tasks.Append(deleteControlPlaneTask)
				} else {
					tasks.Append(&taskWithStackSpec{
						stack: s,
						call:  c.WaitDeleteStackBySpec,
					})
				}
			}
		}
	}
	return tasks, nil
}
