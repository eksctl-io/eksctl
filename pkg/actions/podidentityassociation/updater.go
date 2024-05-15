package podidentityassociation

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/utils/apierrors"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

// An Updater updates pod identity associations.
type Updater struct {
	// ClusterName is the cluster name.
	ClusterName string
	// StackUpdater updates stacks.
	StackUpdater StackUpdater
	// APIDeleter updates pod identity associations using the EKS API.
	APIUpdater APIUpdater
}

// A StackUpdater updates CloudFormation stacks.
//
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_stack_updater.go . StackUpdater
type StackUpdater interface {
	StackLister
	// MustUpdateStack updates the CloudFormation stack.
	MustUpdateStack(ctx context.Context, options manager.UpdateStackOptions) error
}

// APIUpdater updates pod identity associations using the EKS API.
type APIUpdater interface {
	APILister
	DescribePodIdentityAssociation(ctx context.Context, params *eks.DescribePodIdentityAssociationInput, optFns ...func(*eks.Options)) (*eks.DescribePodIdentityAssociationOutput, error)
	UpdatePodIdentityAssociation(ctx context.Context, params *eks.UpdatePodIdentityAssociationInput, optFns ...func(*eks.Options)) (*eks.UpdatePodIdentityAssociationOutput, error)
}

// UpdateConfig holds configuration for updating a pod identity association.
type UpdateConfig struct {
	PodIdentityAssociation api.PodIdentityAssociation
	AssociationID          string
	HasIAMResourcesStack   bool
	StackName              string
}

// Update updates the specified pod identity associations.
func (u *Updater) Update(ctx context.Context, podIdentityAssociations []api.PodIdentityAssociation) error {
	roleStackNames, err := u.StackUpdater.ListPodIdentityStackNames(ctx)
	if err != nil {
		return fmt.Errorf("error listing stack names for pod identity associations: %w", err)
	}
	taskTree := &tasks.TaskTree{
		Parallel: true,
	}
	for _, p := range podIdentityAssociations {
		podIdentityAssociationID := Identifier{
			Namespace:          p.Namespace,
			ServiceAccountName: p.ServiceAccountName,
		}.IDString()
		updateErr := func(err error) error {
			return fmt.Errorf("error updating pod identity association %q: %w", podIdentityAssociationID, err)
		}
		updateConfig, err := u.makeUpdate(ctx, p, roleStackNames)
		if err != nil {
			return updateErr(err)
		}
		taskTree.Append(&tasks.GenericTask{
			Description: fmt.Sprintf("update pod identity association %s", podIdentityAssociationID),
			Doer: func() error {
				if err := u.update(ctx, updateConfig, podIdentityAssociationID); err != nil {
					return updateErr(err)
				}
				return nil
			},
		})
	}
	return runAllTasks(taskTree)
}

func (u *Updater) update(ctx context.Context, updateConfig *UpdateConfig, podIdentityAssociationID string) error {
	roleARN := updateConfig.PodIdentityAssociation.RoleARN
	if updateConfig.HasIAMResourcesStack {
		roleUpdater := &IAMRoleUpdater{
			StackUpdater: u.StackUpdater,
		}
		newRoleARN, hasChanged, err := roleUpdater.Update(ctx, updateConfig.PodIdentityAssociation, updateConfig.StackName, podIdentityAssociationID)
		if err != nil {
			return err
		}
		if !hasChanged {
			return nil
		}
		roleARN = newRoleARN
	}
	return u.updatePodIdentityAssociation(ctx, roleARN, updateConfig, podIdentityAssociationID)
}

