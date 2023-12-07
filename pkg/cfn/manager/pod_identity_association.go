package manager

import (
	"context"
	"fmt"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

func (c *StackCollection) ListPodIdentityStackNames(ctx context.Context) ([]string, error) {
	names := []string{}
	stacks, err := c.ListStacks(ctx)
	if err != nil {
		return names, fmt.Errorf("listing stacks: %w", err)
	}

	for _, s := range stacks {
		isPodIdentityStack := false
		for _, tag := range s.Tags {
			if *tag.Key == api.PodIdentityAssociationNameTag {
				isPodIdentityStack = true
			}
		}
		if isPodIdentityStack {
			names = append(names, *s.StackName)
		}
	}

	return names, nil
}
