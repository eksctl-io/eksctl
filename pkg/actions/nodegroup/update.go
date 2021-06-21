package nodegroup

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/kris-nova/logger"

	"github.com/weaveworks/eksctl/pkg/managed"
)

func (m *Manager) Update() error {
	ngName := m.cfg.ManagedNodeGroups[0].Name
	_, err := m.ctl.Provider.EKS().DescribeNodegroup(&eks.DescribeNodegroupInput{
		ClusterName:   &m.cfg.Metadata.Name,
		NodegroupName: &ngName,
	})

	if err != nil {
		if managed.IsNotFound(err) {
			return fmt.Errorf("could not find managed nodegroup with name %q", ngName)
		}
		return err
	}

	logger.Info("nodegroup successfully updated")
	return nil
}
