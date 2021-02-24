package cmdutils

import (
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/kris-nova/logger"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

func PopulateNodegroup(stackManager manager.StackManager, name string, cfg *api.ClusterConfig, ctl api.ClusterProvider) error {
	nodeGroupType, err := stackManager.GetNodeGroupStackType(name)
	if err != nil {
		logger.Debug("failed to fetch nodegroup %q stack: %v", name, err)
		_, err := ctl.EKS().DescribeNodegroup(&eks.DescribeNodegroupInput{
			ClusterName:   &cfg.Metadata.Name,
			NodegroupName: &name,
		})

		if err != nil {
			return err
		}
		nodeGroupType = api.NodeGroupTypeUnowned
	}
	switch nodeGroupType {
	case api.NodeGroupTypeUnmanaged:
		cfg.NodeGroups = []*api.NodeGroup{
			{
				NodeGroupBase: &api.NodeGroupBase{
					Name: name,
				},
			},
		}
	case api.NodeGroupTypeManaged:
		cfg.ManagedNodeGroups = []*api.ManagedNodeGroup{
			{
				NodeGroupBase: &api.NodeGroupBase{
					Name: name,
				},
			},
		}
	case api.NodeGroupTypeUnowned:
		cfg.ManagedNodeGroups = []*api.ManagedNodeGroup{
			{
				NodeGroupBase: &api.NodeGroupBase{
					Name: name,
				},
				Unowned: true,
			},
		}
	}

	return nil
}
