package eks

import (
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	awseks "github.com/aws/aws-sdk-go/service/eks"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/fargate"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/printers"
	"github.com/weaveworks/eksctl/pkg/vpc"
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

	return ClusterSupportsManagedNodes(c.Status.clusterInfo.cluster)
}

// ClusterSupportsManagedNodes reports whether the EKS cluster supports managed nodes
func ClusterSupportsManagedNodes(cluster *awseks.Cluster) (bool, error) {
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
	return ClusterSupportsFargate(c.Status.clusterInfo.cluster)
}

// ClusterSupportsFargate reports whether an existing cluster supports Fargate.
func ClusterSupportsFargate(cluster *awseks.Cluster) (bool, error) {
	versionSupportsFargate, err := fargate.IsSupportedBy(*cluster.Version)
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
		return false, errors.Wrapf(err, "fetching cluster status to determine if it can be deleted")
	}
	// it must be possible to delete cluster in any state
	return true, nil
}

// CanOperate return true when a cluster can be operated, otherwise it returns false along with an error explaining the reason
func (c *ClusterProvider) CanOperate(spec *api.ClusterConfig) (bool, error) {
	err := c.maybeRefreshClusterStatus(spec)
	if err != nil {
		return false, errors.Wrapf(err, "fetching cluster status to determine operability")
	}

	switch status := *c.Status.clusterInfo.cluster.Status; status {
	case awseks.ClusterStatusCreating, awseks.ClusterStatusDeleting, awseks.ClusterStatusFailed:
		return false, fmt.Errorf("cannot perform Kubernetes API operations on cluster %q in %q region due to status %q", spec.Metadata.Name, spec.Metadata.Region, status)
	default:
		// all other states are considered operable, including UPDGRADING (which is missing from the SDK)
		return true, nil
	}
}

// CanUpdate return true when a cluster or add-ons can be updated, otherwise it returns false along with an error explaining the reason
func (c *ClusterProvider) CanUpdate(spec *api.ClusterConfig) (bool, error) {
	err := c.maybeRefreshClusterStatus(spec)
	if err != nil {
		return false, errors.Wrapf(err, "fetching cluster status to determine update status")
	}

	switch status := *c.Status.clusterInfo.cluster.Status; status {
	case awseks.ClusterStatusActive:
		// only active cluster can be upgraded
		return true, nil
	default:
		return false, fmt.Errorf("cannot update cluster %q in %q region due to status %q", spec.Metadata.Name, spec.Metadata.Region, status)
	}
}

// ControlPlaneVersion returns cached version (EKS API)
func (c *ClusterProvider) ControlPlaneVersion() string {
	if c.Status.clusterInfo.cluster == nil || c.Status.clusterInfo.cluster.Version == nil {
		return ""
	}
	return *c.Status.clusterInfo.cluster.Version
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

	switch c.ControlPlaneVersion() {
	case "", api.Version1_10, api.Version1_11, api.Version1_12:
		return nil, &UnsupportedOIDCError{"OIDC is only supported in Kubernetes versions 1.13 and above"}
	}
	if c.Status.clusterInfo.cluster == nil || c.Status.clusterInfo.cluster.Identity == nil || c.Status.clusterInfo.cluster.Identity.Oidc == nil || c.Status.clusterInfo.cluster.Identity.Oidc.Issuer == nil {
		return nil, &UnsupportedOIDCError{"unknown OIDC issuer URL"}
	}

	parsedARN, err := arn.Parse(spec.Status.ARN)
	if err != nil {
		return nil, errors.Wrapf(err, "unexpected invalid ARN: %q", spec.Status.ARN)
	}
	switch parsedARN.Partition {
	case "aws", "aws-cn":
	default:
		return nil, fmt.Errorf("unknown EKS ARN: %q", spec.Status.ARN)
	}
	accountID := strings.Split(spec.Status.ARN, ":")[4]
	return iamoidc.NewOpenIDConnectManager(c.Provider.IAM(), accountID, *c.Status.clusterInfo.cluster.Identity.Oidc.Issuer)
}

// LoadClusterVPC loads the VPC configuration
func (c *ClusterProvider) LoadClusterVPC(spec *api.ClusterConfig) error {
	stack, err := c.NewStackManager(spec).DescribeClusterStack()
	if err != nil {
		return err
	}

	return vpc.UseFromCluster(c.Provider, stack, spec)
}

// ListClusters display details of all the EKS cluster in your account
func (c *ClusterProvider) ListClusters(clusterName string, chunkSize int, output printers.Type, eachRegion bool) error {
	// NOTE: this needs to be reworked in the future so that the functionality
	// is combined. This require the ability to return details of all clusters
	// in a single call.
	printer, err := printers.NewPrinter(output)
	if err != nil {
		return err
	}

	if clusterName != "" {
		if output == "table" {
			addSummaryTableColumns(printer.(*printers.TablePrinter))
		}
		return c.doGetCluster(clusterName, printer)
	}

	if output == "table" {
		addListTableColumns(printer.(*printers.TablePrinter))
	}
	allClusters := []*api.ClusterMeta{}
	if err := c.doListClusters(int64(chunkSize), printer, &allClusters, eachRegion); err != nil {
		return err
	}
	return printer.PrintObjWithKind("clusters", allClusters, os.Stdout)
}

