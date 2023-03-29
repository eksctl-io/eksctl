package nodegroup

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/eks"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

func (m *Manager) Scale(ctx context.Context, ng *api.NodeGroupBase, wait bool) error {
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
		err = m.scaleUnmanagedNodeGroup(ctx, ng, stackInfo, wait)
	} else {
		err = m.scaleManagedNodeGroup(ctx, ng, wait)
	}

	if err != nil {
		return fmt.Errorf("failed to scale nodegroup %q for cluster %q, error: %v", ng.Name, m.cfg.Metadata.Name, err)
	}

	return nil
}

func (m *Manager) scaleUnmanagedNodeGroup(ctx context.Context, ng *api.NodeGroupBase, stackInfo manager.StackInfo, wait bool) error {
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

	if _, err := m.ctl.AWSProvider.ASG().UpdateAutoScalingGroup(ctx, input); err != nil {
		return err
	}
	logger.Info("initiated scaling of nodegroup")

	if wait {

		if counter, err := eks.GetNodes(m.clientSet, ng); err != nil {
			return errors.Wrap(err, "listing nodes")
		} else if counter >= *ng.MinSize {
			logger.Warning("when scaling down an ASG, passing the --wait flag currently has no effect")
		}

		timeoutCtx, cancel := context.WithTimeout(ctx, m.ctl.AWSProvider.WaitTimeout())
		defer cancel()
		if err := eks.WaitForNodes(timeoutCtx, m.clientSet, ng); err != nil {
			return err
		}

		logger.Info("nodegroup successfully scaled")
	} else {
		logger.Info("to see the status of the scaling run `eksctl get nodegroup --cluster %s --region %s --name %s`", m.cfg.Metadata.Name, m.ctl.AWSProvider.Region(), ng.Name)
	}
	return nil
}

func (m *Manager) scaleManagedNodeGroup(ctx context.Context, ng *api.NodeGroupBase, wait bool) error {

	input := &awseks.UpdateNodegroupConfigInput{
		ScalingConfig: &ekstypes.NodegroupScalingConfig{},
		ClusterName:   &m.cfg.Metadata.Name,
		NodegroupName: &ng.Name,
	}

	if ng.MaxSize != nil {
		input.ScalingConfig.MaxSize = aws.Int32(int32(*ng.MaxSize))
	}

	if ng.MinSize != nil {
		input.ScalingConfig.MinSize = aws.Int32(int32(*ng.MinSize))
	}

	if ng.DesiredCapacity != nil {
		input.ScalingConfig.DesiredSize = aws.Int32(int32(*ng.DesiredCapacity))
	}

	if _, err := m.ctl.AWSProvider.EKS().UpdateNodegroupConfig(ctx, input); err != nil {
		return err
	}
	logger.Info("initiated scaling of nodegroup")

	if wait {
		logger.Info("waiting for scaling of nodegroup %q to complete", ng.Name)

		waiter := awseks.NewNodegroupActiveWaiter(m.ctl.AWSProvider.EKS())
		if err := waiter.Wait(ctx, &awseks.DescribeNodegroupInput{
			ClusterName:   input.ClusterName,
			NodegroupName: input.NodegroupName,
		}, m.ctl.AWSProvider.WaitTimeout()); err != nil {
			return err
		}

		logger.Info("nodegroup successfully scaled")
	} else {
		logger.Info("to see the status of the scaling run `eksctl get nodegroup --cluster %s --region %s --name %s`", *input.ClusterName, m.ctl.AWSProvider.Region(), ng.Name)
	}
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
