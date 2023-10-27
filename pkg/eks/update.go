package eks

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	kubeclient "k8s.io/client-go/kubernetes"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks/waiter"
	"github.com/weaveworks/eksctl/pkg/utils/retry"
)

// ClusterVPCConfig represents a cluster's VPC configuration
type ClusterVPCConfig struct {
	ClusterEndpoints  *api.ClusterEndpoints
	PublicAccessCIDRs []string
}

// GetCurrentClusterConfigForLogging fetches current cluster logging configuration as two sets - enabled and disabled types
func (c *ClusterProvider) GetCurrentClusterConfigForLogging(ctx context.Context, spec *api.ClusterConfig) (sets.String, sets.String, error) {
	enabled := sets.NewString()
	disabled := sets.NewString()

	if ok, err := c.CanOperateWithRefresh(ctx, spec); !ok {
		return nil, nil, errors.Wrap(err, "unable to retrieve current cluster logging configuration")
	}

	for _, logTypeGroup := range c.Status.ClusterInfo.Cluster.Logging.ClusterLogging {
		for _, lt := range logTypeGroup.Types {
			logType := string(lt)
			if api.IsEnabled(logTypeGroup.Enabled) {
				enabled.Insert(logType)
			}
			if api.IsDisabled(logTypeGroup.Enabled) {
				disabled.Insert(logType)
			}
		}
	}
	return enabled, disabled, nil
}