func (c *ClusterProvider) getClustersRequest(chunkSize int64, nextToken string) ([]*string, *string, error) {
	input := &awseks.ListClustersInput{MaxResults: &chunkSize}
	if nextToken != "" {
		input = input.SetNextToken(nextToken)
	}
	output, err := c.Provider.EKS().ListClusters(input)
	if err != nil {
		return nil, nil, errors.Wrap(err, "listing control planes")
	}
	return output.Clusters, output.NextToken, nil
}

func (c *ClusterProvider) doListClusters(chunkSize int64, printer printers.OutputPrinter, allClusters *[]*api.ClusterMeta, eachRegion bool) error {
	if eachRegion {
		// reset region and re-create the client, then make a recursive call
		for _, region := range api.SupportedRegions() {
			spec := &api.ProviderConfig{
				Region:      region,
				Profile:     c.Provider.Profile(),
				WaitTimeout: c.Provider.WaitTimeout(),
			}
			if err := New(spec, nil).doListClusters(chunkSize, printer, allClusters, false); err != nil {
				logger.Critical("error listing clusters in %q region: %s", region, err.Error())
			}
		}
		return nil
	}

	token := ""
	for {
		clusters, nextToken, err := c.getClustersRequest(chunkSize, token)
		if err != nil {
			return err
		}

		for _, clusterName := range clusters {
			*allClusters = append(*allClusters, &api.ClusterMeta{
				Name:   *clusterName,
				Region: c.Provider.Region(),
			})
		}

		if api.IsSetAndNonEmptyString(nextToken) {
			token = *nextToken
		} else {
			break
		}
	}

	return nil
}

func (c *ClusterProvider) doGetCluster(clusterName string, printer printers.OutputPrinter) error {
	input := &awseks.DescribeClusterInput{
		Name: &clusterName,
	}
	output, err := c.Provider.EKS().DescribeCluster(input)
	if err != nil {
		return errors.Wrapf(err, "unable to describe control plane %q", clusterName)
	}
	logger.Debug("cluster = %#v", output)

	clusters := []*awseks.Cluster{output.Cluster} // TODO: in the future this will have multiple clusters
	if err := printer.PrintObjWithKind("clusters", clusters, os.Stdout); err != nil {
		return err
	}

	if *output.Cluster.Status == awseks.ClusterStatusActive {
		if logger.Level >= 4 {
			spec := &api.ClusterConfig{Metadata: &api.ClusterMeta{Name: clusterName}}
			stacks, err := c.NewStackManager(spec).ListStacks()
			if err != nil {
				return errors.Wrapf(err, "listing CloudFormation stack for %q", clusterName)
			}
			for _, s := range stacks {
				logger.Debug("stack = %#v", *s)
			}
		}
	}
	return nil
}

// WaitForControlPlane waits till the control plane is ready
func (c *ClusterProvider) WaitForControlPlane(meta *api.ClusterMeta, clientSet *kubernetes.Clientset) error {
	if _, err := clientSet.ServerVersion(); err == nil {
		return nil
	}

	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	timer := time.NewTimer(c.Provider.WaitTimeout())
	defer timer.Stop()

	for {
		select {
		case <-ticker.C:
			_, err := clientSet.ServerVersion()
			if err == nil {
				return nil
			}
			logger.Debug("control plane not ready yet – %s", err.Error())
		case <-timer.C:
			return fmt.Errorf("timed out waiting for control plane %q after %s", meta.Name, c.Provider.WaitTimeout())
		}
	}
}

func addSummaryTableColumns(printer *printers.TablePrinter) {
	printer.AddColumn("NAME", func(c *awseks.Cluster) string {
		return *c.Name
	})
	printer.AddColumn("VERSION", func(c *awseks.Cluster) string {
		return *c.Version
	})
	printer.AddColumn("STATUS", func(c *awseks.Cluster) string {
		return *c.Status
	})
	printer.AddColumn("CREATED", func(c *awseks.Cluster) string {
		return c.CreatedAt.Format(time.RFC3339)
	})
	printer.AddColumn("VPC", func(c *awseks.Cluster) string {
		return *c.ResourcesVpcConfig.VpcId
	})
	printer.AddColumn("SUBNETS", func(c *awseks.Cluster) string {
		subnets := sets.NewString()
		for _, subnetid := range c.ResourcesVpcConfig.SubnetIds {
			if api.IsSetAndNonEmptyString(subnetid) {
				subnets.Insert(*subnetid)
			}
		}
		return strings.Join(subnets.List(), ",")
	})
	printer.AddColumn("SECURITYGROUPS", func(c *awseks.Cluster) string {
		groups := sets.NewString()
		for _, sg := range c.ResourcesVpcConfig.SecurityGroupIds {
			if api.IsSetAndNonEmptyString(sg) {
				groups.Insert(*sg)
			}
		}
		return strings.Join(groups.List(), ",")
	})
}

func addListTableColumns(printer *printers.TablePrinter) {
	printer.AddColumn("NAME", func(c *api.ClusterMeta) string {
		return c.Name
	})
	printer.AddColumn("REGION", func(c *api.ClusterMeta) string {
		return c.Region
	})
}
