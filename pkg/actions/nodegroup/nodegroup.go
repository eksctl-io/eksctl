package nodegroup

import (
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/eks"
)

type NodeGroup struct {
	manager *manager.StackCollection
	ctl     *eks.ClusterProvider
	cfg     *api.ClusterConfig
}

func New(cfg *api.ClusterConfig, ctl *eks.ClusterProvider) NodeGroup {
	return NodeGroup{
		manager: ctl.NewStackManager(cfg),
		ctl:     ctl,
		cfg:     cfg,
	}
}
