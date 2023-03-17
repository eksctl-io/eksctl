package manager

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

// GetKarpenterStack returns the stack holding the karpenter IAM
// resources
func (c *StackCollection) GetKarpenterStack(ctx context.Context) (*Stack, error) {
	stacks, err := c.ListStacks(ctx)
	if err != nil {
		return nil, err
	}

	for _, s := range stacks {
		if s.StackStatus == types.StackStatusDeleteComplete {
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
