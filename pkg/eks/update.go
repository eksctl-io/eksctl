package eks

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/sets"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/utils/waiters"
)

// ClusterVPCConfig represents a cluster's VPC configuration
type ClusterVPCConfig struct {
	ClusterEndpoints  *api.ClusterEndpoints
	PublicAccessCIDRs []string
}

// GetCurrentClusterConfigForLogging fetches current cluster logging configuration as two sets - enabled and disabled types
func (c *ClusterProvider) GetCurrentClusterConfigForLogging(spec *api.ClusterConfig) (sets.String, sets.String, error) {
	enabled := sets.NewString()
	disabled := sets.NewString()

	if ok, err := c.CanOperate(spec); !ok {
		return nil, nil, errors.Wrap(err, "unable to retrieve current cluster logging configuration")
	}

	for _, logTypeGroup := range c.Status.clusterInfo.cluster.Logging.ClusterLogging {
		for _, logType := range logTypeGroup.Types {
			if logType == nil {
				return nil, nil, fmt.Errorf("unexpected response from EKS API - nil string")
			}
			if api.IsEnabled(logTypeGroup.Enabled) {
				enabled.Insert(*logType)
			}
			if api.IsDisabled(logTypeGroup.Enabled) {
				disabled.Insert(*logType)
			}
		}
	}
	return enabled, disabled, nil
}

// UpdateClusterConfigForLogging calls UpdateClusterConfig to enable logging
func (c *ClusterProvider) UpdateClusterConfigForLogging(cfg *api.ClusterConfig) error {
	all := sets.NewString(api.SupportedCloudWatchClusterLogTypes()...)

	enabled := sets.NewString()
	if cfg.HasClusterCloudWatchLogging() {
		enabled.Insert(cfg.CloudWatch.ClusterLogging.EnableTypes...)
	}

	disabled := all.Difference(enabled)

	input := &awseks.UpdateClusterConfigInput{
		Name: &cfg.Metadata.Name,
		Logging: &awseks.Logging{
			ClusterLogging: []*awseks.LogSetup{
				{
					Enabled: api.Enabled(),
					Types:   aws.StringSlice(enabled.List()),
				},
				{
					Enabled: api.Disabled(),
					Types:   aws.StringSlice(disabled.List()),
				},
			},
		},
	}

	output, err := c.Provider.EKS().UpdateClusterConfig(input)
	if err != nil {
		return err
	}
	if err := c.waitForUpdateToSucceed(cfg.Metadata.Name, output.Update); err != nil {
		return err
	}

	describeEnabledTypes := "no types enabled"
	if len(enabled.List()) > 0 {
		describeEnabledTypes = fmt.Sprintf("enabled types: %s", strings.Join(enabled.List(), ", "))
	}

	describeDisabledTypes := "no types disabled"
	if len(disabled.List()) > 0 {
		describeDisabledTypes = fmt.Sprintf("disabled types: %s", strings.Join(disabled.List(), ", "))
	}

	logger.Success("configured CloudWatch logging for cluster %q in %q (%s & %s)",
		cfg.Metadata.Name, cfg.Metadata.Region, describeEnabledTypes, describeDisabledTypes,
	)
	return nil
}

// GetCurrentClusterVPCConfig fetches current cluster endpoint configuration for public and private access types
func (c *ClusterProvider) GetCurrentClusterVPCConfig(spec *api.ClusterConfig) (*ClusterVPCConfig, error) {
	if ok, err := c.CanOperate(spec); !ok {
		return nil, errors.Wrap(err, "unable to retrieve current cluster VPC configuration")
	}

	vpcConfig := c.Status.clusterInfo.cluster.ResourcesVpcConfig

	return &ClusterVPCConfig{
		ClusterEndpoints: &api.ClusterEndpoints{
			PrivateAccess: vpcConfig.EndpointPrivateAccess,
			PublicAccess:  vpcConfig.EndpointPublicAccess,
		},
		PublicAccessCIDRs: aws.StringValueSlice(vpcConfig.PublicAccessCidrs),
	}, nil
}

