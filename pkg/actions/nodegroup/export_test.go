package nodegroup

import (
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/eks"
)

func (m *Manager) SetWaiter(wait WaitFunc) {
	m.wait = wait
}

func (m *Manager) SetStackManager(stackManager manager.StackManager) {
	m.stackManager = stackManager
}

// MockKubeProvider can be used for passing a mock of the kube provider.
func (m *Manager) MockKubeProvider(k eks.KubeProvider) {
	m.kubeProvider = k
}

// MockNodeGroupService can be used for passing a mock of the nodegroup initialiser.
func (m *Manager) MockNodeGroupService(ngSvc eks.NodeGroupInitialiser) {
	m.init = ngSvc
}
