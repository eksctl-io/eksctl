package nodegroup

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	"github.com/weaveworks/eksctl/pkg/managed"
)

func (m *Manager) Update() error {
	ng := m.cfg.ManagedNodeGroups[0]
	describeNodegroupOutput, err := m.ctl.Provider.EKS().DescribeNodegroup(&eks.DescribeNodegroupInput{
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
		return fmt.Errorf("the submitted config didn't contain changes for nodegroup %s", ng.Name)
	}

	if describeNodegroupOutput.Nodegroup.UpdateConfig == nil {
		return errors.New("cannot update updateConfig because the nodegroup is not configured to use one")
	}

	logger.Info("updating nodegroup %s's UpdateConfig", ng.Name)
	updateConfig := &eks.NodegroupUpdateConfig{}

	if ng.UpdateConfig.MaxUnavailable != nil {
		updateConfig.MaxUnavailable = aws.Int64(int64(*ng.UpdateConfig.MaxUnavailable))
	}

	if ng.UpdateConfig.MaxUnavailablePercentage != nil {
		updateConfig.MaxUnavailablePercentage = aws.Int64(int64(*ng.UpdateConfig.MaxUnavailablePercentage))
	}

	_, err = m.ctl.Provider.EKS().UpdateNodegroupConfig(&eks.UpdateNodegroupConfigInput{
		UpdateConfig:  updateConfig,
		ClusterName:   &m.cfg.Metadata.Name,
		NodegroupName: &ng.Name,
	})
	if err != nil {
		return err
	}

	logger.Info("nodegroup %s successfully updated", ng.Name)
	return nil
}
