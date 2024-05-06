package addon

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/weaveworks/eksctl/pkg/actions/podidentityassociation"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

type EKSPodIdentityDescriber interface {
	ListPodIdentityAssociations(ctx context.Context, params *eks.ListPodIdentityAssociationsInput, optFns ...func(*eks.Options)) (*eks.ListPodIdentityAssociationsOutput, error)
	DescribePodIdentityAssociation(ctx context.Context, params *eks.DescribePodIdentityAssociationInput, optFns ...func(*eks.Options)) (*eks.DescribePodIdentityAssociationOutput, error)
}

type IAMRoleCreator interface {
	Create(ctx context.Context, podIdentityAssociation *api.PodIdentityAssociation, addonName string) (roleARN string, err error)
}

type IAMRoleUpdater interface {
	Update(ctx context.Context, podIdentityAssociation api.PodIdentityAssociation, stackName string, podIdentityAssociationID string) (string, bool, error)
}

// PodIdentityAssociationUpdater creates or updates IAM resources for pod identities associated with an addon.
type PodIdentityAssociationUpdater struct {
	ClusterName             string
	IAMRoleCreator          IAMRoleCreator
	IAMRoleUpdater          IAMRoleUpdater
	EKSPodIdentityDescriber EKSPodIdentityDescriber
	StackDescriber          podidentityassociation.StackDescriber
}

func (p *PodIdentityAssociationUpdater) UpdateRole(ctx context.Context, podIdentityAssociations []api.PodIdentityAssociation, addonName string) ([]ekstypes.AddonPodIdentityAssociations, error) {
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
				if roleARN, err = p.IAMRoleCreator.Create(ctx, &pia, addonName); err != nil {
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
			stackName := podidentityassociation.MakeAddonPodIdentityStackName(p.ClusterName, addonName, pia.ServiceAccountName)
			hasStack := true
			if _, err := p.StackDescriber.DescribeStack(ctx, &manager.Stack{
				StackName: aws.String(stackName),
			}); err != nil {
				if !manager.IsStackDoesNotExistError(err) {
					return nil, fmt.Errorf("describing IAM resources stack for pod identity association %s: %w", pia.NameString(), err)
				}
				hasStack = false
			}

			roleValidator := &podidentityassociation.RoleUpdateValidator{
				StackDescriber: p.StackDescriber,
			}
			if err := roleValidator.ValidateRoleUpdate(pia, *output.Association, hasStack); err != nil {
				return nil, err
			}
			if hasStack {
				// TODO: if no pod identity has changed, skip update.
				newRoleARN, hasChanged, err := p.IAMRoleUpdater.Update(ctx, pia, stackName, *output.Association.AssociationId)
				if err != nil {
					return nil, err
				}
				if hasChanged {
					roleARN = newRoleARN
				} else {
					roleARN = *output.Association.RoleArn
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
