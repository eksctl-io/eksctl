package eks

import (
	"context"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

type NodeGroupStackLister interface {
	ListNodeGroupStacksWithStatuses(ctx context.Context) ([]manager.NodeGroupStack, error)
}

type NodeGroupLister struct {
	NodeGroupStackLister NodeGroupStackLister
}

func (n *NodeGroupLister) List(ctx context.Context) ([]KubeNodeGroup, error) {
	nodeGroupStacks, err := n.NodeGroupStackLister.ListNodeGroupStacksWithStatuses(ctx)
	if err != nil {
		return nil, err
	}
	var kubeNodeGroups []KubeNodeGroup
	for _, stack := range nodeGroupStacks {
		ngBase := &api.NodeGroupBase{
			Name: stack.NodeGroupName,
		}
		switch stack.Type {
		case api.NodeGroupTypeManaged:
			kubeNodeGroups = append(kubeNodeGroups, &api.NodeGroup{
				NodeGroupBase: ngBase,
			})
		case api.NodeGroupTypeUnmanaged:
			kubeNodeGroups = append(kubeNodeGroups, &api.ManagedNodeGroup{
				NodeGroupBase: ngBase,
			})
		}
	}
	return kubeNodeGroups, nil
}
