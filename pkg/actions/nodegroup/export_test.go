package nodegroup

import "github.com/weaveworks/eksctl/pkg/cfn/manager"

func (m *Manager) SetWaiter(wait WaitFunc) {
	m.wait = wait
}

func (m *Manager) SetStackManager(stackManager manager.StackManager) {
	m.stackManager = stackManager
}
