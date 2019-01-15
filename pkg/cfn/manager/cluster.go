package manager

import (
	"strings"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/kris-nova/logger"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha3"
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

	// Unlike with `CreateNodeGroup`, all tags are already set for the cluster stack
	return c.CreateStack(name, stack, nil, nil, errs)
}

// DescribeClusterStack calls DescribeStacks and filters out cluster stack
func (c *StackCollection) DescribeClusterStack() (*Stack, error) {
	stacks, err := c.DescribeStacks()
	if err != nil {
		return nil, err
	}

	for _, s := range stacks {
		if *s.StackStatus == cfn.StackStatusDeleteComplete {
			continue
		}
		if getClusterName(s) != "" {
			return s, nil
		}
	}
	return nil, nil
}

// DeleteCluster deletes the cluster
func (c *StackCollection) DeleteCluster() error {
	_, err := c.DeleteStack(c.makeClusterStackName())
	return err
}

// WaitDeleteCluster waits till the cluster is deleted
func (c *StackCollection) WaitDeleteCluster() error {
	return c.BlockingWaitDeleteStack(c.makeClusterStackName())
}

func getClusterName(s *Stack) string {
	for _, tag := range s.Tags {
		if *tag.Key == api.ClusterNameTag {
			if strings.HasSuffix(*s.StackName, "-cluster") {
				return *tag.Value
			}
		}
	}

	if strings.HasPrefix(*s.StackName, "EKS-") && strings.HasSuffix(*s.StackName, "-ControlPlane") {
		return strings.TrimPrefix("EKS-", strings.TrimSuffix(*s.StackName, "-ControlPlane"))
	}
	return ""
}
