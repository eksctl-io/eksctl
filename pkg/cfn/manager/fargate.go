package manager

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

// GetFargateStack returns the stack holding the fargate IAM
// resources, if any
func (c *StackCollection) GetFargateStack(ctx context.Context) (*Stack, error) {
	stacks, err := c.ListStacks(ctx)
	if err != nil {
		return nil, err
	}

	for _, s := range stacks {
		if s.StackStatus == types.StackStatusDeleteComplete {
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
