package nodegroup

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/eks"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	kubewrapper "github.com/weaveworks/eksctl/pkg/kubernetes"
)

func (m *Manager) GetAll() ([]*manager.NodeGroupSummary, error) {
	summaries, err := m.stackManager.GetUnmanagedNodeGroupSummaries("")
	if err != nil {
		return nil, errors.Wrap(err, "getting nodegroup stack summaries")
	}

	for _, summary := range summaries {
		if summary.DesiredCapacity > 0 {
			summary.Version, err = kubewrapper.GetNodegroupKubernetesVersion(m.clientSet.CoreV1().Nodes(), summary.Name)
			if err != nil {
				return nil, errors.Wrap(err, "getting nodegroup's kubernetes version")
			}
		}
	}

	managedNodeGroups, err := m.ctl.Provider.EKS().ListNodegroups(&eks.ListNodegroupsInput{
		ClusterName: &m.cfg.Metadata.Name,
	})
	if err != nil {
		return nil, err
	}

	for _, managedNodeGroup := range managedNodeGroups.Nodegroups {
		describeOutput, err := m.ctl.Provider.EKS().DescribeNodegroup(&eks.DescribeNodegroupInput{
			ClusterName:   &m.cfg.Metadata.Name,
			NodegroupName: managedNodeGroup,
		})
		if err != nil {
			return nil, err
		}

		var stack *cloudformation.Stack
		stack, err = m.stackManager.DescribeNodeGroupStack(*managedNodeGroup)
		if err != nil {
			stack = &cloudformation.Stack{}
		}

		asgs := []string{}

		if describeOutput.Nodegroup.Resources != nil {
			for _, v := range describeOutput.Nodegroup.Resources.AutoScalingGroups {
				asgs = append(asgs, aws.StringValue(v.Name))
			}
		}

		summaries = append(summaries, &manager.NodeGroupSummary{
			StackName:            aws.StringValue(stack.StackName),
			Name:                 *describeOutput.Nodegroup.NodegroupName,
			Cluster:              *describeOutput.Nodegroup.ClusterName,
			Status:               *describeOutput.Nodegroup.Status,
			MaxSize:              int(*describeOutput.Nodegroup.ScalingConfig.MaxSize),
			MinSize:              int(*describeOutput.Nodegroup.ScalingConfig.MinSize),
			DesiredCapacity:      int(*describeOutput.Nodegroup.ScalingConfig.DesiredSize),
			InstanceType:         getInstanceTypes(describeOutput.Nodegroup),
			ImageID:              *describeOutput.Nodegroup.AmiType,
			CreationTime:         describeOutput.Nodegroup.CreatedAt,
			NodeInstanceRoleARN:  *describeOutput.Nodegroup.NodeRole,
			AutoScalingGroupName: strings.Join(asgs, ","),
			Version:              *describeOutput.Nodegroup.Version,
		})
	}

	return summaries, nil
}

func getInstanceTypes(ng *awseks.Nodegroup) string {
	if len(ng.InstanceTypes) > 0 {
		return strings.Join(aws.StringValueSlice(ng.InstanceTypes), ",")
	}
	logger.Info("no instance types reported by EKS for nodegroup %q", *ng.NodegroupName)
	return "-"
}

func (m *Manager) Get(name string) (*manager.NodeGroupSummary, error) {
	summaries, err := m.stackManager.GetUnmanagedNodeGroupSummaries(name)
	if err != nil {
		return nil, errors.Wrap(err, "getting nodegroup stack summaries")
	}

	if len(summaries) > 0 {
		s := summaries[0]
		if s.DesiredCapacity > 0 {
			s.Version, err = kubewrapper.GetNodegroupKubernetesVersion(m.clientSet.CoreV1().Nodes(), s.Name)
			if err != nil {
				return nil, errors.Wrap(err, "getting nodegroup's kubernetes version")
			}
		}
		return s, nil
	}

	describeOutput, err := m.ctl.Provider.EKS().DescribeNodegroup(&eks.DescribeNodegroupInput{
		ClusterName:   &m.cfg.Metadata.Name,
		NodegroupName: &name,
	})

	if err != nil {
		return nil, err
	}

	var asg string
	if describeOutput.Nodegroup.Resources != nil {
		for _, v := range describeOutput.Nodegroup.Resources.AutoScalingGroups {
			asg = aws.StringValue(v.Name)
		}
	}

	return &manager.NodeGroupSummary{
		Name:                 *describeOutput.Nodegroup.NodegroupName,
		Cluster:              *describeOutput.Nodegroup.ClusterName,
		Status:               *describeOutput.Nodegroup.Status,
		MaxSize:              int(*describeOutput.Nodegroup.ScalingConfig.MaxSize),
		MinSize:              int(*describeOutput.Nodegroup.ScalingConfig.MinSize),
		DesiredCapacity:      int(*describeOutput.Nodegroup.ScalingConfig.DesiredSize),
		InstanceType:         *describeOutput.Nodegroup.InstanceTypes[0],
		ImageID:              *describeOutput.Nodegroup.AmiType,
		CreationTime:         describeOutput.Nodegroup.CreatedAt,
		NodeInstanceRoleARN:  *describeOutput.Nodegroup.NodeRole,
		AutoScalingGroupName: asg,
		Version:              *describeOutput.Nodegroup.Version,
	}, nil
}
