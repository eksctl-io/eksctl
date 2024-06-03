package podidentityassociation

import (
	"context"
	"fmt"

	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/weaveworks/eksctl/pkg/awsapi"
)

type Summary struct {
	AssociationARN     string
	Namespace          string
	ServiceAccountName string
	RoleARN            string
	OwnerARN           string
}

type Getter struct {
	clusterName string
	eksAPI      awsapi.EKS
}

func NewGetter(clusterName string, eksAPI awsapi.EKS) *Getter {
	return &Getter{
		clusterName: clusterName,
		eksAPI:      eksAPI,
	}
}

func (g *Getter) GetPodIdentityAssociations(ctx context.Context, namespace, serviceAccountName string) ([]Summary, error) {
	var (
		summaries []Summary
	)

	input := &awseks.ListPodIdentityAssociationsInput{ClusterName: &g.clusterName}
	if namespace != "" {
		input.Namespace = &namespace
	}
	if serviceAccountName != "" {
		input.ServiceAccount = &serviceAccountName
	}

	listOut, err := g.eksAPI.ListPodIdentityAssociations(ctx, input)
	if err != nil {
		return summaries, fmt.Errorf("failed to list pod identity associations: %w", err)
	}

	for _, a := range listOut.Associations {
		describeOut, err := g.eksAPI.DescribePodIdentityAssociation(ctx, &awseks.DescribePodIdentityAssociationInput{
			ClusterName:   &g.clusterName,
			AssociationId: a.AssociationId,
		})
		if err != nil {
			return summaries, fmt.Errorf("failed to describe pod identity association with associationID: %s", *a.AssociationId)
		}

		summary := Summary{
			AssociationARN:     *describeOut.Association.AssociationArn,
			Namespace:          *describeOut.Association.Namespace,
			ServiceAccountName: *describeOut.Association.ServiceAccount,
			RoleARN:            *describeOut.Association.RoleArn,
		}
		if describeOut.Association.OwnerArn != nil {
			summary.OwnerARN = *describeOut.Association.OwnerArn
		}
		summaries = append(summaries, summary)
	}

	return summaries, nil
}
