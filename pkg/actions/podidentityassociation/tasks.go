package podidentityassociation

import (
	"context"
	"fmt"

	awseks "github.com/aws/aws-sdk-go-v2/service/eks"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_stack_manager.go . StackManager
type StackManager interface {
	CreateStack(ctx context.Context, name string, stack builder.ResourceSetReader, tags, parameters map[string]string, errs chan error) error
}

type createIAMRoleTask struct {
	ctx                    context.Context
	info                   string
	clusterName            string
	podIdentityAssociation *api.PodIdentityAssociation
	stackManager           StackManager
}

func (t *createIAMRoleTask) Describe() string {
	return t.info
}

func (t *createIAMRoleTask) Do(errorCh chan error) error {
	rs := builder.NewIAMRoleResourceSetForPodIdentity(t.podIdentityAssociation)
	if err := rs.AddAllResources(); err != nil {
		return err
	}
	if err := t.stackManager.CreateStack(t.ctx,
		MakeStackName(
			t.clusterName,
			t.podIdentityAssociation.Namespace,
			t.podIdentityAssociation.ServiceAccountName),
		rs, nil, nil, errorCh); err != nil {
		return fmt.Errorf("creating IAM role for pod identity association for service account %s in namespace %s: %w",
			t.podIdentityAssociation.ServiceAccountName, t.podIdentityAssociation.Namespace, err)
	}
	return nil
}

type createPodIdentityAssociationTask struct {
	ctx                    context.Context
	info                   string
	clusterName            string
	podIdentityAssociation *api.PodIdentityAssociation
	eksAPI                 awsapi.EKS
}

func (t *createPodIdentityAssociationTask) Describe() string {
	return t.info
}

func (t *createPodIdentityAssociationTask) Do(errorCh chan error) error {
	defer close(errorCh)

	if _, err := t.eksAPI.CreatePodIdentityAssociation(t.ctx, &awseks.CreatePodIdentityAssociationInput{
		ClusterName:    &t.clusterName,
		Namespace:      &t.podIdentityAssociation.Namespace,
		RoleArn:        &t.podIdentityAssociation.RoleARN,
		ServiceAccount: &t.podIdentityAssociation.ServiceAccountName,
		Tags:           t.podIdentityAssociation.Tags,
	}); err != nil {
		return fmt.Errorf(
			"creating pod identity association for service account %s in namespace %s: %w",
			t.podIdentityAssociation.ServiceAccountName, t.podIdentityAssociation.Namespace, err)
	}

	return nil
}

func makeStackNamePrefix(clusterName string) string {
	return fmt.Sprintf("eksctl-%s-podidentityrole-ns-", clusterName)
}

// MakeStackName creates a stack name for the specified access entry.
func MakeStackName(clusterName, namespace, serviceAccountName string) string {
	return fmt.Sprintf("%s%s-sa-%s", makeStackNamePrefix(clusterName), namespace, serviceAccountName)
}
