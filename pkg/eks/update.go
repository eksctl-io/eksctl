package eks

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	kubeclient "k8s.io/client-go/kubernetes"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/utils/retry"
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

	if ok, err := c.CanOperateWithRefresh(spec); !ok {
		return nil, nil, errors.Wrap(err, "unable to retrieve current cluster logging configuration")
	}

	for _, logTypeGroup := range c.Status.ClusterInfo.Cluster.Logging.ClusterLogging {
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

	input := &eks.UpdateClusterConfigInput{
		Name: &cfg.Metadata.Name,
		Logging: &eks.Logging{
			ClusterLogging: []*eks.LogSetup{
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
	if ok, err := c.CanOperateWithRefresh(spec); !ok {
		return nil, errors.Wrap(err, "unable to retrieve current cluster VPC configuration")
	}

	vpcConfig := c.Status.ClusterInfo.Cluster.ResourcesVpcConfig

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

	input := &eks.UpdateClusterConfigInput{
		Name: &cfg.Metadata.Name,
		ResourcesVpcConfig: &eks.VpcConfigRequest{
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
	input := &eks.UpdateClusterConfigInput{
		Name: &clusterConfig.Metadata.Name,
		ResourcesVpcConfig: &eks.VpcConfigRequest{
			PublicAccessCidrs: aws.StringSlice(clusterConfig.VPC.PublicAccessCIDRs),
		},
	}
	output, err := c.Provider.EKS().UpdateClusterConfig(input)
	if err != nil {
		return err
	}
	return c.waitForUpdateToSucceed(clusterConfig.Metadata.Name, output.Update)
}

// EnableKMSEncryption enables KMS encryption for the specified cluster
func (c *ClusterProvider) EnableKMSEncryption(ctx context.Context, clusterConfig *api.ClusterConfig) error {
	clusterName := aws.String(clusterConfig.Metadata.Name)
	clusterOutput, err := c.Provider.EKS().DescribeCluster(&eks.DescribeClusterInput{
		Name: clusterName,
	})
	if err != nil {
		return errors.Wrap(err, "error describing cluster")
	}
	for _, e := range clusterOutput.Cluster.EncryptionConfig {
		if len(e.Resources) == 1 && *e.Resources[0] == "secrets" {
			if existingKey := *e.Provider.KeyArn; existingKey != clusterConfig.SecretsEncryption.KeyARN {
				return errors.Errorf("KMS encryption is already enabled with key %q, changing the key is not supported", existingKey)
			}
			logger.Info("KMS encryption is already enabled on the cluster")
			return nil
		}
	}

	output, err := c.Provider.EKS().AssociateEncryptionConfigWithContext(ctx, &eks.AssociateEncryptionConfigInput{
		ClusterName: clusterName,
		EncryptionConfig: []*eks.EncryptionConfig{
			{
				Resources: aws.StringSlice([]string{"secrets"}),
				Provider: &eks.Provider{
					KeyArn: aws.String(clusterConfig.SecretsEncryption.KeyARN),
				},
			},
		},
	})

	if err != nil {
		return errors.Wrap(err, "error enabling KMS encryption")
	}

	logger.Info("initiated KMS encryption, this may take up to 45 minutes to complete")

	err = waitForUpdate(ctx, c.Provider.EKS(), &eks.DescribeUpdateInput{
		Name:     clusterName,
		UpdateId: output.Update.Id,
	})

	switch e := err.(type) {
	case *updateFailedError:
		if e.Status == eks.UpdateStatusCancelled {
			return errors.Errorf("request to enable KMS encryption was cancelled: %s", e.UpdateError)
		}
		return errors.Errorf("failed to enable KMS encryption: %s", e.UpdateError)

	case nil:
		logger.Info("KMS encryption successfully enabled on cluster %q", clusterConfig.Metadata.Name)
		return nil

	default:
		return err
	}
}

type updateFailedError struct {
	Status      string
	UpdateError string
}

func (u *updateFailedError) Error() string {
	return fmt.Sprintf("update failed with status %q: %s", u.Status, u.UpdateError)
}

func waitForUpdate(ctx context.Context, eksAPI eksiface.EKSAPI, input *eks.DescribeUpdateInput) error {
	logger.Debug("waiting for update to complete (updateID: %v)", *input.UpdateId)

	const retryAfter = 20 * time.Second

	for {
		describeOutput, err := eksAPI.DescribeUpdate(input)

		if err != nil {
			describeErr := errors.Wrap(err, "error describing nodegroup update")
			if !request.IsErrorRetryable(err) {
				return describeErr
			}
			logger.Warning(describeErr.Error())
		}

		logger.Debug("DescribeUpdate output: %v", describeOutput.Update.String())

		switch status := *describeOutput.Update.Status; status {
		case eks.UpdateStatusSuccessful:
			return nil

		case eks.UpdateStatusCancelled, eks.UpdateStatusFailed:
			return &updateFailedError{
				Status:      status,
				UpdateError: fmt.Sprintf("update errors:\n%s", aggregateErrors(describeOutput.Update.Errors)),
			}

		case eks.UpdateStatusInProgress:
			logger.Debug("update in progress")

		default:
			return errors.Errorf("unexpected update status: %q", status)

		}

		timer := time.NewTimer(retryAfter)
		select {
		case <-ctx.Done():
			timer.Stop()
			return errors.Errorf("timed out waiting for update to complete: %v", ctx.Err())
		case <-timer.C:
		}
	}
}

func aggregateErrors(errorDetails []*eks.ErrorDetail) string {
	var aggregatedErrors []string
	for _, err := range errorDetails {
		aggregatedErrors = append(aggregatedErrors, fmt.Sprintf("- %s", err.String()))
	}
	return strings.Join(aggregatedErrors, "\n")
}

// UpdateClusterVersion calls eks.UpdateClusterVersion and updates to cfg.Metadata.Version,
// it will return update ID along with an error (if it occurs)
func (c *ClusterProvider) UpdateClusterVersion(cfg *api.ClusterConfig) (*eks.Update, error) {
	input := &eks.UpdateClusterVersionInput{
		Name:    &cfg.Metadata.Name,
		Version: &cfg.Metadata.Version,
	}
	output, err := c.Provider.EKS().UpdateClusterVersion(input)
	if err != nil {
		return nil, err
	}
	return output.Update, nil
}

// UpdateClusterVersionBlocking calls UpdateClusterVersion and blocks until update
// operation is successful
func (c *ClusterProvider) UpdateClusterVersionBlocking(cfg *api.ClusterConfig) error {
	id, err := c.UpdateClusterVersion(cfg)
	if err != nil {
		return err
	}

	if err := c.waitForUpdateToSucceed(cfg.Metadata.Name, id); err != nil {
		return err
	}

	return c.waitForControlPlaneVersion(cfg)
}

func (c *ClusterProvider) waitForUpdateToSucceed(clusterName string, update *eks.Update) error {
	newRequest := func() *request.Request {
		input := &eks.DescribeUpdateInput{
			Name:     &clusterName,
			UpdateId: update.Id,
		}
		req, _ := c.Provider.EKS().DescribeUpdateRequest(input)
		return req
	}

	acceptors := waiters.MakeAcceptors(
		"Update.Status",
		eks.UpdateStatusSuccessful,
		[]string{
			eks.UpdateStatusCancelled,
			eks.UpdateStatusFailed,
		},
	)

	msg := fmt.Sprintf("waiting for requested %q in cluster %q to succeed", *update.Type, clusterName)

	return waiters.Wait(clusterName, msg, acceptors, newRequest, c.Provider.WaitTimeout(), nil)
}

func controlPlaneIsVersion(clientSet *kubeclient.Clientset, version string) (bool, error) {
	serverVersion, err := clientSet.ServerVersion()
	if err != nil {
		return false, err
	}
	return fmt.Sprintf("%s.%s", serverVersion.Major, strings.TrimSuffix(serverVersion.Minor, "+")) == version, nil
}

func (c *ClusterProvider) waitForControlPlaneVersion(cfg *api.ClusterConfig) error {
	retryPolicy := retry.TimingOutExponentialBackoff{
		Timeout:  c.Provider.WaitTimeout(),
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
