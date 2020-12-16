package nodegroup

import (
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/managed"
)

func (m *Manager) Upgrade(nodeGroupName, version, launchTemplateVersion string, forceUpgrade bool) error {
	stackCollection := manager.NewStackCollection(m.ctl.Provider, m.cfg)
	hasStacks, err := m.hasStacks(nodeGroupName)
	if err != nil {
		return err
	}

	if hasStacks {
		managedService := managed.NewService(m.ctl.Provider, stackCollection, m.cfg.Metadata.Name)
		return managedService.UpgradeNodeGroup(managed.UpgradeOptions{
			NodegroupName:         nodeGroupName,
			KubernetesVersion:     version,
			LaunchTemplateVersion: launchTemplateVersion,
			ForceUpgrade:          forceUpgrade,
		})
	}

	_, err = m.ctl.Provider.EKS().UpdateNodegroupVersion(&eks.UpdateNodegroupVersionInput{
		ClusterName:   &m.cfg.Metadata.Name,
		Force:         &forceUpgrade,
		NodegroupName: &nodeGroupName,
		Version:       &version,
	})

	return err
}
