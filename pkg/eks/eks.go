package eks

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/version"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

// DescribeControlPlane describes the cluster control plane
func (c *ClusterProvider) DescribeControlPlane(ctx context.Context, meta *api.ClusterMeta) (*ekstypes.Cluster, error) {
	input := &awseks.DescribeClusterInput{
		Name: &meta.Name,
	}
	output, err := c.AWSProvider.EKS().DescribeCluster(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "unable to describe cluster control plane")
	}
	return output.Cluster, nil
}

// RefreshClusterStatus calls c.DescribeControlPlane and caches the results;
// it parses the credentials (endpoint, CA certificate) and stores them in ClusterConfig.Status,
// so that a Kubernetes client can be constructed; additionally it caches Kubernetes
// version (use ctl.ControlPlaneVersion to retrieve it) and other properties in
// c.Status.cachedClusterInfo.
// It also updates ClusterConfig to reflect the current cluster state.
func (c *ClusterProvider) RefreshClusterStatus(ctx context.Context, spec *api.ClusterConfig) error {
	cluster, err := c.DescribeControlPlane(ctx, spec.Metadata)
	if err != nil {
		return err
	}
	logger.Debug("cluster = %#v", cluster)

	if isNonEKSCluster(cluster) {
		return errors.Errorf("cannot perform this operation on a non-EKS cluster; please follow the documentation for "+
			"cluster %s's Kubernetes provider", spec.Metadata.Name)
	}

	if spec.Status == nil {
		spec.Status = &api.ClusterStatus{}
	}

	c.Status.ClusterInfo = &ClusterInfo{
		Cluster: cluster,
	}

	switch cluster.Status {
	case ekstypes.ClusterStatusCreating, ekstypes.ClusterStatusDeleting, ekstypes.ClusterStatusFailed:
		return nil
	default:
		return spec.SetClusterState(cluster)
	}
}

// isNonEKSCluster returns true if the cluster is external
func isNonEKSCluster(cluster *ekstypes.Cluster) bool {
	return cluster.ConnectorConfig != nil
}

var (
	platformVersionRegex = regexp.MustCompile(`^eks\.(\d+)$`)
)

// PlatformVersion extracts the digit X in the provided platform version eks.X
func PlatformVersion(platformVersion string) (int, error) {
	match := platformVersionRegex.FindStringSubmatch(platformVersion)
	if len(match) != 2 {
		return -1, fmt.Errorf("failed to parse cluster's platform version: %q", platformVersion)
	}
	versionStr := match[1]
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return -1, err
	}
	return version, nil
}

// RefreshClusterStatusIfStale refreshes the cluster status if enough time has passed since the last refresh
func (c *ClusterProvider) RefreshClusterStatusIfStale(ctx context.Context, spec *api.ClusterConfig) error {
	if c.Status.ClusterInfo == nil {
		return c.RefreshClusterStatus(ctx, spec)
	}
	return nil
}

// CanOperate returns true when a cluster can be operated, otherwise it returns false along with an error explaining the reason
func (c *ClusterProvider) CanOperate(spec *api.ClusterConfig) (bool, error) {
	// if the check before calling this failed, it won't have a clusterInfo meaning,
	// we either ignored this error during delete, or the Refresh failed anyway. In both cases the cluster is NOT operable.
	if c.Status.ClusterInfo == nil {
		return false, fmt.Errorf("cluster info not available")
	}
	switch status := c.Status.ClusterInfo.Cluster.Status; status {
	case ekstypes.ClusterStatusCreating, ekstypes.ClusterStatusDeleting, ekstypes.ClusterStatusFailed:
		return false, fmt.Errorf("cannot perform Kubernetes API operations on cluster %q in %q region due to status %q", spec.Metadata.Name, spec.Metadata.Region, status)
	default:
		return true, nil
	}
}

// CanOperateWithRefresh returns true when a cluster can be operated, otherwise it returns false along with an error explaining the reason
func (c *ClusterProvider) CanOperateWithRefresh(ctx context.Context, spec *api.ClusterConfig) (bool, error) {
	if err := c.RefreshClusterStatusIfStale(ctx, spec); err != nil {
		return false, errors.Wrapf(err, "unable to fetch cluster status to determine operability")
	}

	switch status := c.Status.ClusterInfo.Cluster.Status; status {
	case ekstypes.ClusterStatusCreating, ekstypes.ClusterStatusDeleting, ekstypes.ClusterStatusFailed:
		return false, fmt.Errorf("cannot perform Kubernetes API operations on cluster %q in %q region due to status %q", spec.Metadata.Name, spec.Metadata.Region, status)
	default:
		return true, nil
	}
}

