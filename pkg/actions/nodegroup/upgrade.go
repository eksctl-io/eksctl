package nodegroup

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/blang/semver"
	"github.com/kris-nova/logger"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/managed"
	"github.com/weaveworks/eksctl/pkg/utils/waiters"
)

func (m *Manager) Upgrade(nodeGroupName, version, launchTemplateVersion string, forceUpgrade bool) error {
	stackCollection := manager.NewStackCollection(m.ctl.Provider, m.cfg)
	hasStacks, err := m.hasStacks(nodeGroupName)
	if err != nil {
		return err
	}

	if version != "" {
		if _, err := semver.ParseTolerant(version); err != nil {
			return errors.Wrap(err, "invalid Kubernetes version")
		}
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

	return m.upgradeAndWait(nodeGroupName, version, launchTemplateVersion, forceUpgrade)
}

func (m *Manager) upgradeAndWait(nodeGroupName, version, launchTemplateVersion string, forceUpgrade bool) error {
	input := &eks.UpdateNodegroupVersionInput{
		ClusterName:   &m.cfg.Metadata.Name,
		Force:         &forceUpgrade,
		NodegroupName: &nodeGroupName,
		Version:       &version,
	}

	describeNodegroupOutput, err := m.ctl.Provider.EKS().DescribeNodegroup(&eks.DescribeNodegroupInput{
		ClusterName:   &m.cfg.Metadata.Name,
		NodegroupName: &nodeGroupName,
	})

	if err != nil {
		return err
	}

	if launchTemplateVersion != "" {
		lt := describeNodegroupOutput.Nodegroup.LaunchTemplate
		if lt == nil || (lt.Id == nil && lt.Name == nil) {
			return errors.New("cannot update launch template version because the nodegroup is not configured to use one")
		}

		input.LaunchTemplate = &eks.LaunchTemplateSpecification{
			Version: &launchTemplateVersion,
		}

		if lt.Id != nil {
			input.LaunchTemplate.Id = describeNodegroupOutput.Nodegroup.LaunchTemplate.Id
		} else {
			input.LaunchTemplate.Name = describeNodegroupOutput.Nodegroup.LaunchTemplate.Name

		}
	}

	if version == "" {
		// Use the current Kubernetes version
		version, err := semver.ParseTolerant(*describeNodegroupOutput.Nodegroup.Version)
		if err != nil {
			return errors.Wrapf(err, "unexpected error parsing Kubernetes version %q", *describeNodegroupOutput.Nodegroup.Version)
		}
		input.Version = aws.String(fmt.Sprintf("%v.%v", version.Major, version.Minor))
	}

	upgradeResponse, err := m.ctl.Provider.EKS().UpdateNodegroupVersion(input)

	if err != nil {
		return err
	}

	if upgradeResponse != nil {
		logger.Debug("upgrade response for %q: %s", nodeGroupName, upgradeResponse.String())
	}

	logger.Info("upgrade of nodegroup %q in progress", nodeGroupName)

	newRequest := func() *request.Request {
		input := &eks.DescribeNodegroupInput{
			ClusterName:   &m.cfg.Metadata.Name,
			NodegroupName: &nodeGroupName,
		}
		req, _ := m.ctl.Provider.EKS().DescribeNodegroupRequest(input)
		return req
	}

	msg := fmt.Sprintf("waiting for upgrade of nodegroup %q to complete", nodeGroupName)

	acceptors := waiters.MakeAcceptors(
		"Nodegroup.Status",
		eks.NodegroupStatusActive,
		[]string{
			eks.NodegroupStatusDegraded,
		},
	)

	err = m.wait(nodeGroupName, msg, acceptors, newRequest, m.ctl.Provider.WaitTimeout(), nil)
	if err != nil {
		return err
	}
	logger.Info("nodegroup successfully upgraded")
	return nil
}
