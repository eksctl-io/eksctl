package fargate

import (
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/kris-nova/logger"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/utils/retry"
	"github.com/weaveworks/eksctl/pkg/utils/strings"
)

func NewFromProvider(clusterName string, provider api.ClusterProvider) Client {
	retry := retry.NewTimingOutExponentialBackoff(provider.WaitTimeout())
	return NewWithRetryPolicy(clusterName, provider.EKS(), &retry)
}

// NewWithRetryPolicy returns a new Fargate client configured with the
// provided retry policy for blocking/waiting operations.
func NewWithRetryPolicy(clusterName string, api eksiface.EKSAPI, retryPolicy retry.Policy) Client {
	return Client{
		clusterName: clusterName,
		api:         api,
		retryPolicy: retryPolicy,
	}
}

// Client wraps around an EKS API client to expose high-level methods.
type Client struct {
	clusterName string
	api         eksiface.EKSAPI
	retryPolicy retry.Policy
}

func describeRequest(clusterName string, profileName string) *eks.DescribeFargateProfileInput {
	request := &eks.DescribeFargateProfileInput{
		ClusterName:        &clusterName,
		FargateProfileName: &profileName,
	}
	logger.Debug("Fargate profile: describe request: sending: %#v", request)
	return request
}

func toSelectors(in []*eks.FargateProfileSelector) []api.FargateProfileSelector {
	out := make([]api.FargateProfileSelector, len(in))
	for i, selector := range in {
		out[i] = api.FargateProfileSelector{
			Namespace: *selector.Namespace,
			Labels:    strings.ToValuesMap(selector.Labels),
		}
	}
	return out
}
