package accessentry

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/eks"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
)

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

type Summary struct {
	PrincipalARN     string             `json:"principalARN"`
	KubernetesGroups []string           `json:"kubernetesGroups,omitempty"`
	AccessPolicies   []api.AccessPolicy `json:"accessPolicies,omitempty"`
}

func (aeg *Getter) Get(ctx context.Context, principalARN api.ARN) ([]Summary, error) {

	toBeFetched := []string{principalARN.String()}
	// if no principal ARN was specified, we fetch all entries for the cluster
	if principalARN.IsZero() {
		out, err := aeg.eksAPI.ListAccessEntries(ctx, &eks.ListAccessEntriesInput{
			ClusterName: &aeg.clusterName,
		})
		if err != nil {
			return nil, fmt.Errorf("calling EKS API to list access entries: %w", err)
		}
		toBeFetched = out.AccessEntries
	}

	var summaries []Summary
	for _, pARN := range toBeFetched {
		summary, err := aeg.getIndividualEntry(ctx, pARN)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}
	return summaries, nil
}

func (aeg *Getter) getIndividualEntry(ctx context.Context, principalARN string) (Summary, error) {
	summary := Summary{
		PrincipalARN:   principalARN,
		AccessPolicies: []api.AccessPolicy{},
	}

	// fetch kubernetes groups
	entry, err := aeg.eksAPI.DescribeAccessEntry(ctx, &eks.DescribeAccessEntryInput{
		ClusterName:  &aeg.clusterName,
		PrincipalArn: &principalARN,
	})
	if err != nil {
		return Summary{}, fmt.Errorf("calling EKS API to describe access entry with principal ARN %s: %w", principalARN, err)
	}
	summary.KubernetesGroups = entry.AccessEntry.KubernetesGroups

	// fetch associated polices
	policies, err := aeg.eksAPI.ListAssociatedAccessPolicies(ctx, &eks.ListAssociatedAccessPoliciesInput{
		ClusterName:  &aeg.clusterName,
		PrincipalArn: &principalARN,
	})
	if err != nil {
		return Summary{}, fmt.Errorf("calling EKS API to list associated access policies for entry with principal ARN %s: %w", principalARN, err)
	}
	for _, policy := range policies.AssociatedAccessPolicies {
		p := api.AccessPolicy{
			PolicyARN: api.MustParseARN(*policy.PolicyArn),
			AccessScope: api.AccessScope{
				Type:       policy.AccessScope.Type,
				Namespaces: policy.AccessScope.Namespaces,
			},
		}
		summary.AccessPolicies = append(summary.AccessPolicies, p)
	}

	return summary, nil
}