// UpdateClusterConfigForLogging calls UpdateClusterConfig to enable logging
func (c *ClusterProvider) UpdateClusterConfigForLogging(ctx context.Context, cfg *api.ClusterConfig) error {
	all := sets.NewString(api.SupportedCloudWatchClusterLogTypes()...)

	enabled := sets.NewString()
	if cfg.HasClusterCloudWatchLogging() {
		enabled.Insert(cfg.CloudWatch.ClusterLogging.EnableTypes...)
	}

	disabled := all.Difference(enabled)

	toLogTypes := func(logTypes sets.String) []ekstypes.LogType {
		ret := make([]ekstypes.LogType, len(logTypes))
		for i, logType := range logTypes.List() {
			ret[i] = ekstypes.LogType(logType)
		}
		return ret
	}

	input := &eks.UpdateClusterConfigInput{
		Name: &cfg.Metadata.Name,
		Logging: &ekstypes.Logging{
			ClusterLogging: []ekstypes.LogSetup{
				{
					Enabled: api.Enabled(),
					Types:   toLogTypes(enabled),
				},
				{
					Enabled: api.Disabled(),
					Types:   toLogTypes(disabled),
				},
			},
		},
	}
	if err := c.UpdateClusterConfig(ctx, input); err != nil {
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

	if logRetentionInDays := cfg.CloudWatch.ClusterLogging.LogRetentionInDays; logRetentionInDays > 0 {
		if _, err := c.AWSProvider.CloudWatchLogs().PutRetentionPolicy(ctx, &cloudwatchlogs.PutRetentionPolicyInput{
			// The format for log group name is documented here: https://docs.aws.amazon.com/eks/latest/userguide/control-plane-logs.html
			LogGroupName:    aws.String(fmt.Sprintf("/aws/eks/%s/cluster", cfg.Metadata.Name)),
			RetentionInDays: aws.Int32(int32(logRetentionInDays)),
		}); err != nil {
			return fmt.Errorf("error updating log retention settings: %w", err)
		}
		logger.Success("configured CloudWatch log retention to %d days for CloudWatch logging", logRetentionInDays)
	}
	return nil
}

// GetCurrentClusterVPCConfig fetches current cluster endpoint configuration for public and private access types
func (c *ClusterProvider) GetCurrentClusterVPCConfig(ctx context.Context, spec *api.ClusterConfig) (*ClusterVPCConfig, error) {
	if ok, err := c.CanOperateWithRefresh(ctx, spec); !ok {
		return nil, errors.Wrap(err, "unable to retrieve current cluster VPC configuration")
	}

	vpcConfig := c.Status.ClusterInfo.Cluster.ResourcesVpcConfig

	return &ClusterVPCConfig{
		ClusterEndpoints: &api.ClusterEndpoints{
			PrivateAccess: &vpcConfig.EndpointPrivateAccess,
			PublicAccess:  &vpcConfig.EndpointPublicAccess,
		},
		PublicAccessCIDRs: vpcConfig.PublicAccessCidrs,
	}, nil
}

// UpdateClusterConfigForEndpoints calls eks.UpdateClusterConfig and updates access to API endpoints
func (c *ClusterProvider) UpdateClusterConfigForEndpoints(ctx context.Context, cfg *api.ClusterConfig) error {

	input := &eks.UpdateClusterConfigInput{
		Name: &cfg.Metadata.Name,
		ResourcesVpcConfig: &ekstypes.VpcConfigRequest{
			EndpointPrivateAccess: cfg.VPC.ClusterEndpoints.PrivateAccess,
			EndpointPublicAccess:  cfg.VPC.ClusterEndpoints.PublicAccess,
		},
	}

	return c.UpdateClusterConfig(ctx, input)
}

// UpdatePublicAccessCIDRs calls eks.UpdateClusterConfig and updates the CIDRs for public access
func (c *ClusterProvider) UpdatePublicAccessCIDRs(ctx context.Context, clusterConfig *api.ClusterConfig) error {
	input := &eks.UpdateClusterConfigInput{
		Name: &clusterConfig.Metadata.Name,
		ResourcesVpcConfig: &ekstypes.VpcConfigRequest{
			PublicAccessCidrs: clusterConfig.VPC.PublicAccessCIDRs,
		},
	}
	return c.UpdateClusterConfig(ctx, input)
}

// UpdateClusterConfig calls EKS.UpdateClusterConfig and waits for the update to complete.
func (c *ClusterProvider) UpdateClusterConfig(ctx context.Context, input *eks.UpdateClusterConfigInput) error {
	output, err := c.AWSProvider.EKS().UpdateClusterConfig(ctx, input)
	if err != nil {
		return err
	}
	return c.waitForUpdateToSucceed(ctx, *input.Name, output.Update)
}

// EnableKMSEncryption enables KMS encryption for the specified cluster
func (c *ClusterProvider) EnableKMSEncryption(ctx context.Context, clusterConfig *api.ClusterConfig) error {
	clusterName := aws.String(clusterConfig.Metadata.Name)
	clusterOutput, err := c.AWSProvider.EKS().DescribeCluster(ctx, &eks.DescribeClusterInput{
		Name: clusterName,
	})
	if err != nil {
		return errors.Wrap(err, "error describing cluster")
	}
	for _, e := range clusterOutput.Cluster.EncryptionConfig {
		if len(e.Resources) == 1 && e.Resources[0] == "secrets" {
			if existingKey := *e.Provider.KeyArn; existingKey != clusterConfig.SecretsEncryption.KeyARN {
				return errors.Errorf("KMS encryption is already enabled with key %q, changing the key is not supported", existingKey)
			}
			logger.Info("KMS encryption is already enabled on the cluster")
			return nil
		}
	}

	output, err := c.AWSProvider.EKS().AssociateEncryptionConfig(ctx, &eks.AssociateEncryptionConfigInput{
		ClusterName: clusterName,
		EncryptionConfig: []ekstypes.EncryptionConfig{
			{
				Resources: []string{"secrets"},
				Provider: &ekstypes.Provider{
					KeyArn: aws.String(clusterConfig.SecretsEncryption.KeyARN),
				},
			},
		},
	})

	if err != nil {
		return errors.Wrap(err, "error enabling KMS encryption")
	}

	logger.Info("initiated KMS encryption, this may take up to 45 minutes to complete")

	updateWaiter := waiter.NewUpdateWaiter(c.AWSProvider.EKS(), func(options *waiter.UpdateWaiterOptions) {
		options.RetryAttemptLogMessage = fmt.Sprintf("waiting for update %q in cluster %q to complete", *output.Update.Id, *clusterName)
	})
	err = updateWaiter.Wait(ctx, &eks.DescribeUpdateInput{
		Name:     clusterName,
		UpdateId: output.Update.Id,
	}, c.AWSProvider.WaitTimeout())

	switch e := err.(type) {
	case *waiter.UpdateFailedError:
		if e.Status == string(ekstypes.UpdateStatusCancelled) {
			return fmt.Errorf("request to enable KMS encryption was cancelled: %s", e.UpdateError)
		}
		return fmt.Errorf("failed to enable KMS encryption: %s", e.UpdateError)

	case nil:
		logger.Info("KMS encryption successfully enabled on cluster %q", clusterConfig.Metadata.Name)
		return nil

	default:
		return err
	}
}

// UpdateClusterVersion calls eks.UpdateClusterVersion and updates to cfg.Metadata.Version,
// it will return update ID along with an error (if it occurs)
func (c *ClusterProvider) UpdateClusterVersion(ctx context.Context, cfg *api.ClusterConfig) (*ekstypes.Update, error) {
	input := &eks.UpdateClusterVersionInput{
		Name:    &cfg.Metadata.Name,
		Version: &cfg.Metadata.Version,
	}
	output, err := c.AWSProvider.EKS().UpdateClusterVersion(ctx, input)
	if err != nil {
		return nil, err
	}
	return output.Update, nil
}

// UpdateClusterVersionBlocking calls UpdateClusterVersion and blocks until update
// operation is successful
func (c *ClusterProvider) UpdateClusterVersionBlocking(ctx context.Context, cfg *api.ClusterConfig) error {
	id, err := c.UpdateClusterVersion(ctx, cfg)
	if err != nil {
		return err
	}

	if err := c.waitForUpdateToSucceed(ctx, cfg.Metadata.Name, id); err != nil {
		return err
	}

	return c.waitForControlPlaneVersion(cfg)
}

func (c *ClusterProvider) waitForUpdateToSucceed(ctx context.Context, clusterName string, update *ekstypes.Update) error {
	updateWaiter := waiter.NewUpdateWaiter(c.AWSProvider.EKS(), func(options *waiter.UpdateWaiterOptions) {
		options.RetryAttemptLogMessage = fmt.Sprintf("waiting for requested %q in cluster %q to succeed", update.Type, clusterName)
	})
	return updateWaiter.Wait(ctx, &eks.DescribeUpdateInput{
		Name:     &clusterName,
		UpdateId: update.Id,
	}, c.AWSProvider.WaitTimeout())
}

func controlPlaneIsVersion(clientSet kubeclient.Interface, version string) (bool, error) {
	serverVersion, err := clientSet.Discovery().ServerVersion()
	if err != nil {
		return false, err
	}
	return fmt.Sprintf("%s.%s", serverVersion.Major, strings.TrimSuffix(serverVersion.Minor, "+")) == version, nil
}

func (c *ClusterProvider) waitForControlPlaneVersion(cfg *api.ClusterConfig) error {
	retryPolicy := retry.TimingOutExponentialBackoff{
		Timeout:  c.AWSProvider.WaitTimeout(),
		TimeUnit: time.Second,
	}

	clientSet, err := c.NewStdClientSet(cfg)
	if err != nil {
		return err
	}

	for !retryPolicy.Done() {
		isUpdated, err := controlPlaneIsVersion(clientSet, cfg.Metadata.Version)
		if err != nil {
			return err
		}
		if isUpdated {
			return nil
		}
		time.Sleep(retryPolicy.Duration())
	}
	return errors.New("timed out while waiting for control plane to report updated version")
}
