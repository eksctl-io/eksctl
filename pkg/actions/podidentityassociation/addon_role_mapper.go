package podidentityassociation

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
)

// AddonServiceAccountRoleMapper maps service account role ARNs to EKS addons.
type AddonServiceAccountRoleMapper map[string]*ekstypes.Addon

// CreateAddonServiceAccountRoleMapper creates an AddonServiceAccountRoleMapper that maps service account role ARNs to EKS addons.
func CreateAddonServiceAccountRoleMapper(ctx context.Context, clusterName string, eksAddonsAPI EKSAddonsAPI) (AddonServiceAccountRoleMapper, error) {
	addonMapper := AddonServiceAccountRoleMapper{}
	paginator := eks.NewListAddonsPaginator(eksAddonsAPI, &eks.ListAddonsInput{
		ClusterName: aws.String(clusterName),
	})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("listing addons: %w", err)
		}
		for _, addonName := range output.Addons {
			addon, err := eksAddonsAPI.DescribeAddon(ctx, &eks.DescribeAddonInput{
				ClusterName: aws.String(clusterName),
				AddonName:   aws.String(addonName),
			})
			if err != nil {
				return nil, err
			}
			if roleARN := addon.Addon.ServiceAccountRoleArn; roleARN != nil {
				addonMapper[*roleARN] = addon.Addon
			}
		}
	}
	return addonMapper, nil
}

// AddonForServiceAccountRole returns the addon used by roleARN.
func (m AddonServiceAccountRoleMapper) AddonForServiceAccountRole(roleARN string) *ekstypes.Addon {
	return m[roleARN]
}
