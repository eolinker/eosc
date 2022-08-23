package variable

import (
	"fmt"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	"reflect"
	"strings"
)

var _ IVariable = (*Variables)(nil)

type IVariable interface {
	SetByNamespace(namespace string, variables map[string]string) error
	GetByNamespace(namespace string) (map[string]string, bool)
	SetVariablesById(id string, variables []string)
	Namespaces() []string
	Unmarshal(buf []byte, typ reflect.Type) (interface{}, []string, error)
	Check(namespace string, variables map[string]string) ([]string, error)
}

type Variables struct {
	// variables 变量数据
	variables      eosc.IUntyped
	requireManager IRequires
}

func (m *Variables) Unmarshal(buf []byte, typ reflect.Type) (interface{}, []string, error) {
	return NewParse(m.getAll()).Unmarshal(buf, typ)
}

func NewVariables(data map[string][]byte) IVariable {
	v := &Variables{variables: eosc.NewUntyped(), requireManager: NewRequireManager()}
	for namespace, value := range data {
		v.variables.Set(namespace, value)
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

func (m *Variables) Check(namespace string, variables map[string]string) ([]string, error) {
	// variables的key为：{变量名}@{namespace}，如：v1@default
	old, has := m.getByNamespace(namespace)
	if !has {
		m.variables.Set(namespace, variables)
		// 此时变量都是新的，没有受影响的配置id
		return nil, nil
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
	}
	for key := range old {
		// 删除的key
		if m.requireManager.RequireByCount(key) > 0 {
			return nil, fmt.Errorf("variable %s %w", key, eosc.ErrorRequire)
		}
	}
	
	return affectIds, nil
}

func (m *Variables) SetByNamespace(namespace string, variables map[string]string) error {
	_, err := m.Check(namespace, variables)
	if err != nil {
		return err
	}
	m.variables.Set(namespace, variables)
	return nil
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
	newMap := make(map[string]string)
	for key, value := range v {
		newMap[key] = value
	}
	return newMap, ok
}

func (m *Variables) GetByNamespace(namespace string) (map[string]string, bool) {
	return m.getByNamespace(namespace)
}

//
//func (m *Variables) GetAll() map[string]string {
//	return m.getAll()
//}

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
