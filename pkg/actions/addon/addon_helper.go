package addon

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// Helper is a helper for validating nodegroup creation.
type Helper struct {
	// ClusterName holds the cluster name.
	ClusterName string
	// Lister lists addons.
	Lister Lister
}

// A Lister lists addons.
type Lister interface {
	ListAddons(ctx context.Context, params *eks.ListAddonsInput, optFns ...func(*eks.Options)) (*eks.ListAddonsOutput, error)
}

// ValidateNodeGroupCreation validates whether the cluster has core networking addons.
func (a *Helper) ValidateNodeGroupCreation(ctx context.Context) error {
	output, err := a.Lister.ListAddons(ctx, &eks.ListAddonsInput{
		ClusterName: aws.String(a.ClusterName),
	})
	if err != nil {
		return fmt.Errorf("listing addons: %w", err)
	}
	if !api.HasAllDefaultAddons(output.Addons) {
		return fmt.Errorf("core networking addons which are required to create nodegroups are missing in the cluster; " +
			"please create them using `eksctl create addon`")
	}
	return nil
}
