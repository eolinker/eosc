package variable

import (
	"encoding/json"
	"fmt"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	"strings"
)

type IVariable interface {
	SetByNamespace(namespace string, variables map[string]string) (map[string]string, []string, error)
	GetByNamespace(namespace string) (map[string]string, bool)
	SetVariablesById(id string, variables []string)
	GetVariablesById(id string) []string
	GetIdsByVariable(variable string) []string
	GetAll() map[string]string
	Namespaces() []string
}

type Variables struct {
	// variables 变量数据
	variables      eosc.IUntyped
	requireManager IRequires
}

func NewVariables(data map[string][]byte) IVariable {
	v := &Variables{variables: eosc.NewUntyped(), requireManager: NewRequireManager()}
	for namespace, value := range data {
		var variables map[string]string
		json.Unmarshal(value, &variables)
		v.SetByNamespace(namespace, variables)
	}
	return v
}

func (m *Variables) SetVariablesById(id string, variables []string) {
	m.requireManager.Set(id, variables)
}

func (m *Variables) GetVariablesById(id string) []string {
	return m.requireManager.WorkerIDs(id)
}

func (m *Variables) GetIdsByVariable(variable string) []string {
	return m.requireManager.RequireIDs(variable)
}

func (m *Variables) SetByNamespace(namespace string, variables map[string]string) (map[string]string, []string, error) {
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

func (m *Variables) getAll() map[string]string {
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

func (m *Variables) getByNamespace(namespace string) (map[string]string, bool) {
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

func (m *Variables) GetByNamespace(namespace string) (map[string]string, bool) {
	return m.getByNamespace(namespace)
}

func (m *Variables) GetAll() map[string]string {
	return m.getAll()
}

func (m *Variables) Namespaces() []string {
	return m.variables.Keys()
}

func TrimNamespace(origin map[string]string) map[string]string {
	target := make(map[string]string)
	for key, value := range origin {
		index := strings.Index(key, "@")
		if index < 0 {
			continue
		}
		key = key[:index]

		target[key] = value
	}
	return target
}

func FillNamespace(namespace string, origin map[string]string) map[string]string {
	target := make(map[string]string)
	for key, value := range origin {
		target[fmt.Sprintf("%s@%s", key, namespace)] = value
	}
	return target
}
