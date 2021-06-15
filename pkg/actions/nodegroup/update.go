package nodegroup

import (
	"fmt"

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

	if ng.UpdateConfig != nil {
		if describeNodegroupOutput.Nodegroup.UpdateConfig == nil {
			return errors.New("cannot update updateConfig because the nodegroup is not configured to use one")
		}
	}

	logger.Info("nodegroup successfully updated")
	return nil
}
