package podidentityassociation

import (
	"context"
	"fmt"
	"strings"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
)

// IAMRoleCreator creates IAM resources for a pod identity association.
type IAMRoleCreator struct {
	ClusterName  string
	StackCreator StackCreator
}

// Create creates IAM resources for podIdentityAssociation. If podIdentityAssociation belongs to an addon, addonName
// must be non-empty.
func (r *IAMRoleCreator) Create(ctx context.Context, podIdentityAssociation *api.PodIdentityAssociation, addonName string) (string, error) {
	rs := builder.NewIAMRoleResourceSetForPodIdentity(podIdentityAssociation)
	if err := rs.AddAllResources(); err != nil {
		return "", err
	}
	if podIdentityAssociation.Tags == nil {
		podIdentityAssociation.Tags = make(map[string]string)
	}
	podID := Identifier{
		Namespace:          podIdentityAssociation.Namespace,
		ServiceAccountName: podIdentityAssociation.ServiceAccountName,
	}.IDString()

	// If a target role ARN is specified for cross-account access, we need to:
	// 1. Add permission to assume the target role
	// 2. Configure the external ID for the target role trust relationship
	if podIdentityAssociation.TargetRoleARN != "" {
		// Extract account ID and role name from the target role ARN
		// ARN format: arn:aws:iam::ACCOUNT_ID:role/ROLE_NAME
		targetRoleARNParts := strings.Split(podIdentityAssociation.TargetRoleARN, ":")
		if len(targetRoleARNParts) >= 5 {
			targetAccountID := targetRoleARNParts[4]
			targetRoleName := strings.TrimPrefix(targetRoleARNParts[5], "role/")
			
			// Add permission to assume the target role
			// This will be added to the role's permission policy
			rs.AddAssumeRolePermission(podIdentityAssociation.TargetRoleARN)
			
			// Add a tag to indicate this role is configured for cross-account access
			podIdentityAssociation.Tags["eksctl.io/cross-account-role"] = "true"
			podIdentityAssociation.Tags["eksctl.io/target-account-id"] = targetAccountID
			podIdentityAssociation.Tags["eksctl.io/target-role-name"] = targetRoleName
		}
	}

	var stackName string
	if addonName != "" {
		podIdentityAssociation.Tags[api.AddonNameTag] = addonName
		podIdentityAssociation.Tags[api.AddonPodIdentityAssociationNameTag] = podID
		stackName = MakeAddonPodIdentityStackName(r.ClusterName, addonName, podIdentityAssociation.ServiceAccountName)
	} else {
		podIdentityAssociation.Tags[api.PodIdentityAssociationNameTag] = podID
		stackName = MakeStackName(r.ClusterName, podIdentityAssociation.Namespace, podIdentityAssociation.ServiceAccountName)
	}

	stackCh := make(chan error)
	if err := r.StackCreator.CreateStack(ctx, stackName, rs, podIdentityAssociation.Tags, nil, stackCh); err != nil {
		return "", fmt.Errorf("creating IAM role for pod identity association for service account %s in namespace %s: %w",
			podIdentityAssociation.ServiceAccountName, podIdentityAssociation.Namespace, err)
	}
	select {
	case err := <-stackCh:
		if err != nil {
			return "", err
		}
		return podIdentityAssociation.RoleARN, nil
	case <-ctx.Done():
		return "", fmt.Errorf("timed out waiting for creation of IAM resources for pod identity association %s: %w",
			podIdentityAssociation.NameString(), ctx.Err())
	}
}

// MakeStackName creates a stack name for the specified access entry.
func MakeStackName(clusterName, namespace, serviceAccountName string) string {
	return fmt.Sprintf("eksctl-%s-podidentityrole-%s-%s", clusterName, namespace, serviceAccountName)
}

func MakeAddonPodIdentityStackName(clusterName, addonName, serviceAccountName string) string {
	return fmt.Sprintf("eksctl-%s-addon-%s-podidentityrole-%s", clusterName, addonName, serviceAccountName)
}
