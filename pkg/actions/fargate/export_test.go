package fargate

import "k8s.io/client-go/kubernetes"

func (m *Manager) SetNewClientSet(newClientSet func() (kubernetes.Interface, error)) {
	m.newStdClientSet = newClientSet
}
