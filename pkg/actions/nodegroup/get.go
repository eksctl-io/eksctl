package nodegroup

import (
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

func (ng *NodeGroup) GetAll() ([]*manager.NodeGroupSummary, error) {
	summaries, err := ng.manager.GetNodeGroupSummaries("")
	if err != nil {
		return nil, errors.Wrap(err, "getting nodegroup stack summaries")
	}

	nodeGroups, err := ng.ctl.Provider.EKS().ListNodegroups(&eks.ListNodegroupsInput{
		ClusterName: &ng.cfg.Metadata.Name,
	})

	if err != nil {
		return nil, err
	}

	var unownedNodeGroups []string
	for _, ng := range nodeGroups.Nodegroups {
		found := false
		for _, summary := range summaries {
			if summary.Name == *ng {
				found = true
			}
		}

		if !found {
			unownedNodeGroups = append(unownedNodeGroups, *ng)
		}
	}

	for _, unownedNodeGroup := range unownedNodeGroups {
		describeOutput, err := ng.ctl.Provider.EKS().DescribeNodegroup(&eks.DescribeNodegroupInput{
			ClusterName:   &ng.cfg.Metadata.Name,
			NodegroupName: &unownedNodeGroup,
		})
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, &manager.NodeGroupSummary{
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
		})
	}

	return summaries, nil
}

func (ng *NodeGroup) Get(name string) (*manager.NodeGroupSummary, error) {
	summaries, err := ng.manager.GetNodeGroupSummaries(name)
	if err != nil {
		return nil, errors.Wrap(err, "getting nodegroup stack summaries")
	}

	if len(summaries) > 0 {
		return summaries[0], nil
	}

	describeOutput, err := ng.ctl.Provider.EKS().DescribeNodegroup(&eks.DescribeNodegroupInput{
		ClusterName:   &ng.cfg.Metadata.Name,
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
