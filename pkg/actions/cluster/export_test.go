package cluster

import "github.com/weaveworks/eksctl/pkg/kubernetes"

func (c *UnownedCluster) SetNewClientSet(newClientSet func() (kubernetes.Interface, error)) {
	c.newClientSet = newClientSet
}

func (c *OwnedCluster) SetNewClientSet(newClientSet func() (kubernetes.Interface, error)) {
	c.newClientSet = newClientSet
}
