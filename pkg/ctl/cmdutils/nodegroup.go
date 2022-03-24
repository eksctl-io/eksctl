package cmdutils

import (
	"context"

	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

func PopulateNodegroup(ctx context.Context, stackManager manager.StackManager, name string, cfg *api.ClusterConfig, ctl api.ClusterProvider) error {
	nodeGroupType, err := stackManager.GetNodeGroupStackType(ctx, manager.GetNodegroupOption{
		NodeGroupName: name,
	})
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

	return PopulateNodegroupFromStack(nodeGroupType, name, cfg)
}

// PopulateNodegroupFromStack populates the nodegroup field of an api.ClusterConfig by type from its CF stack.
func PopulateNodegroupFromStack(nodeGroupType api.NodeGroupType, nodeGroupName string, cfg *api.ClusterConfig) error {
	switch nodeGroupType {
	case api.NodeGroupTypeUnmanaged:
		PopulateUnmanagedNodegroup(nodeGroupName, cfg)
	case api.NodeGroupTypeManaged:
		cfg.ManagedNodeGroups = append(cfg.ManagedNodeGroups, &api.ManagedNodeGroup{
			NodeGroupBase: &api.NodeGroupBase{
				Name: nodeGroupName,
			},
		})
	case api.NodeGroupTypeUnowned:
		cfg.ManagedNodeGroups = append(cfg.ManagedNodeGroups, &api.ManagedNodeGroup{
			NodeGroupBase: &api.NodeGroupBase{
				Name: nodeGroupName,
			},
			Unowned: true,
		})
	}
	return nil
}

// PopulateUnmanagedNodegroup populates the unmanaged nodegroup field of a ClucterConfig.
func PopulateUnmanagedNodegroup(nodeGroupName string, cfg *api.ClusterConfig) {
	cfg.NodeGroups = append(cfg.NodeGroups, &api.NodeGroup{
		NodeGroupBase: &api.NodeGroupBase{
			Name: nodeGroupName,
		},
	})
}
