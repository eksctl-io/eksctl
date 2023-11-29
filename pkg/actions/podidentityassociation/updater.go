package podidentityassociation

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"golang.org/x/exp/slices"

	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"

	"github.com/kris-nova/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
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

type updateConfig struct {
	podIdentityAssociation api.PodIdentityAssociation
	associationID          string
	hasIAMResourcesStack   bool
	stackName              string
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

func (u *Updater) update(ctx context.Context, updateConfig *updateConfig, podIdentityAssociationID string) error {
	if !updateConfig.hasIAMResourcesStack {
		return u.updatePodIdentityAssociation(ctx, updateConfig, podIdentityAssociationID)
	}

	stack, err := u.StackUpdater.DescribeStack(ctx, &manager.Stack{
		StackName: aws.String(updateConfig.stackName),
	})
	if err != nil {
		return fmt.Errorf("describing IAM resources stack %q: %w", updateConfig.stackName, err)
	}
	if updateConfig.podIdentityAssociation.RoleName != "" && !slices.Contains(stack.Capabilities, cfntypes.CapabilityCapabilityNamedIam) {
		return errors.New("cannot update role name if the pod identity association was not created with a role name")
	}
	rs := builder.NewIAMRoleResourceSetForPodIdentity(&updateConfig.podIdentityAssociation)
	if err := rs.AddAllResources(); err != nil {
		return fmt.Errorf("adding resources to CloudFormation template: %w", err)
	}
	template, err := rs.RenderJSON()
	if err != nil {
		return fmt.Errorf("generating CloudFormation template: %w", err)
	}
	if err := u.StackUpdater.MustUpdateStack(ctx, manager.UpdateStackOptions{
		StackName:     updateConfig.stackName,
		ChangeSetName: fmt.Sprintf("eksctl-%s-%s-update-%d", updateConfig.podIdentityAssociation.Namespace, updateConfig.podIdentityAssociation.ServiceAccountName, time.Now().Unix()),
		Description:   fmt.Sprintf("updating IAM resources stack %q for pod identity association %q", updateConfig.stackName, podIdentityAssociationID),
		TemplateData:  manager.TemplateBody(template),
		Wait:          true,
	}); err != nil {
		if _, ok := err.(*manager.NoChangeError); ok {
			logger.Info("IAM resources for %q are already up-to-date", podIdentityAssociationID)
			return nil
		}
		return fmt.Errorf("updating IAM resources for pod identity association: %w", err)
	}
	logger.Info("updated IAM resources stack %q for %q", updateConfig.stackName, podIdentityAssociationID)
	stack, err = u.StackUpdater.DescribeStack(ctx, &manager.Stack{
		StackName: aws.String(updateConfig.stackName),
	})
	if err != nil {
		return fmt.Errorf("describing IAM resources stack: %w", err)
	}
	if err := rs.GetAllOutputs(*stack); err != nil {
		return fmt.Errorf("error getting IAM role output from IAM resources stack: %w", err)
	}

	return u.updatePodIdentityAssociation(ctx, updateConfig, podIdentityAssociationID)
}

func (u *Updater) updatePodIdentityAssociation(ctx context.Context, updateConfig *updateConfig, podIdentityAssociationID string) error {
	roleARN := updateConfig.podIdentityAssociation.RoleARN
	if _, err := u.APIUpdater.UpdatePodIdentityAssociation(ctx, &eks.UpdatePodIdentityAssociationInput{
		AssociationId: aws.String(updateConfig.associationID),
		ClusterName:   aws.String(u.ClusterName),
		RoleArn:       aws.String(roleARN),
	}); err != nil {
		return fmt.Errorf("updating pod identity association (associationID: %s, roleARN: %s): %w", updateConfig.associationID, roleARN, err)
	}
	logger.Info("updated role ARN %q for pod identity association %q", roleARN, podIdentityAssociationID)
	return nil
}

func (u *Updater) makeUpdate(ctx context.Context, p api.PodIdentityAssociation, roleStackNames []string) (*updateConfig, error) {
	const notFoundErrMsg = "pod identity association does not exist"
	output, err := u.APIUpdater.ListPodIdentityAssociations(ctx, &eks.ListPodIdentityAssociationsInput{
		ClusterName:    aws.String(u.ClusterName),
		Namespace:      aws.String(p.Namespace),
		ServiceAccount: aws.String(p.ServiceAccountName),
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
		describeOutput, err := u.APIUpdater.DescribePodIdentityAssociation(ctx, &eks.DescribePodIdentityAssociationInput{
			ClusterName:   aws.String(u.ClusterName),
			AssociationId: output.Associations[0].AssociationId,
		})
		if err != nil {
			return nil, fmt.Errorf("error describing pod identity association: %w", err)
		}
		stackName, hasStack := getIAMResourcesStack(roleStackNames, Identifier{
			Namespace:          p.Namespace,
			ServiceAccountName: p.ServiceAccountName,
		})
		if hasStack {
			if describeOutput.Association.RoleArn != nil && p.RoleARN != "" && p.RoleARN != *describeOutput.Association.RoleArn {
				return nil, errors.New("cannot change podIdentityAssociation.roleARN since the role was created by eksctl")
			}
		} else {
			if p.RoleARN == "" {
				return nil, errors.New("podIdentityAssociation.roleARN is required since the role was not created by eksctl")
			}
			podIDWithRoleARN := api.PodIdentityAssociation{
				Namespace:          p.Namespace,
				ServiceAccountName: p.ServiceAccountName,
				RoleARN:            p.RoleARN,
			}
			if !reflect.DeepEqual(p, podIDWithRoleARN) {
				return nil, errors.New("only namespace, serviceAccountName and roleARN can be specified if the role was not created by eksctl")
			}
		}
		return &updateConfig{
			podIdentityAssociation: p,
			associationID:          *describeOutput.Association.AssociationId,
			hasIAMResourcesStack:   hasStack,
			stackName:              stackName,
		}, nil
	}
}
