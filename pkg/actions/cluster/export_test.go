package cluster

import (
	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

func (c *UnownedCluster) SetNewClientSet(newClientSet func() (kubernetes.Interface, error)) {
	c.newClientSet = newClientSet
}

func (c *OwnedCluster) SetNewClientSet(newClientSet func() (kubernetes.Interface, error)) {
	c.newClientSet = newClientSet
}

func (c *UnownedCluster) SetNewNodeGroupManager(newNodeGroupManager func(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, clientSet kubernetes.Interface) *nodegroup.Manager) {
	c.newNodeGroupManager = newNodeGroupManager
}

func (c *OwnedCluster) SetNewNodeGroupManager(newNodeGroupManager func(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, clientSet kubernetes.Interface) *nodegroup.Manager) {
	c.newNodeGroupManager = newNodeGroupManager
}
