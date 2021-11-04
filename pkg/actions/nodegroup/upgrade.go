package nodegroup

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/blang/semver"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/managed"
	"github.com/weaveworks/eksctl/pkg/utils/waiters"
)

func (m *Manager) Upgrade(options managed.UpgradeOptions) error {
	stackCollection := manager.NewStackCollection(m.ctl.Provider, m.cfg)
	hasStacks, err := m.hasStacks(options.NodegroupName)
	if err != nil {
		return err
	}

	if options.KubernetesVersion != "" {
		if _, err := semver.ParseTolerant(options.KubernetesVersion); err != nil {
			return errors.Wrap(err, "invalid Kubernetes version")
		}
	}

	if hasStacks {
		managedService := managed.NewService(m.ctl.Provider.EKS(), m.ctl.Provider.SSM(), m.ctl.Provider.EC2(), stackCollection, m.cfg.Metadata.Name)
		return managedService.UpgradeNodeGroup(options)
	}

	if err := m.upgrade(options); err != nil {
		return err
	}

	if options.Wait {
		return m.waitForUpgrade(options)
	}

	logger.Info("nodegroup upgrade request submitted successfully")

	return nil

}

func (m *Manager) upgrade(options managed.UpgradeOptions) error {
	input := &eks.UpdateNodegroupVersionInput{
		ClusterName:   &m.cfg.Metadata.Name,
		Force:         &options.ForceUpgrade,
		NodegroupName: &options.NodegroupName,
		Version:       &options.KubernetesVersion,
	}

	describeNodegroupOutput, err := m.ctl.Provider.EKS().DescribeNodegroup(&eks.DescribeNodegroupInput{
		ClusterName:   &m.cfg.Metadata.Name,
		NodegroupName: &options.NodegroupName,
	})

	if err != nil {
		return err
	}

	if options.LaunchTemplateVersion != "" {
		lt := describeNodegroupOutput.Nodegroup.LaunchTemplate
		if lt == nil || (lt.Id == nil && lt.Name == nil) {
			return errors.New("cannot update launch template version because the nodegroup is not configured to use one")
		}

		input.LaunchTemplate = &eks.LaunchTemplateSpecification{
			Version: &options.LaunchTemplateVersion,
		}

		if lt.Id != nil {
			input.LaunchTemplate.Id = describeNodegroupOutput.Nodegroup.LaunchTemplate.Id
		} else {
			input.LaunchTemplate.Name = describeNodegroupOutput.Nodegroup.LaunchTemplate.Name

		}
	}

	if options.KubernetesVersion == "" {
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
		logger.Debug("upgrade response for %q: %s", options.NodegroupName, upgradeResponse.String())
	}

	logger.Info("upgrade of nodegroup %q in progress", options.NodegroupName)
	return nil
}

func (m *Manager) waitForUpgrade(options managed.UpgradeOptions) error {

	newRequest := func() *request.Request {
		input := &eks.DescribeNodegroupInput{
			ClusterName:   &m.cfg.Metadata.Name,
			NodegroupName: &options.NodegroupName,
		}
		req, _ := m.ctl.Provider.EKS().DescribeNodegroupRequest(input)
		return req
	}

	msg := fmt.Sprintf("waiting for upgrade of nodegroup %q to complete", options.NodegroupName)

	acceptors := waiters.MakeAcceptors(
		"Nodegroup.Status",
		eks.NodegroupStatusActive,
		[]string{
			eks.NodegroupStatusDegraded,
		},
	)

	err := m.wait(options.NodegroupName, msg, acceptors, newRequest, m.ctl.Provider.WaitTimeout(), nil)
	if err != nil {
		return err
	}
	logger.Info("nodegroup successfully upgraded")
	return nil
}