// UpdateClusterConfigForEndpoints calls eks.UpdateClusterConfig and updates access to API endpoints
func (c *ClusterProvider) UpdateClusterConfigForEndpoints(cfg *api.ClusterConfig) error {

	input := &awseks.UpdateClusterConfigInput{
		Name: &cfg.Metadata.Name,
		ResourcesVpcConfig: &awseks.VpcConfigRequest{
			EndpointPrivateAccess: cfg.VPC.ClusterEndpoints.PrivateAccess,
			EndpointPublicAccess:  cfg.VPC.ClusterEndpoints.PublicAccess,
		},
	}

	output, err := c.Provider.EKS().UpdateClusterConfig(input)
	if err != nil {
		return err
	}

	return c.waitForUpdateToSucceed(cfg.Metadata.Name, output.Update)
}

// UpdatePublicAccessCIDRs calls eks.UpdateClusterConfig and updates the CIDRs for public access
func (c *ClusterProvider) UpdatePublicAccessCIDRs(clusterConfig *api.ClusterConfig) error {
	input := &awseks.UpdateClusterConfigInput{
		Name: &clusterConfig.Metadata.Name,
		ResourcesVpcConfig: &awseks.VpcConfigRequest{
			PublicAccessCidrs: aws.StringSlice(clusterConfig.VPC.PublicAccessCIDRs),
		},
	}
	output, err := c.Provider.EKS().UpdateClusterConfig(input)
	if err != nil {
		return err
	}
	return c.waitForUpdateToSucceed(clusterConfig.Metadata.Name, output.Update)
}

// UpdateClusterVersion calls eks.UpdateClusterVersion and updates to cfg.Metadata.Version,
// it will return update ID along with an error (if it occurs)
func (c *ClusterProvider) UpdateClusterVersion(cfg *api.ClusterConfig) (*awseks.Update, error) {
	input := &awseks.UpdateClusterVersionInput{
		Name:    &cfg.Metadata.Name,
		Version: &cfg.Metadata.Version,
	}
	output, err := c.Provider.EKS().UpdateClusterVersion(input)
	if err != nil {
		return nil, err
	}
	return output.Update, nil
}

// UpdateClusterTags calls eks.TagResource and tags the cluster
func (c *ClusterProvider) UpdateClusterTags(clusterConfig *api.ClusterConfig) error {
	tags := make(map[string]*string)
	for key, value := range clusterConfig.Metadata.Tags {
		tags[key] = &value
	}
	if err := c.RefreshClusterStatus(clusterConfig); err != nil {
		return err
	}
	input := &awseks.TagResourceInput{
		ResourceArn: c.Status.clusterInfo.cluster.Arn,
		Tags:        tags,
	}
	_, err := c.Provider.EKS().TagResource(input)
	if err != nil {
		return err
	}
	var tagStrings []string
	for k, v := range tags {
		tagStrings = append(tagStrings, fmt.Sprintf("%s=%s", k, *v))
	}
	logger.Success("tagged EKS cluster (%s)", strings.Join(tagStrings, ", "))
	return nil
}

// UpdateClusterVersionBlocking calls UpdateClusterVersion and blocks until update
// operation is successful
func (c *ClusterProvider) UpdateClusterVersionBlocking(cfg *api.ClusterConfig) error {
	id, err := c.UpdateClusterVersion(cfg)
	if err != nil {
		return err
	}

	return c.waitForUpdateToSucceed(cfg.Metadata.Name, id)
}

func (c *ClusterProvider) waitForUpdateToSucceed(clusterName string, update *awseks.Update) error {
	newRequest := func() *request.Request {
		input := &awseks.DescribeUpdateInput{
			Name:     &clusterName,
			UpdateId: update.Id,
		}
		req, _ := c.Provider.EKS().DescribeUpdateRequest(input)
		return req
	}

	acceptors := waiters.MakeAcceptors(
		"Update.Status",
		awseks.UpdateStatusSuccessful,
		[]string{
			awseks.UpdateStatusCancelled,
			awseks.UpdateStatusFailed,
		},
	)

	msg := fmt.Sprintf("waiting for requested %q in cluster %q to succeed", *update.Type, clusterName)

	return waiters.Wait(clusterName, msg, acceptors, newRequest, c.Provider.WaitTimeout(), nil)
}
