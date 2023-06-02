package fargate

import (
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/utils/retry"
)

// Client wraps around an EKS API client to expose high-level methods.
type Client struct {
	clusterName  string
	api          awsapi.EKS
	retryPolicy  retry.Policy
	stackManager manager.StackManager
}

func NewFromProvider(clusterName string, provider api.ClusterProvider, stackManager manager.StackManager) Client {
	retry := retry.NewTimingOutExponentialBackoff(provider.WaitTimeout())
	return NewWithRetryPolicy(clusterName, provider.EKS(), &retry, stackManager)
}

// NewWithRetryPolicy returns a new Fargate client configured with the
// provided retry policy for blocking/waiting operations.
func NewWithRetryPolicy(clusterName string, api awsapi.EKS, retryPolicy retry.Policy, stackManager manager.StackManager) Client {
	return Client{
		clusterName:  clusterName,
		api:          api,
		retryPolicy:  retryPolicy,
		stackManager: stackManager,
	}
}

func describeRequest(clusterName, profileName string) *eks.DescribeFargateProfileInput {
	request := &eks.DescribeFargateProfileInput{
		ClusterName:        &clusterName,
		FargateProfileName: &profileName,
	}
	logger.Debug("Fargate profile: describe request: sending: %#v", request)
	return request
}

func toSelectors(in []ekstypes.FargateProfileSelector) []api.FargateProfileSelector {
	out := make([]api.FargateProfileSelector, len(in))
	for i, selector := range in {
		out[i] = api.FargateProfileSelector{
			Namespace: *selector.Namespace,
			Labels:    selector.Labels,
		}
	}
	return out
}
