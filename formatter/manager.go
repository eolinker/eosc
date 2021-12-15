package formatter

import "sync"

var manager = NewManager()

type Manager struct {
	factory map[string]IFormatterFactory
	locker  sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		factory: make(map[string]IFormatterFactory),
		locker:  sync.RWMutex{},
	}
}

func (m *Manager) Get(name string) (IFormatterFactory, bool) {
	m.locker.RLock()

	defer m.locker.RUnlock()
	if v, ok := m.factory[name]; ok {
		return v, ok
	}
	return nil, false
}

func (m *Manager) Set(name string, factory IFormatterFactory) {
	m.locker.Lock()
	defer m.locker.Unlock()
	m.factory[name] = factory
}

func Register(name string, factory IFormatterFactory) {
	manager.Set(name, factory)
}

func GetFormatterFactory(name string) (IFormatterFactory, bool) {
	return manager.Get(name)
}
