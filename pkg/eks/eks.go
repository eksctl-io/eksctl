package eks

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/utils/waiters"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	awseks "github.com/aws/aws-sdk-go/service/eks"

	"k8s.io/client-go/kubernetes"

	"github.com/weaveworks/eksctl/pkg/actions/fargate"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/utils"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

const (
	eksctlCreatedTrue    api.EKSCTLCreated = "True"
	eksctlCreatedFalse   api.EKSCTLCreated = "False"
	eksctlCreatedUnknown api.EKSCTLCreated = "Unknown"
)

// DescribeControlPlane describes the cluster control plane
func (c *ClusterProvider) DescribeControlPlane(meta *api.ClusterMeta) (*awseks.Cluster, error) {
	input := &awseks.DescribeClusterInput{
		Name: &meta.Name,
	}
	output, err := c.Provider.EKS().DescribeCluster(input)
	if err != nil {
		return nil, errors.Wrap(err, "unable to describe cluster control plane")
	}
	return output.Cluster, nil
}

// RefreshClusterStatus calls c.DescribeControlPlane and caches the results;
// it parses the credentials (endpoint, CA certificate) and stores them in spec.Status,
// so that a Kubernetes client can be constructed; additionally it caches Kubernetes
// version (use ctl.ControlPlaneVersion to retrieve it) and other properties in
// c.Status.cachedClusterInfo
func (c *ClusterProvider) RefreshClusterStatus(spec *api.ClusterConfig) error {
	cluster, err := c.DescribeControlPlane(spec.Metadata)
	if err != nil {
		return err
	}
	logger.Debug("cluster = %#v", cluster)

	if spec.Status == nil {
		spec.Status = &api.ClusterStatus{}
	}

	c.setClusterInfo(cluster)

	switch *cluster.Status {
	case awseks.ClusterStatusCreating, awseks.ClusterStatusDeleting, awseks.ClusterStatusFailed:
		return nil
	default:
		data, err := base64.StdEncoding.DecodeString(*cluster.CertificateAuthority.Data)
		if err != nil {
			return errors.Wrap(err, "decoding certificate authority data")
		}
		spec.Status.Endpoint = *cluster.Endpoint
		spec.Status.CertificateAuthorityData = data
		spec.Status.ARN = *cluster.Arn
		return nil
	}
}

// SupportsManagedNodes reports whether an existing cluster supports Managed Nodes
// The minimum required control plane version and platform version are 1.14 and eks.3 respectively
func (c *ClusterProvider) SupportsManagedNodes(clusterConfig *api.ClusterConfig) (bool, error) {
	if err := c.maybeRefreshClusterStatus(clusterConfig); err != nil {
		return false, err
	}

	return ClusterSupportsManagedNodes(c.Status.ClusterInfo.Cluster)
}

// ClusterSupportsManagedNodes reports whether the EKS cluster supports managed nodes
func ClusterSupportsManagedNodes(cluster *awseks.Cluster) (bool, error) {
	supportsManagedNodes, err := utils.IsMinVersion(api.Version1_15, *cluster.Version)
	if err != nil {
		return false, err
	}
	if supportsManagedNodes {
		return true, nil
	}

	versionSupportsManagedNodes, err := VersionSupportsManagedNodes(*cluster.Version)
	if err != nil {
		return false, err
	}

	if !versionSupportsManagedNodes {
		return false, nil
	}

	if cluster.PlatformVersion == nil {
		logger.Warning("could not find cluster's platform version")
		return false, nil
	}
	version, err := PlatformVersion(*cluster.PlatformVersion)
	if err != nil {
		return false, err
	}
	minSupportedVersion := 3
	return version >= minSupportedVersion, nil
}

// SupportsFargate reports whether an existing cluster supports Fargate.
func (c *ClusterProvider) SupportsFargate(clusterConfig *api.ClusterConfig) (bool, error) {
	if err := c.maybeRefreshClusterStatus(clusterConfig); err != nil {
		return false, err
	}
	return ClusterSupportsFargate(c.Status.ClusterInfo.Cluster)
}

