package formatter

import (
	"github.com/eolinker/eosc/formatter/json"
	"github.com/eolinker/eosc/formatter/line"
	"sync"

	"github.com/eolinker/eosc"
)

var manager = NewManager()

type Manager struct {
	factory map[string]eosc.IFormatterFactory
	locker  sync.RWMutex
}

func init() {
	Register(line.Name, line.NewFactory())
	Register(json.Name, json.NewFactory())
}

func NewManager() *Manager {
	return &Manager{
		factory: make(map[string]eosc.IFormatterFactory),
		locker:  sync.RWMutex{},
	}
}

func (m *Manager) Get(name string) (eosc.IFormatterFactory, bool) {
	m.locker.RLock()

	defer m.locker.RUnlock()
	if v, ok := m.factory[name]; ok {
		return v, ok
	}
	return nil, false
}

func (m *Manager) Set(name string, factory eosc.IFormatterFactory) {
	m.locker.Lock()
	defer m.locker.Unlock()
	m.factory[name] = factory
}

func Register(name string, factory eosc.IFormatterFactory) {
	manager.Set(name, factory)
}

func GetFormatterFactory(name string) (eosc.IFormatterFactory, bool) {
	return manager.Get(name)
}
