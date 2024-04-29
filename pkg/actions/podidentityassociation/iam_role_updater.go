package podidentityassociation

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"golang.org/x/exp/slices"
	"time"
)

// IAMRoleUpdater updates IAM resources for pod identity associations.
type IAMRoleUpdater struct {
	// StackUpdater updates CloudFormation stacks.
	StackUpdater StackUpdater
}

// Update updates IAM resources for updateConfig and returns an IAM role ARN upon success. The boolean return value reports
// whether the IAM resources have changed or not.
func (u *IAMRoleUpdater) Update(ctx context.Context, updateConfig *UpdateConfig, podIdentityAssociationID string) (string, bool, error) {
	stack, err := u.StackUpdater.DescribeStack(ctx, &manager.Stack{
		StackName: aws.String(updateConfig.StackName),
	})
	if err != nil {
		return "", false, fmt.Errorf("describing IAM resources stack %q: %w", updateConfig.StackName, err)
	}
	if updateConfig.PodIdentityAssociation.RoleName != "" && !slices.Contains(stack.Capabilities, cfntypes.CapabilityCapabilityNamedIam) {
		return "", false, errors.New("cannot update role name if the pod identity association was not created with a role name")
	}
	rs := builder.NewIAMRoleResourceSetForPodIdentity(&updateConfig.PodIdentityAssociation)
	if err := rs.AddAllResources(); err != nil {
		return "", false, fmt.Errorf("adding resources to CloudFormation template: %w", err)
	}
	template, err := rs.RenderJSON()
	if err != nil {
		return "", false, fmt.Errorf("generating CloudFormation template: %w", err)
	}
	if err := u.StackUpdater.MustUpdateStack(ctx, manager.UpdateStackOptions{
		StackName:     updateConfig.StackName,
		ChangeSetName: fmt.Sprintf("eksctl-%s-%s-update-%d", updateConfig.PodIdentityAssociation.Namespace, updateConfig.PodIdentityAssociation.ServiceAccountName, time.Now().Unix()),
		Description:   fmt.Sprintf("updating IAM resources stack %q for pod identity association %q", updateConfig.StackName, podIdentityAssociationID),
		TemplateData:  manager.TemplateBody(template),
		Wait:          true,
	}); err != nil {
		var noChangeErr *manager.NoChangeError
		if errors.As(err, &noChangeErr) {
			logger.Info("IAM resources for %q are already up-to-date", podIdentityAssociationID)
			return updateConfig.PodIdentityAssociation.RoleARN, false, nil
		}
		return "", false, fmt.Errorf("updating IAM resources for pod identity association: %w", err)
	}
	logger.Info("updated IAM resources stack %q for %q", updateConfig.StackName, podIdentityAssociationID)

	stack, err = u.StackUpdater.DescribeStack(ctx, &manager.Stack{
		StackName: aws.String(updateConfig.StackName),
	})
	if err != nil {
		return "", false, fmt.Errorf("describing IAM resources stack: %w", err)
	}
	if err := rs.GetAllOutputs(*stack); err != nil {
		return "", false, fmt.Errorf("error getting IAM role output from IAM resources stack: %w", err)
	}
	return updateConfig.PodIdentityAssociation.RoleARN, true, nil
}

func (u *IAMRoleUpdater) updateStack(ctx context.Context, updateConfig *UpdateConfig, podIdentityAssociationID string) error {
	stack, err := u.StackUpdater.DescribeStack(ctx, &manager.Stack{
		StackName: aws.String(updateConfig.StackName),
	})
	if err != nil {
		return fmt.Errorf("describing IAM resources stack %q: %w", updateConfig.StackName, err)
	}
	if updateConfig.PodIdentityAssociation.RoleName != "" && !slices.Contains(stack.Capabilities, cfntypes.CapabilityCapabilityNamedIam) {
		return errors.New("cannot update role name if the pod identity association was not created with a role name")
	}
	rs := builder.NewIAMRoleResourceSetForPodIdentity(&updateConfig.PodIdentityAssociation)
	if err := rs.AddAllResources(); err != nil {
		return fmt.Errorf("adding resources to CloudFormation template: %w", err)
	}
	template, err := rs.RenderJSON()
	if err != nil {
		return fmt.Errorf("generating CloudFormation template: %w", err)
	}
	if err := u.StackUpdater.MustUpdateStack(ctx, manager.UpdateStackOptions{
		StackName:     updateConfig.StackName,
		ChangeSetName: fmt.Sprintf("eksctl-%s-%s-update-%d", updateConfig.PodIdentityAssociation.Namespace, updateConfig.PodIdentityAssociation.ServiceAccountName, time.Now().Unix()),
		Description:   fmt.Sprintf("updating IAM resources stack %q for pod identity association %q", updateConfig.StackName, podIdentityAssociationID),
		TemplateData:  manager.TemplateBody(template),
		Wait:          true,
	}); err != nil {
		var noChangeErr *manager.NoChangeError
		if errors.As(err, &noChangeErr) {
			logger.Info("IAM resources for %q are already up-to-date", podIdentityAssociationID)
			return nil
		}
		return fmt.Errorf("updating IAM resources for pod identity association: %w", err)
	}
	logger.Info("updated IAM resources stack %q for %q", updateConfig.StackName, podIdentityAssociationID)

	stack, err = u.StackUpdater.DescribeStack(ctx, &manager.Stack{
		StackName: aws.String(updateConfig.StackName),
	})
	if err != nil {
		return fmt.Errorf("describing IAM resources stack: %w", err)
	}
	if err := rs.GetAllOutputs(*stack); err != nil {
		return fmt.Errorf("error getting IAM role output from IAM resources stack: %w", err)
	}
	return nil
}