func (u *Updater) updatePodIdentityAssociation(ctx context.Context, roleARN string, updateConfig *UpdateConfig, podIdentityAssociationID string) error {
	if _, err := u.APIUpdater.UpdatePodIdentityAssociation(ctx, &eks.UpdatePodIdentityAssociationInput{
		AssociationId: aws.String(updateConfig.AssociationID),
		ClusterName:   aws.String(u.ClusterName),
		RoleArn:       aws.String(roleARN),
	}); err != nil {
		return fmt.Errorf("(associationID: %s, roleARN: %s): %w", updateConfig.AssociationID, roleARN, err)
	}
	logger.Info("updated role ARN %q for pod identity association %q", roleARN, podIdentityAssociationID)
	return nil
}

func (u *Updater) makeUpdate(ctx context.Context, pia api.PodIdentityAssociation, roleStackNames []string) (*UpdateConfig, error) {
	const notFoundErrMsg = "pod identity association does not exist"
	output, err := u.APIUpdater.ListPodIdentityAssociations(ctx, &eks.ListPodIdentityAssociationsInput{
		ClusterName:    aws.String(u.ClusterName),
		Namespace:      aws.String(pia.Namespace),
		ServiceAccount: aws.String(pia.ServiceAccountName),
	})
	if err != nil {
		if apierrors.IsNotFoundError(err) {
			return nil, fmt.Errorf("%s: %w", notFoundErrMsg, err)
		}
		return nil, fmt.Errorf("error listing pod identity associations: %w", err)
	}
	switch len(output.Associations) {
	default:
		return nil, fmt.Errorf("expected to find only 1 pod identity association; got %d", len(output.Associations))
	case 0:
		return nil, errors.New(notFoundErrMsg)
	case 1:
		association := output.Associations[0]
		if association.OwnerArn != nil {
			return nil, fmt.Errorf("cannot update podidentityassociation %s as it is in use by addon %s; "+
				"please use `eksctl update addon` instead", pia.NameString(), *association.OwnerArn)
		}
		describeOutput, err := u.APIUpdater.DescribePodIdentityAssociation(ctx, &eks.DescribePodIdentityAssociationInput{
			ClusterName:   aws.String(u.ClusterName),
			AssociationId: association.AssociationId,
		})
		if err != nil {
			return nil, fmt.Errorf("error describing pod identity association: %w", err)
		}
		stackName, hasStack := getIAMResourcesStack(roleStackNames, Identifier{
			Namespace:          pia.Namespace,
			ServiceAccountName: pia.ServiceAccountName,
		})
		updateValidator := &RoleUpdateValidator{
			StackDescriber: u.StackUpdater,
		}
		if err := updateValidator.ValidateRoleUpdate(pia, *describeOutput.Association, hasStack); err != nil {
			return nil, err
		}
		return &UpdateConfig{
			PodIdentityAssociation: pia,
			AssociationID:          *describeOutput.Association.AssociationId,
			HasIAMResourcesStack:   hasStack,
			StackName:              stackName,
		}, nil
	}
}

type StackDescriber interface {
	DescribeStack(context.Context, *manager.Stack) (*manager.Stack, error)
}

type RoleUpdateValidator struct {
	StackDescriber StackDescriber
}

// ValidateRoleUpdate validates the role associated with pia.
func (r *RoleUpdateValidator) ValidateRoleUpdate(pia api.PodIdentityAssociation, association ekstypes.PodIdentityAssociation, hasStack bool) error {
	if hasStack {
		if association.RoleArn != nil && pia.RoleARN != "" && pia.RoleARN != *association.RoleArn {
			return errors.New("cannot change podIdentityAssociation.roleARN since the role was created by eksctl")
		}
	} else {
		if pia.RoleARN == "" {
			return errors.New("podIdentityAssociation.roleARN is required since the role was not created by eksctl")
		}
		podIDWithRoleARN := api.PodIdentityAssociation{
			Namespace:          pia.Namespace,
			ServiceAccountName: pia.ServiceAccountName,
			RoleARN:            pia.RoleARN,
		}
		if !reflect.DeepEqual(pia, podIDWithRoleARN) {
			return errors.New("only namespace, serviceAccountName and roleARN can be specified if the role was not created by eksctl")
		}
	}
	return nil
}
