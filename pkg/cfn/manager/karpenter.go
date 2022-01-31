package manager

import (
	"strings"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
)

// GetKarpenterStack returns the stack holding the karpenter IAM
// resources
func (c *StackCollection) GetKarpenterStack() (*Stack, error) {
	stacks, err := c.DescribeStacks()
	if err != nil {
		return nil, err
	}

	for _, s := range stacks {
		if *s.StackStatus == cfn.StackStatusDeleteComplete {
			continue
		}
		if isKarpenterStack(s) {
			return s, nil
		}
	}

	return nil, nil
}

func isKarpenterStack(s *Stack) bool {
	return strings.HasSuffix(*s.StackName, "-karpenter")
}
