package nodegroup

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

func (m *Manager) GetAll() ([]*manager.NodeGroupSummary, error) {
	summaries, err := m.stackManager.GetNodeGroupSummaries("")
	if err != nil {
		return nil, errors.Wrap(err, "getting nodegroup stack summaries")
	}

	nodeGroups, err := m.ctl.Provider.EKS().ListNodegroups(&eks.ListNodegroupsInput{
		ClusterName: &m.cfg.Metadata.Name,
	})
	if err != nil {
		return nil, err
	}

	var nodeGroupsWithoutStacks []string
	for _, ng := range nodeGroups.Nodegroups {
		found := false
		for _, summary := range summaries {
			if summary.Name == *ng {
				found = true
			}
		}

		if !found {
			nodeGroupsWithoutStacks = append(nodeGroupsWithoutStacks, *ng)
		}
	}

	for _, nodeGroupWithoutStack := range nodeGroupsWithoutStacks {
		describeOutput, err := m.ctl.Provider.EKS().DescribeNodegroup(&eks.DescribeNodegroupInput{
			ClusterName:   &m.cfg.Metadata.Name,
			NodegroupName: &nodeGroupWithoutStack,
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
			InstanceType:         *describeOutput.Nodegroup.InstanceTypes[0],
			ImageID:              *describeOutput.Nodegroup.AmiType,
			CreationTime:         describeOutput.Nodegroup.CreatedAt,
			NodeInstanceRoleARN:  *describeOutput.Nodegroup.NodeRole,
			AutoScalingGroupName: strings.Join(asgs, ","),
		})
	}

	return summaries, nil
}

func (m *Manager) Get(name string) (*manager.NodeGroupSummary, error) {
	summaries, err := m.stackManager.GetNodeGroupSummaries(name)
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
