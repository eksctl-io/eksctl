package nodegroup

import (
	"context"
	"fmt"
	"time"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/waiter"
	"github.com/weaveworks/eksctl/pkg/managed"
)

func (m *Manager) Update(ctx context.Context, wait bool) error {
	for _, ng := range m.cfg.ManagedNodeGroups {
		if err := m.updateNodegroup(ctx, ng, wait); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) updateNodegroup(ctx context.Context, ng *api.ManagedNodeGroup, wait bool) error {
	logger.Info("checking that nodegroup %s is a managed nodegroup", ng.Name)

	_, err := m.ctl.AWSProvider.EKS().DescribeNodegroup(ctx, &eks.DescribeNodegroupInput{
		ClusterName:   &m.cfg.Metadata.Name,
		NodegroupName: &ng.Name,
	})

	if err != nil {
		if managed.IsNotFound(err) {
			return fmt.Errorf("could not find managed nodegroup with name %q", ng.Name)
		}
		return err
	}

	if ng.UpdateConfig == nil {
		return fmt.Errorf("the submitted config does not contain an 'updateConfig' field for nodegroup %s", ng.Name)
	}

	updateConfig, err := updateUpdateConfig(ng)
	if err != nil {
		return err
	}

	output, err := m.ctl.AWSProvider.EKS().UpdateNodegroupConfig(ctx, &eks.UpdateNodegroupConfigInput{
		UpdateConfig:  updateConfig,
		ClusterName:   &m.cfg.Metadata.Name,
		NodegroupName: &ng.Name,
	})
	if err != nil {
		return fmt.Errorf("failed to update nodegroup %s: %w", ng.Name, err)
	}

	if wait {
		if status, err := waiter.WaitForNodegroupUpdate(ctx, string(output.Update.Status), m.ctl.AWSProvider.EKS(), m.ctl.AWSProvider.WaitTimeout(), func(attempts int) time.Duration {
			return 30 * time.Second
		}); err != nil {
			return fmt.Errorf("failed to wait for nodegroup %s to update; last observed status was %s with error: %w", ng.Name, status, err)
		}
	}

	logger.Info("nodegroup %s successfully updated", ng.Name)
	return nil
}

func updateUpdateConfig(ng *api.ManagedNodeGroup) (*ekstypes.NodegroupUpdateConfig, error) {
	logger.Info("updating nodegroup %s's UpdateConfig", ng.Name)
	updateConfig := &ekstypes.NodegroupUpdateConfig{}

	if ng.UpdateConfig.MaxUnavailable != nil {
		updateConfig.MaxUnavailable = aws.Int32(int32(*ng.UpdateConfig.MaxUnavailable))
	}

	if ng.UpdateConfig.MaxUnavailablePercentage != nil {
		updateConfig.MaxUnavailablePercentage = aws.Int32(int32(*ng.UpdateConfig.MaxUnavailablePercentage))
	}

	return updateConfig, nil
}
