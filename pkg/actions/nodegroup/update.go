package nodegroup

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/kris-nova/logger"

	"github.com/weaveworks/eksctl/pkg/managed"
)

func (m *Manager) Update(options managed.UpdateOptions) error {
	_, err := m.ctl.Provider.EKS().DescribeNodegroup(&eks.DescribeNodegroupInput{
		ClusterName:   &m.cfg.Metadata.Name,
		NodegroupName: &options.NodegroupName,
	})

	if err != nil {
		if managed.IsNotFound(err) {
			return fmt.Errorf("update is only supported for managed nodegroups; could not find one with name %q", options.NodegroupName)
		}
		return err
	}

	logger.Info("nodegroup successfully updated")
	return nil
}
