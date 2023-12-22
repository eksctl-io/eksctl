package cluster

import (
	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

var (
	DrainAllNodeGroups = drainAllNodeGroups
)

func (c *UnownedCluster) SetNewClientSet(newClientSet func() (kubernetes.Interface, error)) {
	c.newClientSet = newClientSet
}

func (c *OwnedCluster) SetNewClientSet(newClientSet func() (kubernetes.Interface, error)) {
	c.newClientSet = newClientSet
}

func (c *UnownedCluster) SetNewNodeGroupDrainer(newNodeGroupDrainer func(clientSet kubernetes.Interface) NodeGroupDrainer) {
	c.newNodeGroupDrainer = newNodeGroupDrainer
}

func (c *OwnedCluster) SetNewNodeGroupDrainer(newNodeGroupDrainer func(kubernetes.Interface) NodeGroupDrainer) {
	c.newNodeGroupDrainer = newNodeGroupDrainer
}

func SetProviderConstructor(f ProviderConstructor) {
	newClusterProvider = f
}

func SetStackManagerConstructor(f StackManagerConstructor) {
	newStackCollection = f
}
