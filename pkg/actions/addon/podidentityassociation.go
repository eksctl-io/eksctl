package addon

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/kris-nova/logger"

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
	StackDeleter            podidentityassociation.StackDeleter
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
			return nil, fmt.Errorf("expected to find exactly 1 pod identity association for %s; got %d", pia.NameString(), len(output.Associations))
		case 0:
			// Create IAM resources.
			if roleARN == "" {
				var err error
				if roleARN, err = p.IAMRoleCreator.Create(ctx, &pia, addonName); err != nil {
					return nil, err
				}
				stack, err := p.getStack(ctx, manager.MakeAddonStackName(p.ClusterName, addonName), pia.ServiceAccountName)
				if err != nil {
					return nil, fmt.Errorf("getting old IRSA stack for addon %s: %w", addonName, err)
				}
				if stack != nil {
					logger.Info("deleting old IRSA stack for addon %s", addonName)
					if err := p.deleteStack(ctx, stack); err != nil {
						return nil, fmt.Errorf("deleting old IRSA stack for addon %s: %w", addonName, err)
					}
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
			stack, err := p.getAddonStack(ctx, addonName, pia.ServiceAccountName)
			if err != nil {
				return nil, fmt.Errorf("getting IAM resources stack for addon %s with pod identity association %s: %w", addonName, pia.NameString(), err)
			}

			roleValidator := &podidentityassociation.RoleUpdateValidator{
				StackDescriber: p.StackDeleter,
			}
			hasStack := stack != nil
			if err := roleValidator.ValidateRoleUpdate(pia, *output.Association, hasStack); err != nil {
				return nil, err
			}
			if hasStack {
				// TODO: if no pod identity has changed, skip update.
				newRoleARN, hasChanged, err := p.IAMRoleUpdater.Update(ctx, pia, *stack.StackName, *output.Association.AssociationId)
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

func (p *PodIdentityAssociationUpdater) getAddonStack(ctx context.Context, addonName, serviceAccount string) (*manager.Stack, error) {
	for _, stackName := range []string{podidentityassociation.MakeAddonPodIdentityStackName(p.ClusterName, addonName, serviceAccount),
		manager.MakeAddonStackName(p.ClusterName, addonName)} {
		stack, err := p.getStack(ctx, stackName, serviceAccount)
		if err != nil {
			return nil, err
		}
		if stack != nil {
			return stack, nil
		}
	}
	return nil, nil
}

func (p *PodIdentityAssociationUpdater) getStack(ctx context.Context, stackName, serviceAccount string) (*manager.Stack, error) {
	switch stack, err := p.StackDeleter.DescribeStack(ctx, &manager.Stack{
		StackName: aws.String(stackName),
	}); {
	case err == nil:
		return stack, nil
	case manager.IsStackDoesNotExistError(err):
		return nil, nil
	default:
		return nil, fmt.Errorf("describing IAM resources stack for service account %s: %w", serviceAccount, err)
	}
}

func (p *PodIdentityAssociationUpdater) DeleteRole(ctx context.Context, addonName, serviceAccountName string) (bool, error) {
	stack, err := p.getAddonStack(ctx, addonName, serviceAccountName)
	if err != nil {
		return false, fmt.Errorf("getting IAM resources stack for addon %s with service account %s: %w", addonName, serviceAccountName, err)
	}
	if err := p.deleteStack(ctx, stack); err != nil {
		return false, err
	}
	return true, nil
}

func (p *PodIdentityAssociationUpdater) deleteStack(ctx context.Context, stack *manager.Stack) error {
	errCh := make(chan error)
	if err := p.StackDeleter.DeleteStackBySpecSync(ctx, stack, errCh); err != nil {
		return fmt.Errorf("deleting stack %s: %w", *stack.StackName, err)
	}
	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("deleting stack %s: %w", *stack.StackName, err)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("timed out waiting for deletion of stack %s: %w", *stack.StackName, ctx.Err())
	}
}
