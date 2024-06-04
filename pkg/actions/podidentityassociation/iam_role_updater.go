package podidentityassociation

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	"time"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"golang.org/x/exp/slices"
)

// IAMRoleUpdater updates IAM resources for pod identity associations.
type IAMRoleUpdater struct {
	// StackUpdater updates CloudFormation stacks.
	StackUpdater StackUpdater
}

// Update updates IAM resources for podIdentityAssociation and returns an IAM role ARN upon success. The boolean return value reports
// whether the IAM resources have changed or not.
func (u *IAMRoleUpdater) Update(ctx context.Context, podIdentityAssociation api.PodIdentityAssociation, stackName, podIdentityAssociationID string) (string, bool, error) {
	stack, err := u.StackUpdater.DescribeStack(ctx, &manager.Stack{
		StackName: aws.String(stackName),
	})
	if err != nil {
		return "", false, fmt.Errorf("describing IAM resources stack %q: %w", stackName, err)
	}
	if podIdentityAssociation.RoleName != "" && !slices.Contains(stack.Capabilities, cfntypes.CapabilityCapabilityNamedIam) {
		return "", false, errors.New("cannot update role name if the pod identity association was not created with a role name")
	}
	rs := builder.NewIAMRoleResourceSetForPodIdentity(&podIdentityAssociation)
	if err := rs.AddAllResources(); err != nil {
		return "", false, fmt.Errorf("adding resources to CloudFormation template: %w", err)
	}
	template, err := rs.RenderJSON()
	if err != nil {
		return "", false, fmt.Errorf("generating CloudFormation template: %w", err)
	}
	if err := u.StackUpdater.MustUpdateStack(ctx, manager.UpdateStackOptions{
		StackName:     stackName,
		ChangeSetName: fmt.Sprintf("eksctl-%s-%s-update-%d", podIdentityAssociation.Namespace, podIdentityAssociation.ServiceAccountName, time.Now().Unix()),
		Description:   fmt.Sprintf("updating IAM resources stack %q for pod identity association %q", stackName, podIdentityAssociationID),
		TemplateData:  manager.TemplateBody(template),
		Wait:          true,
	}); err != nil {
		var noChangeErr *manager.NoChangeError
		if errors.As(err, &noChangeErr) {
			logger.Info("IAM resources for %s (pod identity association ID: %s) are already up-to-date", podIdentityAssociation.NameString(), podIdentityAssociationID)
			return podIdentityAssociation.RoleARN, false, nil
		}
		return "", false, fmt.Errorf("updating IAM resources for pod identity association: %w", err)
	}
	logger.Info("updated IAM resources stack %q for %q", stackName, podIdentityAssociationID)

	stack, err = u.StackUpdater.DescribeStack(ctx, &manager.Stack{
		StackName: aws.String(stackName),
	})
	if err != nil {
		return "", false, fmt.Errorf("describing IAM resources stack: %w", err)
	}

	if err := populateRoleARN(rs, stack); err != nil {
		return "", false, err
	}
	return podIdentityAssociation.RoleARN, true, nil
}

func populateRoleARN(rs builder.ResourceSet, stack *manager.Stack) error {
	if err := rs.GetAllOutputs(*stack); err != nil {
		return fmt.Errorf("error getting IAM role output from IAM resources stack: %w", err)
	}
	return nil
}
