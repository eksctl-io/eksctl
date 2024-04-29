package addon

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/weaveworks/eksctl/pkg/actions/podidentityassociation"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

type PodIdentityStackLister interface {
	ListPodIdentityStackNames(ctx context.Context) ([]string, error)
}

type EKSPodIdentityDescriber interface {
	ListPodIdentityAssociations(ctx context.Context, params *eks.ListPodIdentityAssociationsInput, optFns ...func(*eks.Options)) (*eks.ListPodIdentityAssociationsOutput, error)
	DescribePodIdentityAssociation(ctx context.Context, params *eks.DescribePodIdentityAssociationInput, optFns ...func(*eks.Options)) (*eks.DescribePodIdentityAssociationOutput, error)
}

type IAMRoleCreator interface {
	Create(ctx context.Context, podIdentityAssociation *api.PodIdentityAssociation) (roleARN string, err error)
}

type IAMRoleUpdater interface {
	Update(ctx context.Context, updateConfig *podidentityassociation.UpdateConfig, podIdentityAssociationID string) (roleARN string, hasChanged bool, err error)
}

// PodIdentityAssociationUpdater creates or updates IAM resources for pod identities associated with an addon.
type PodIdentityAssociationUpdater struct {
	ClusterName             string
	IAMRoleCreator          IAMRoleCreator
	IAMRoleUpdater          IAMRoleUpdater
	PodIdentityStackLister  PodIdentityStackLister
	EKSPodIdentityDescriber EKSPodIdentityDescriber
}

// TODO

func (p *PodIdentityAssociationUpdater) UpdateRole(ctx context.Context, podIdentityAssociations []api.PodIdentityAssociation) ([]ekstypes.AddonPodIdentityAssociations, error) {
	var addonPodIdentityAssociations []ekstypes.AddonPodIdentityAssociations
	for _, pia := range podIdentityAssociations {
		output, err := p.EKSPodIdentityDescriber.ListPodIdentityAssociations(ctx, &eks.ListPodIdentityAssociationsInput{
			ClusterName:    aws.String(p.ClusterName),
			Namespace:      aws.String(pia.Namespace),
			ServiceAccount: aws.String(pia.ServiceAccountName),
		})
		if err != nil {
			return nil, fmt.Errorf("listing pod identity associations: %w", err)
		}
		roleARN := pia.RoleARN
		switch len(output.Associations) {
		default:
			// TODO: does the API return a not found error if no association exists?
			return nil, fmt.Errorf("expected to find exactly 1 pod identity association for %s; got %d", pia.NameString(), len(output.Associations))
		case 0:
			// Create IAM resources.
			if roleARN == "" {
				var err error
				if roleARN, err = p.IAMRoleCreator.Create(ctx, &pia); err != nil {
					return nil, err
				}
			}
		case 1:
			// Update IAM resources if required.
			output, err := p.EKSPodIdentityDescriber.DescribePodIdentityAssociation(ctx, &eks.DescribePodIdentityAssociationInput{
				ClusterName:   aws.String(p.ClusterName),
				AssociationId: output.Associations[0].AssociationId,
			})
			if err != nil {
				return nil, err
			}
			// TODO: avoid repeating this call.
			roleStackNames, err := p.PodIdentityStackLister.ListPodIdentityStackNames(ctx)
			if err != nil {
				return nil, fmt.Errorf("error listing stack names for pod identity associations: %w", err)
			}
			updateConfig, err := podidentityassociation.MakeRoleUpdateConfig(pia, *output.Association, roleStackNames)
			if err != nil {
				return nil, err
			}
			if updateConfig.HasIAMResourcesStack {
				// TODO: if no pod identity has changed, skip update?
				if roleARN, _, err = p.IAMRoleUpdater.Update(ctx, updateConfig, *output.Association.AssociationId); err != nil {
					return nil, err
				}
			}
		}
		addonPodIdentityAssociations = append(addonPodIdentityAssociations, ekstypes.AddonPodIdentityAssociations{
			RoleArn:        aws.String(roleARN),
			ServiceAccount: aws.String(pia.ServiceAccountName),
		})
	}
	return addonPodIdentityAssociations, nil
}
