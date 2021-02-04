package nodegroup

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/kris-nova/logger"
	"github.com/weaveworks/eksctl/pkg/utils/waiters"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/eks"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

func (m *Manager) Scale(ng *api.NodeGroup) error {
	logger.Info("scaling nodegroup %q in cluster %s", ng.Name, m.cfg.Metadata.Name)
	stackManager := m.ctl.NewStackManager(m.cfg)

	hasStacks, err := m.hasStacks(ng.Name)
	if err != nil {
		return err
	}

	if hasStacks {
		err = stackManager.ScaleNodeGroup(ng)
	} else {
		err = m.scale(ng)
	}

	if err != nil {
		return fmt.Errorf("failed to scale nodegroup for cluster %q, error: %v", m.cfg.Metadata.Name, err)
	}

	return nil
}

func (m *Manager) scale(ng *api.NodeGroup) error {
	scalingConfig := &eks.NodegroupScalingConfig{}

	if ng.MaxSize != nil {
		scalingConfig.MaxSize = aws.Int64(int64(*ng.MaxSize))
	}

	if ng.MinSize != nil {
		scalingConfig.MinSize = aws.Int64(int64(*ng.MinSize))
	}

	if ng.DesiredCapacity != nil {
		scalingConfig.DesiredSize = aws.Int64(int64(*ng.DesiredCapacity))
	}

	_, err := m.ctl.Provider.EKS().UpdateNodegroupConfig(&eks.UpdateNodegroupConfigInput{
		ScalingConfig: scalingConfig,
		ClusterName:   &m.cfg.Metadata.Name,
		NodegroupName: &ng.Name,
	})

	if err != nil {
		return err
	}

	newRequest := func() *request.Request {
		input := &eks.DescribeNodegroupInput{
			ClusterName:   &m.cfg.Metadata.Name,
			NodegroupName: &ng.Name,
		}
		req, _ := m.ctl.Provider.EKS().DescribeNodegroupRequest(input)
		return req
	}

	msg := fmt.Sprintf("waiting for scaling of nodegroup %q to complete", ng.Name)

	acceptors := waiters.MakeAcceptors(
		"Nodegroup.Status",
		eks.NodegroupStatusActive,
		[]string{
			eks.NodegroupStatusDegraded,
		},
	)

	err = m.wait(ng.Name, msg, acceptors, newRequest, m.ctl.Provider.WaitTimeout(), nil)
	if err != nil {
		return err
	}
	logger.Info("nodegroup successfully scaled")
	return nil
}
