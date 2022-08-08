package variable

import (
	"fmt"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
)

type IVariable interface {
	SetByNamespace(namespace string, variables map[string]string) (map[string]string, []string, error)
	GetByNamespace(namespace string) (map[string]string, bool)
	SetVariablesById(id string, variables []string)
	GetVariablesById(id string) []string
	GetIdsByVariable(variable string) []string
	GetAll() map[string]string
}

type Manager struct {
	// variables 变量数据
	variables      eosc.IUntyped
	requireManager IRequires
}

func (m *Manager) SetVariablesById(id string, variables []string) {
	m.requireManager.Set(id, variables)
}

func (m *Manager) GetVariablesById(id string) []string {
	return m.requireManager.WorkerIDs(id)
}

func (m *Manager) GetIdsByVariable(variable string) []string {
	return m.requireManager.RequireIDs(variable)
}

func (m *Manager) SetByNamespace(namespace string, variables map[string]string) (map[string]string, []string, error) {
	// variables的key为：{变量名}@{namespace}，如：v1@default
	old, has := m.getByNamespace(namespace)
	if !has {
		m.variables.Set(namespace, variables)
		// 此时变量都是新的，没有受影响的配置id
		return m.getAll(), nil, nil
	}
	affectIds := make([]string, 0, len(variables))
	for key, value := range variables {
		if v, ok := old[key]; ok {
			if v != value {
				// 将更新的key记录下来
				affectIds = append(affectIds, m.requireManager.RequireIDs(key)...)
			}
			delete(old, key)
			continue
		}
		// 将新增的key记录下来
		affectIds = append(affectIds, m.requireManager.RequireIDs(key)...)
	}
	for key := range old {
		if m.requireManager.RequireByCount(key) > 0 {
			return nil, nil, fmt.Errorf("variable %s %w", key, eosc.ErrorRequire)
		}
	}
	m.variables.Set(namespace, variables)
	return m.getAll(), affectIds, nil
}

func (m *Manager) getAll() map[string]string {
	newVariables := make(map[string]string)
	for _, key := range m.variables.Keys() {
		vs, ok := m.getByNamespace(key)
		if !ok {
			log.Error("fail to get variable,namespace is ", key)
			continue
		}
		for k, v := range vs {
			newVariables[k] = v
		}
	}
	return newVariables
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
	return m.getByNamespace(namespace)
}

func (m *Manager) GetAll() map[string]string {
	return m.getAll()
}
