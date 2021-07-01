package nodegroup

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/managed"
)

func (m *Manager) Update() error {
	for _, ng := range m.cfg.ManagedNodeGroups {
		if err := m.updateNodegroup(ng); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) updateNodegroup(ng *api.ManagedNodeGroup) error {
	logger.Info("checking that nodegroup %s is a managed nodegroup", ng.Name)

	_, err := m.ctl.Provider.EKS().DescribeNodegroup(&eks.DescribeNodegroupInput{
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
		return fmt.Errorf("the submitted config does not contain any changes for nodegroup %s", ng.Name)
	}

	updateConfig, err := updateUpdateConfig(ng)
	if err != nil {
		return err
	}

	_, err = m.ctl.Provider.EKS().UpdateNodegroupConfig(&eks.UpdateNodegroupConfigInput{
		UpdateConfig:  updateConfig,
		ClusterName:   &m.cfg.Metadata.Name,
		NodegroupName: &ng.Name,
	})
	if err != nil {
		return fmt.Errorf("failed to update nodegroup %s: %w", ng.Name, err)
	}

	logger.Info("nodegroup %s successfully updated", ng.Name)
	return nil
}

func updateUpdateConfig(ng *api.ManagedNodeGroup) (*eks.NodegroupUpdateConfig, error) {
	logger.Info("updating nodegroup %s's UpdateConfig", ng.Name)
	updateConfig := &eks.NodegroupUpdateConfig{}

	if ng.UpdateConfig.MaxUnavailable != nil {
		updateConfig.MaxUnavailable = aws.Int64(int64(*ng.UpdateConfig.MaxUnavailable))
	}

	if ng.UpdateConfig.MaxUnavailablePercentage != nil {
		updateConfig.MaxUnavailablePercentage = aws.Int64(int64(*ng.UpdateConfig.MaxUnavailablePercentage))
	}

	return updateConfig, nil
}
