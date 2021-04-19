package nodegroup

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

func (m *Manager) GetAll() ([]*manager.NodeGroupSummary, error) {
	summaries, err := m.stackManager.GetUnmanagedNodeGroupSummaries("")
	if err != nil {
		return nil, errors.Wrap(err, "getting nodegroup stack summaries")
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

		asgs := []string{}

		if describeOutput.Nodegroup.Resources != nil {
			for _, v := range describeOutput.Nodegroup.Resources.AutoScalingGroups {
				asgs = append(asgs, aws.StringValue(v.Name))
			}
		}

		summaries = append(summaries, &manager.NodeGroupSummary{
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
		return summaries[0], nil
	}

	describeOutput, err := m.ctl.Provider.EKS().DescribeNodegroup(&eks.DescribeNodegroupInput{
		ClusterName:   &m.cfg.Metadata.Name,
		NodegroupName: &name,
	})

	if err != nil {
		return nil, err
	}

	return &manager.NodeGroupSummary{
		Name:                *describeOutput.Nodegroup.NodegroupName,
		Cluster:             *describeOutput.Nodegroup.ClusterName,
		Status:              *describeOutput.Nodegroup.Status,
		MaxSize:             int(*describeOutput.Nodegroup.ScalingConfig.MaxSize),
		MinSize:             int(*describeOutput.Nodegroup.ScalingConfig.MinSize),
		DesiredCapacity:     int(*describeOutput.Nodegroup.ScalingConfig.DesiredSize),
		InstanceType:        *describeOutput.Nodegroup.InstanceTypes[0],
		ImageID:             *describeOutput.Nodegroup.AmiType,
		CreationTime:        describeOutput.Nodegroup.CreatedAt,
		NodeInstanceRoleARN: *describeOutput.Nodegroup.NodeRole,
	}, nil
}
