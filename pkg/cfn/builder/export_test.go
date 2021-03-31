package builder

import "github.com/weaveworks/eksctl/pkg/nodebootstrap"

func (n *NodeGroupResourceSet) SetBootstrapper(bootstrapper nodebootstrap.Bootstrapper) {
	n.bootstrapper = bootstrapper
}
