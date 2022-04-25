package nodegroup

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/kris-nova/logger"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/utils/waiters"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/eks"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

func (m *Manager) Scale(ctx context.Context, ng *api.NodeGroupBase) error {
	logger.Info("scaling nodegroup %q in cluster %s", ng.Name, m.cfg.Metadata.Name)

	nodegroupStackInfos, err := m.stackManager.DescribeNodeGroupStacksAndResources(ctx)
	if err != nil {
		return err
	}

	var stackInfo manager.StackInfo
	var ok, isUnmanagedNodegroup bool
	stackInfo, ok = nodegroupStackInfos[ng.Name]
	if ok {
		nodegroupType, err := manager.GetNodeGroupType(stackInfo.Stack.Tags)
		if err != nil {
			return err
		}
		isUnmanagedNodegroup = nodegroupType == api.NodeGroupTypeUnmanaged
	}

	if isUnmanagedNodegroup {
		err = m.scaleUnmanagedNodeGroup(ctx, ng, stackInfo)
	} else {
		err = m.scaleManagedNodeGroup(ng)
	}

	if err != nil {
		return fmt.Errorf("failed to scale nodegroup for cluster %q, error: %v", m.cfg.Metadata.Name, err)
	}

	return nil
}

func (m *Manager) scaleUnmanagedNodeGroup(ctx context.Context, ng *api.NodeGroupBase, stackInfo manager.StackInfo) error {
	asgName := ""
	for _, resource := range stackInfo.Resources {
		if *resource.LogicalResourceId == "NodeGroup" {
			asgName = *resource.PhysicalResourceId
			break
		}
	}

	if asgName == "" {
		return fmt.Errorf("failed to find NodeGroup auto scaling group")
	}

	input := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: &asgName,
	}

	if ng.MaxSize != nil {
		input.MaxSize = aws.Int32(int32(*ng.MaxSize))
	}

	if ng.MinSize != nil {
		input.MinSize = aws.Int32(int32(*ng.MinSize))
	}

	if ng.DesiredCapacity != nil {
		input.DesiredCapacity = aws.Int32(int32(*ng.DesiredCapacity))
	}

	_, err := m.ctl.Provider.ASG().UpdateAutoScalingGroup(ctx, input)
	if err != nil {
		return err
	}
	logger.Info("nodegroup successfully scaled")

	return nil
}

func (m *Manager) scaleManagedNodeGroup(ng *api.NodeGroupBase) error {
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
