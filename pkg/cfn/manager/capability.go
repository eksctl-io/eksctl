package manager

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// GetIAMCapabilitiesStacks returns all capability IAM stacks for the cluster.
func (c *StackCollection) GetIAMCapabilitiesStacks(ctx context.Context) ([]*Stack, error) {
	stacks, err := c.ListStacks(ctx)
	if err != nil {
		return nil, err
	}

	iamCapabilityStacks := []*Stack{}
	for _, s := range stacks {
		if s.StackStatus == types.StackStatusDeleteComplete {
			continue
		}
		if GetIAMCapabilityName(s) != "" {
			iamCapabilityStacks = append(iamCapabilityStacks, s)
		}
	}
	return iamCapabilityStacks, nil
}

// GetIAMCapabilityName returns the capability name for stack.
func GetIAMCapabilityName(stack *types.Stack) string {
	for _, tag := range stack.Tags {
		if *tag.Key == api.CapabilityNameTag {
			return *tag.Value
		}
	}
	return ""
}
