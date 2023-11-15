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
		summaries      []Summary
		associationIDs []*string
	)

	output, err := g.eksAPI.ListPodIdentityAssociations(ctx, &awseks.ListPodIdentityAssociationsInput{
		ClusterName: &g.clusterName,
	})
	if err != nil {
		return summaries, fmt.Errorf("failed to list pod identity associations: %w", err)
	}

	for _, a := range output.Associations {
		associationIDs = append(associationIDs, a.AssociationId)
	}

	for _, id := range associationIDs {
		output, err := g.eksAPI.DescribePodIdentityAssociation(ctx, &awseks.DescribePodIdentityAssociationInput{
			ClusterName:   &g.clusterName,
			AssociationId: id,
		})
		if err != nil {
			return summaries, fmt.Errorf("failed to describe pod identity association with associationID: %s", *id)
		}

		if !shouldFetchPodIdentityAssociation(output, namespace, serviceAccountName) {
			continue
		}

		summaries = append(summaries, Summary{
			AssociationARN:     *output.Association.AssociationArn,
			Namespace:          *output.Association.Namespace,
			ServiceAccountName: *output.Association.ServiceAccount,
			RoleARN:            *output.Association.RoleArn,
		})
	}

	return summaries, nil
}

func shouldFetchPodIdentityAssociation(output *awseks.DescribePodIdentityAssociationOutput, namespace, serviceAccount string) bool {
	if namespace == "" ||
		(*output.Association.Namespace == namespace && serviceAccount == "") ||
		(*output.Association.Namespace == namespace && *output.Association.ServiceAccount == serviceAccount) {
		return true
	}
	return false
}
