package podidentityassociation

import (
	"context"
	"fmt"

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
