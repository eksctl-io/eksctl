package addon

import "time"

func (m *Manager) SetTimeout(timeout time.Duration) {
	m.timeout = timeout
}
