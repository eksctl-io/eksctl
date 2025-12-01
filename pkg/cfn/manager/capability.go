package manager

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// GetIAMCapabilitiesStacks returns all capability IAM stacks for the cluster.
func (c *StackCollection) ListCapabilitiesIAMStacks(ctx context.Context) ([]*Stack, error) {
	stacks, err := c.ListStacks(ctx)
	if err != nil {
		return nil, err
	}

	iamCapabilityStacks := []*Stack{}
	for _, s := range stacks {
		if s.StackStatus == types.StackStatusDeleteComplete {
			continue
		}
		if GetCapabilityNameFromIAMStack(s) != "" {
			iamCapabilityStacks = append(iamCapabilityStacks, s)
		}
	}
	return iamCapabilityStacks, nil
}

// ListCapabilityStacks returns capability stacks for the cluster using tags
func (c *StackCollection) ListCapabilityStacks(ctx context.Context) ([]*Stack, error) {
	stacks, err := c.ListStacks(ctx)
	if err != nil {
		return nil, err
	}

	capabilityStacks := []*Stack{}
	for _, s := range stacks {
		if s.StackStatus == types.StackStatusDeleteComplete {
			continue
		}
		if GetCapabilityNameFromStack(s) != "" {
			capabilityStacks = append(capabilityStacks, s)
		}
	}
	return capabilityStacks, nil
}

// GetIAMCapabilityName returns the capability name for stack.
func GetCapabilityNameFromIAMStack(stack *types.Stack) string {
	for _, tag := range stack.Tags {
		if *tag.Key == api.CapabilityIAMRoleTag {
			return *tag.Value
		}
	}
	return ""
}

// GetIAMCapabilityName returns the capability name for stack.
func GetCapabilityNameFromStack(stack *types.Stack) string {
	for _, tag := range stack.Tags {
		if *tag.Key == api.CapabilityNameTag {
			return *tag.Value
		}
	}
	return ""
}

// // ListCapabilityStackNames returns a list of capability stack names for the cluster
// func (c *StackCollection) ListCapabilityStackNames(ctx context.Context, clusterName string) ([]string, error) {
// 	stacks, err := c.ListCapabilityStacks(ctx, clusterName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var stackNames []string
// 	for _, stack := range stacks {
// 		if stack.StackName != nil {
// 			stackNames = append(stackNames, *stack.StackName)
// 		}
// 	}
// 	return stackNames, nil
// }
