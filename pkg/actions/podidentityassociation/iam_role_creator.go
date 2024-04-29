package podidentityassociation

import (
	"context"
	"fmt"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
)

type IAMRoleCreator struct {
	ClusterName  string
	StackCreator StackCreator
}

func (r *IAMRoleCreator) Create(ctx context.Context, podIdentityAssociation *api.PodIdentityAssociation) (string, error) {
	rs := builder.NewIAMRoleResourceSetForPodIdentity(podIdentityAssociation)
	if err := rs.AddAllResources(); err != nil {
		return "", err
	}
	if podIdentityAssociation.Tags == nil {
		podIdentityAssociation.Tags = make(map[string]string)
	}
	podIdentityAssociation.Tags[api.PodIdentityAssociationNameTag] = Identifier{
		Namespace:          podIdentityAssociation.Namespace,
		ServiceAccountName: podIdentityAssociation.ServiceAccountName,
	}.IDString()

	stackName := MakeStackName(r.ClusterName, podIdentityAssociation.Namespace, podIdentityAssociation.ServiceAccountName)
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
