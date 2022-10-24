import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/version"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

func (c *ClusterManager) NewOpenIDConnectManager(spec *api.ClusterConfig) (*iamoidc.OpenIDConnectManager, error) {
	if _, err := c.CanOperateWithRefresh(ctx, spec); err != nil {
		return nil, err
	}

	if c.Status.ClusterInfo.Cluster == nil || c.Status.ClusterInfo.Cluster.Identity == nil || c.Status.ClusterInfo.Cluster.Identity.Oidc == nil || c.Status.ClusterInfo.Cluster.Identity.Oidc.Issuer == nil {
		return nil, &iamoidc.UnsupportedOIDCError{Message: "unknown OIDC issuer URL"}
	}

	parsedARN, err := arn.Parse(spec.Status.ARN)
	if err != nil {
		return nil, errors.Wrapf(err, "unexpected invalid ARN: %q", spec.Status.ARN)
	}
	switch parsedARN.Partition {
	case "aws", "aws-cn", "aws-us-gov":
	default:
		return nil, fmt.Errorf("unknown EKS ARN: %q", spec.Status.ARN)
	}

	return iamoidc.NewOpenIDConnectManager(c.AWSProvider.IAM(), parsedARN.AccountID,
		*c.Status.ClusterInfo.Cluster.Identity.Oidc.Issuer, parsedARN.Partition, sharedTags(c.Status.ClusterInfo.Cluster))
}

func (c *ClusterManager) LoadClusterVPC(spec *api.ClusterConfig, stack *manager.Stack) error {
	return vpc.UseFromClusterStack(ctx, c.AWSProvider, stack, spec)
}