package nodegroup

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/kris-nova/logger"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"

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
		err = m.scaleManagedNodeGroup(ctx, ng)
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

	_, err := m.ctl.AWSProvider.ASG().UpdateAutoScalingGroup(ctx, input)
	if err != nil {
		return err
	}
	logger.Info("nodegroup successfully scaled")

	return nil
}

func (m *Manager) scaleManagedNodeGroup(ctx context.Context, ng *api.NodeGroupBase) error {
	scalingConfig := &ekstypes.NodegroupScalingConfig{}

	if ng.MaxSize != nil {
		scalingConfig.MaxSize = aws.Int32(int32(*ng.MaxSize))
	}

	if ng.MinSize != nil {
		scalingConfig.MinSize = aws.Int32(int32(*ng.MinSize))
	}

	if ng.DesiredCapacity != nil {
		scalingConfig.DesiredSize = aws.Int32(int32(*ng.DesiredCapacity))
	}

	_, err := m.ctl.AWSProvider.EKS().UpdateNodegroupConfig(ctx, &eks.UpdateNodegroupConfigInput{
		ScalingConfig: scalingConfig,
		ClusterName:   &m.cfg.Metadata.Name,
		NodegroupName: &ng.Name,
	})

	if err != nil {
		return err
	}

	logger.Info("waiting for scaling of nodegroup %q to complete", ng.Name)

	waiter := eks.NewNodegroupActiveWaiter(m.ctl.AWSProvider.EKS())
	if err := waiter.Wait(ctx, &eks.DescribeNodegroupInput{
		ClusterName:   &m.cfg.Metadata.Name,
		NodegroupName: &ng.Name,
	}, m.ctl.AWSProvider.WaitTimeout()); err != nil {
		return err
	}

	logger.Info("nodegroup successfully scaled")
	return nil
}
