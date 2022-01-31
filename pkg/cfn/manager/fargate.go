package manager

import (
	"strings"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
)

// GetFargateStack returns the stack holding the fargate IAM
// resources, if any
func (c *StackCollection) GetFargateStack() (*Stack, error) {
	stacks, err := c.DescribeStacks()
	if err != nil {
		return nil, err
	}

	for _, s := range stacks {
		if *s.StackStatus == cfn.StackStatusDeleteComplete {
			continue
		}
		if isFargateStack(s) {
			return s, nil
		}
	}

	return nil, nil
}

func isFargateStack(s *Stack) bool {
	return strings.HasSuffix(*s.StackName, "-fargate")
}
