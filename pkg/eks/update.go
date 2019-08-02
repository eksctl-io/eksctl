package eks

import (
	"fmt"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	awseks "github.com/aws/aws-sdk-go/service/eks"

	"k8s.io/apimachinery/pkg/util/sets"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/utils/waiters"
)

type updateClusterConfigTask struct {
	info string
	skip bool
	spec *api.ClusterConfig
	call func(*api.ClusterConfig) error
}

func (t *updateClusterConfigTask) Skip() bool { return t.skip }
func (t *updateClusterConfigTask) Describe() string {
	if t.skip {
		return "(skip) " + t.info
	}
	return t.info
}
func (t *updateClusterConfigTask) Do(errs chan error) error {
	err := t.call(t.spec)
	close(errs)
	return err
}

// GetCurrentClusterConfigForLogging fetches current cluster logging configuration as two sets - enabled and disabled types
func (c *ClusterProvider) GetCurrentClusterConfigForLogging(cl *api.ClusterMeta) (sets.String, sets.String, error) {
	enabled := sets.NewString()
	disabled := sets.NewString()

	cluster, err := c.DescribeControlPlaneMustBeActive(cl)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to retrieve current cluster logging configuration")
	}

	for _, logTypeGroup := range cluster.Logging.ClusterLogging {
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

// GetUpdateClusterConfigTasks returns all tasks for updating cluster configuration or nil if there are no tasks
func (c *ClusterProvider) GetUpdateClusterConfigTasks(cfg *api.ClusterConfig) *manager.TaskTree {
	if !cfg.HasClusterCloudWatchLogging() {
		logger.Info("CloudWatch logging will not be enabled for cluster %q in %q", cfg.Metadata.Name, cfg.Metadata.Region)
		logger.Info("you can enable it with 'eksctl utils update-cluster-logging --region=%s --name=%s'", cfg.Metadata.Region, cfg.Metadata.Name)
		return nil
	}

	loggingTasks := &manager.TaskTree{Parallel: false}
	loggingTasks.Append(&updateClusterConfigTask{
		info: "update CloudWatch logging configuration",
	})
	return loggingTasks
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
