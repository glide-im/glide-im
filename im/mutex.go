package im

import "sync"

type mutex struct {
	*sync.Mutex
}

func (m *mutex) LockUtilReturn() func() {
	m.Lock()
	return func() {
		m.Unlock()
	}
}