// CanUpdate return true when a cluster or add-ons can be updated, otherwise it returns false along with an error explaining the reason
func (c *ClusterProvider) CanUpdate(spec *api.ClusterConfig) (bool, error) {
	if c.Status.ClusterInfo == nil {
		return false, nil
	}
	switch status := c.Status.ClusterInfo.Cluster.Status; status {
	case ekstypes.ClusterStatusActive:
		// only active cluster can be upgraded
		return true, nil
	default:
		return false, fmt.Errorf("cannot update cluster %q in %q region due to status %q", spec.Metadata.Name, spec.Metadata.Region, status)
	}
}

// ControlPlaneVersion returns cached version (EKS API)
func (c *ClusterProvider) ControlPlaneVersion() string {
	if c.Status.ClusterInfo == nil || c.Status.ClusterInfo.Cluster == nil || c.Status.ClusterInfo.Cluster.Version == nil {
		return ""
	}
	return *c.Status.ClusterInfo.Cluster.Version
}

// ControlPlaneVPCInfo returns cached version (EKS API)
func (c *ClusterProvider) ControlPlaneVPCInfo() ekstypes.VpcConfigResponse {
	if c.Status.ClusterInfo == nil || c.Status.ClusterInfo.Cluster == nil || c.Status.ClusterInfo.Cluster.ResourcesVpcConfig == nil {
		return ekstypes.VpcConfigResponse{}
	}
	return *c.Status.ClusterInfo.Cluster.ResourcesVpcConfig
}

// NewOpenIDConnectManager returns OpenIDConnectManager
func (c *ClusterProvider) NewOpenIDConnectManager(ctx context.Context, spec *api.ClusterConfig) (*iamoidc.OpenIDConnectManager, error) {
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
	if !api.Partitions.IsSupported(parsedARN.Partition) {
		return nil, fmt.Errorf("unknown EKS ARN: %q", spec.Status.ARN)
	}

	return iamoidc.NewOpenIDConnectManager(c.AWSProvider.IAM(), parsedARN.AccountID,
		*c.Status.ClusterInfo.Cluster.Identity.Oidc.Issuer, parsedARN.Partition, sharedTags(c.Status.ClusterInfo.Cluster))
}

func sharedTags(cluster *ekstypes.Cluster) map[string]string {
	return map[string]string{
		api.ClusterNameTag:   *cluster.Name,
		api.EksctlVersionTag: version.GetVersion(),
	}

}

// LoadClusterVPC loads the VPC configuration.
func (c *ClusterProvider) LoadClusterVPC(ctx context.Context, spec *api.ClusterConfig, stack *manager.Stack, ignoreDrift bool) error {
	return vpc.UseFromClusterStack(ctx, c.AWSProvider, stack, spec, ignoreDrift)
}

// GetCluster display details of an EKS cluster in your account
func (c *ClusterProvider) GetCluster(ctx context.Context, clusterName string) (*ekstypes.Cluster, error) {
	input := &awseks.DescribeClusterInput{
		Name: &clusterName,
	}

	output, err := c.AWSProvider.EKS().DescribeCluster(ctx, input)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to describe control plane %q", clusterName)
	}
	logger.Debug("cluster = %#v", output)

	if output.Cluster.Status == ekstypes.ClusterStatusActive {
		if logger.Level >= 4 {
			spec := &api.ClusterConfig{Metadata: &api.ClusterMeta{Name: clusterName}}
			stacks, err := c.NewStackManager(spec).ListStacksWithStatuses(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "listing CloudFormation stack for %q", clusterName)
			}
			for _, s := range stacks {
				logger.Debug("stack = %#v", *s)
			}
		}
	}
	return output.Cluster, nil
}

// GetClusterState returns the EKS cluster state.
func (c *ClusterProvider) GetClusterState() *ekstypes.Cluster {
	return c.Status.ClusterInfo.Cluster
}

// IsAccessEntryEnabled reports whether the cluster has access entries enabled.
func (c *ClusterProvider) IsAccessEntryEnabled() bool {
	return IsAccessEntryEnabled(c.Status.ClusterInfo.Cluster.AccessConfig)
}

// IsAccessEntryEnabled reports whether the specified accessConfig has access entries enabled.
func IsAccessEntryEnabled(accessConfig *ekstypes.AccessConfigResponse) bool {
	return accessConfig != nil && accessConfig.AuthenticationMode != ekstypes.AuthenticationModeConfigMap
}
