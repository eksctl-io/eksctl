package outposts

import (
	"context"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"

	"github.com/aws/aws-sdk-go-v2/service/eks"

	"github.com/weaveworks/eksctl/pkg/awsapi"
)

// WrapClusterProvider wraps ClusterProvider with a custom EKS service that overrides operations unsupported on Outposts
// with no-ops.
func WrapClusterProvider(p api.ClusterProvider) api.ClusterProvider {
	return &clusterProvider{
		ClusterProvider: p,
		eksService: &eksService{
			EKS: p.EKS(),
		},
	}
}

type clusterProvider struct {
	api.ClusterProvider
	eksService *eksService
}

func (c *clusterProvider) EKS() awsapi.EKS {
	return c.eksService
}

type eksService struct {
	awsapi.EKS
}

func (e *eksService) ListNodegroups(_ context.Context, _ *eks.ListNodegroupsInput, _ ...func(options *eks.Options)) (*eks.ListNodegroupsOutput, error) {
	return &eks.ListNodegroupsOutput{
		Nodegroups: []string{},
	}, nil
}
