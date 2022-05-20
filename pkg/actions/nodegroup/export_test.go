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

func (m *Manager) GetStackManager() manager.StackManager {
	return m.stackManager
}

// MockKubeProvider can be used for passing a mock of the kube provider.
func (m *Manager) MockKubeProvider(k eks.KubeProvider) {
	m.ctl.KubeProvider = k
}
