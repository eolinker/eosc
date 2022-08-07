package variable

import (
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/workers/require"
	"sync"
)

type IVariable interface {
	SetByNamespace(namespace string, variables map[string]string) (map[string]string, []string)
	GetByNamespace(namespace string) (map[string]string, bool)
	Get() map[string]string
}

type Manager struct {
	// variables 变量数据
	variables      eosc.IUntyped
	requireManager require.IRequires
	locker         sync.RWMutex
}

func (m *Manager) SetByNamespace(namespace string, variables map[string]string) (map[string]string, []string) {
	// variables的key为：{变量名}@{namespace}，如：v1@default
	old, has := m.getByNamespace(namespace)
	if !has {
		// 此时变量都是新的，没有受影响的配置id
		m.variables.Set(namespace, variables)
		newVariables := make(map[string]string)
		for _, key := range m.variables.Keys() {
			vs, ok := m.getByNamespace(key)
			if !ok {
				log.Error("fail to get variable,namespace is ", namespace)
				continue
			}
			for k, v := range vs {
				newVariables[k] = v
			}
		}
		return newVariables, nil
	}

	return newVariables, keys
}

func (m *Manager) getByNamespace(namespace string) (map[string]string, bool) {
	variables, has := m.variables.Get(namespace)
	if !has {
		return nil, false
	}
	v, ok := variables.(map[string]string)
	if !ok {
		return nil, false
	}
	return v, ok
}

func (m *Manager) GetByNamespace(namespace string) (map[string]string, bool) {
	//TODO implement me
	panic("implement me")
}

func (m *Manager) Get() map[string]string {
	//TODO implement me
	panic("implement me")
}
