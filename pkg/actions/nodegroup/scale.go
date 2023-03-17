package nodegroup

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
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
		return fmt.Errorf("failed to scale nodegroup %q for cluster %q, error: %v", ng.Name, m.cfg.Metadata.Name, err)
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

	if err := validateNodeGroupAMI(ctx, m.ctl.AWSProvider, asgName); err != nil {
		return err
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

func validateNodeGroupAMI(ctx context.Context, awsProvider api.ClusterProvider, asgName string) error {
	asg, err := awsProvider.ASG().DescribeAutoScalingGroups(ctx, &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{asgName},
	})
	if err != nil {
		return fmt.Errorf("error describing Auto Scaling group %q for nodegroup: %w", asgName, err)
	}
	if len(asg.AutoScalingGroups) != 1 {
		return fmt.Errorf("expected to find exactly one Auto Scaling group for nodegroup; got %d", len(asg.AutoScalingGroups))
	}
	lt := asg.AutoScalingGroups[0].LaunchTemplate
	if lt == nil {
		logger.Warning("nodegroup with Auto Scaling group %q does not have a launch template", asgName)
		return nil
	}

	ltData, err := awsProvider.EC2().DescribeLaunchTemplateVersions(ctx, &ec2.DescribeLaunchTemplateVersionsInput{
		LaunchTemplateId: lt.LaunchTemplateId,
		Versions:         []string{aws.ToString(lt.Version)},
	})
	if err != nil {
		return fmt.Errorf("error describing launch template %q for Auto Scaling group %q: %w", aws.ToString(lt.LaunchTemplateId), asgName, err)
	}
	if len(ltData.LaunchTemplateVersions) != 1 {
		return fmt.Errorf("expected to find exactly one launch template %q with version %q for Auto Scaling group %q; got %d", aws.ToString(lt.LaunchTemplateId), aws.ToString(lt.Version), asgName, len(ltData.LaunchTemplateVersions))
	}
	imageID := ltData.LaunchTemplateVersions[0].LaunchTemplateData.ImageId
	if imageID == nil {
		logger.Warning("nodegroup with launch template %q does not have an AMI", aws.ToString(lt.LaunchTemplateId))
		return nil
	}

	describeImagesOutput, err := awsProvider.EC2().DescribeImages(ctx, &ec2.DescribeImagesInput{
		ImageIds: []string{aws.ToString(imageID)},
	})
	if err != nil {
		return fmt.Errorf("error describing AMI for launch template %q: %w", aws.ToString(lt.LaunchTemplateId), err)
	}
	if len(describeImagesOutput.Images) == 0 {
		return errors.New("AMI associated with the nodegroup is either deprecated or removed; please upgrade the nodegroup before scaling it: https://eksctl.io/usage/nodegroup-upgrade/")
	}
	return nil
}