// ClusterSupportsFargate reports whether an existing cluster supports Fargate.
func ClusterSupportsFargate(cluster *awseks.Cluster) (bool, error) {
	supportsFargate, err := utils.IsMinVersion(api.Version1_15, *cluster.Version)
	if err != nil {
		return false, err
	}
	if supportsFargate {
		return true, nil
	}

	versionSupportsFargate, err := utils.IsMinVersion(fargate.MinKubernetesVersion, *cluster.Version)
	if err != nil {
		return false, err
	}
	if !versionSupportsFargate {
		return false, nil
	}

	if cluster.PlatformVersion == nil {
		logger.Warning("could not find cluster's platform version")
		return false, nil
	}
	version, err := PlatformVersion(*cluster.PlatformVersion)
	if err != nil {
		return false, err
	}
	return version >= fargate.MinPlatformVersion, nil
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

func (c *ClusterProvider) maybeRefreshClusterStatus(spec *api.ClusterConfig) error {
	if c.clusterInfoNeedsUpdate() {
		return c.RefreshClusterStatus(spec)
	}
	return nil
}

// CanDelete return true when a cluster can be deleted, otherwise it returns false along with an error explaining the reason
func (c *ClusterProvider) CanDelete(spec *api.ClusterConfig) (bool, error) {
	err := c.maybeRefreshClusterStatus(spec)
	if err != nil {
		if awsError, ok := errors.Unwrap(errors.Unwrap(err)).(awserr.Error); ok &&
			awsError.Code() == awseks.ErrCodeResourceNotFoundException {
			return true, nil
		}
		return false, errors.Wrapf(err, "fetching cluster status to determine if it can be deleted")
	}
	// it must be possible to delete cluster in any state
	return true, nil
}

// CanOperate return true when a cluster can be operated, otherwise it returns false along with an error explaining the reason
func (c *ClusterProvider) CanOperate(spec *api.ClusterConfig) (bool, error) {
	err := c.maybeRefreshClusterStatus(spec)
	if err != nil {
		return false, errors.Wrapf(err, "unable to fetch cluster status to determine operability")
	}

	switch status := *c.Status.ClusterInfo.Cluster.Status; status {
	case awseks.ClusterStatusCreating, awseks.ClusterStatusDeleting, awseks.ClusterStatusFailed:
		return false, fmt.Errorf("cannot perform Kubernetes API operations on cluster %q in %q region due to status %q", spec.Metadata.Name, spec.Metadata.Region, status)
	default:
		return true, nil
	}
}

// CanUpdate return true when a cluster or add-ons can be updated, otherwise it returns false along with an error explaining the reason
func (c *ClusterProvider) CanUpdate(spec *api.ClusterConfig) (bool, error) {
	err := c.maybeRefreshClusterStatus(spec)
	if err != nil {
		return false, errors.Wrapf(err, "fetching cluster status to determine update status")
	}

	switch status := *c.Status.ClusterInfo.Cluster.Status; status {
	case awseks.ClusterStatusActive:
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

// UnsupportedOIDCError represents an unsupported OIDC error
type UnsupportedOIDCError struct {
	msg string
}

func (u *UnsupportedOIDCError) Error() string {
	return u.msg
}

// NewOpenIDConnectManager returns OpenIDConnectManager
func (c *ClusterProvider) NewOpenIDConnectManager(spec *api.ClusterConfig) (*iamoidc.OpenIDConnectManager, error) {
	if _, err := c.CanOperate(spec); err != nil {
		return nil, err
	}

	if c.Status.ClusterInfo.Cluster == nil || c.Status.ClusterInfo.Cluster.Identity == nil || c.Status.ClusterInfo.Cluster.Identity.Oidc == nil || c.Status.ClusterInfo.Cluster.Identity.Oidc.Issuer == nil {
		return nil, &UnsupportedOIDCError{"unknown OIDC issuer URL"}
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

	return iamoidc.NewOpenIDConnectManager(c.Provider.IAM(), parsedARN.AccountID, *c.Status.ClusterInfo.Cluster.Identity.Oidc.Issuer, parsedARN.Partition)
}

// LoadClusterIntoSpec loads the cluster configuration into the spec
// At the moment VPC and KubernetesNetworkConfig are respected
func (c *ClusterProvider) LoadClusterIntoSpec(spec *api.ClusterConfig) error {
	if err := c.LoadClusterVPC(spec); err != nil {
		return err
	}
	if err := c.RefreshClusterStatus(spec); err != nil {
		return err
	}
	if err := c.loadClusterKubernetesNetworkConfig(spec); err != nil {
		return err
	}
	return nil
}

// LoadClusterVPC loads the VPC configuration
func (c *ClusterProvider) LoadClusterVPC(spec *api.ClusterConfig) error {
	stack, err := c.NewStackManager(spec).DescribeClusterStack()
	if err != nil {
		return err
	}

	return vpc.UseFromCluster(c.Provider, stack, spec)
}

// loadClusterKubernetesNetworkConfig gets the network config of an existing
// cluster, note status must be refreshed!
func (c *ClusterProvider) loadClusterKubernetesNetworkConfig(spec *api.ClusterConfig) error {
	if spec.Status == nil {
		return errors.New("cluster hasn't been refreshed")
	}
	knCfg := c.Status.ClusterInfo.Cluster.KubernetesNetworkConfig
	if knCfg != nil {
		spec.KubernetesNetworkConfig = &api.KubernetesNetworkConfig{
			ServiceIPv4CIDR: aws.StringValue(knCfg.ServiceIpv4Cidr),
		}
	}
	return nil
}

// ListClusters returns a list of the EKS cluster in your account
func (c *ClusterProvider) ListClusters(chunkSize int, listAllRegions bool) ([]*api.ClusterConfig, error) {
	if listAllRegions {
		var clusters []*api.ClusterConfig
		// reset region and re-create the client, then make a recursive call
		for _, region := range api.SupportedRegions() {
			spec := &api.ProviderConfig{
				Region:      region,
				Profile:     c.Provider.Profile(),
				WaitTimeout: c.Provider.WaitTimeout(),
			}
			newClusters, err := New(spec, nil).listClusters(int64(chunkSize))
			if err != nil {
				logger.Critical("error listing clusters in %q region: %s", region, err.Error())
			}

			clusters = append(clusters, newClusters...)
		}
		return clusters, nil
	}

	return c.listClusters(int64(chunkSize))
}

func (c *ClusterProvider) listClusters(chunkSize int64) ([]*api.ClusterConfig, error) {
	allClusters := []*api.ClusterConfig{}

	spec := &api.ClusterConfig{Metadata: &api.ClusterMeta{Name: ""}}
	allStacks, err := c.NewStackManager(spec).ListStackNamesMatching("eksctl-.*-cluster")
	if err != nil {
		return nil, err
	}

	token := ""
	for {
		clusters, nextToken, err := c.getClustersRequest(chunkSize, token)
		if err != nil {
			return nil, err
		}

		for _, clusterName := range clusters {
			spec := &api.ClusterConfig{Metadata: &api.ClusterMeta{Name: *clusterName}}
			isClusterStack, err := c.NewStackManager(spec).IsClusterStackUsingCachedList(allStacks)
			managed := eksctlCreatedFalse
			if err != nil {
				managed = eksctlCreatedUnknown
				logger.Warning("error fetching stacks for cluster %s: %v", clusterName, err)
			} else if isClusterStack {
				managed = eksctlCreatedTrue
			}
			allClusters = append(allClusters, &api.ClusterConfig{
				Metadata: &api.ClusterMeta{
					Name:   *clusterName,
					Region: c.Provider.Region(),
				},
				Status: &api.ClusterStatus{
					EKSCTLCreated: managed,
				},
			})
		}

		if api.IsSetAndNonEmptyString(nextToken) {
			token = *nextToken
		} else {
			break
		}
	}

	return allClusters, nil
}

// GetCluster display details of an EKS cluster in your account
func (c *ClusterProvider) GetCluster(clusterName string) (*awseks.Cluster, error) {
	input := &awseks.DescribeClusterInput{
		Name: &clusterName,
	}

	output, err := c.Provider.EKS().DescribeCluster(input)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to describe control plane %q", clusterName)
	}
	logger.Debug("cluster = %#v", output)

	if *output.Cluster.Status == awseks.ClusterStatusActive {
		if logger.Level >= 4 {
			spec := &api.ClusterConfig{Metadata: &api.ClusterMeta{Name: clusterName}}
			stacks, err := c.NewStackManager(spec).ListStacks()
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

func (c *ClusterProvider) getClustersRequest(chunkSize int64, nextToken string) ([]*string, *string, error) {
	input := &awseks.ListClustersInput{MaxResults: &chunkSize}
	if nextToken != "" {
		input = input.SetNextToken(nextToken)
	}
	manager.ListCount = manager.ListCount + 1
	output, err := c.Provider.EKS().ListClusters(input)
	if err != nil {
		return nil, nil, errors.Wrap(err, "listing control planes")
	}
	return output.Clusters, output.NextToken, nil
}

// WaitForControlPlane waits till the control plane is ready
func (c *ClusterProvider) WaitForControlPlane(meta *api.ClusterMeta, clientSet *kubernetes.Clientset) error {
	if _, err := clientSet.ServerVersion(); err == nil {
		return nil
	}

	condition := func() (bool, error) {
		_, err := clientSet.ServerVersion()
		if err == nil {
			return true, nil
		}
		logger.Debug("control plane not ready yet – %s", err.Error())
		return false, nil
	}

	return waiters.WaitForCondition(c.Provider.WaitTimeout(), time.Second*20, fmt.Errorf("timed out waiting for control plane %q after %s", meta.Name, c.Provider.WaitTimeout()), condition)
}
