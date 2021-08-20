package comm

import "sync"

type Mutex struct {
	sync.Mutex
}

func NewMutex() *Mutex {
	return &Mutex{sync.Mutex{}}
}

func (m *Mutex) LockUtilReturn() func() {
	m.Lock()
	return func() {
		m.Unlock()
	}
}
