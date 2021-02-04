package nodegroup

func (m *Manager) SetWaiter(wait WaitFunc) {
	m.wait = wait
}
