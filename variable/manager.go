package variable

import (
	"encoding/json"
	"fmt"
	"github.com/eolinker/eosc"
	"reflect"
	"strings"
	"sync"
)

var _ IVariable = (*Variables)(nil)

type IVariable interface {
	SetByNamespace(namespace string, variables map[string]string) error
	GetByNamespace(namespace string) (map[string]string, bool)
	SetVariablesById(id string, variables []string)
	Unmarshal(buf []byte, typ reflect.Type) (interface{}, []string, error)
	Check(namespace string, variables map[string]string) ([]string, IVariable, error)
	Get(id string) (string, bool)
	Len() int
}

type Variables struct {
	// data 变量数据
	lock           sync.RWMutex
	data           map[string]map[string]string
	requireManager IRequires
}

func (m *Variables) Get(id string) (string, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	namespace, key := readId(id)
	vs, has := m.data[namespace]
	if has {
		val, has := vs[key]
		return val, has
	}
	return "", false
}
func readId(id string) (namespace, key string) {

	if i := strings.Index(id, "@"); i > 0 {
		namespace = id[i+1:]
		key = id[:i]
		if len(key) == 0 {
			// "@xxxx"
			key = namespace
			namespace = "default"
		}
		if len(namespace) == 0 {
			namespace = "default"
		}

	} else {
		key = id

	}
	return
}
func (m *Variables) Len() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	l := 0
	for _, vs := range m.data {
		l += len(vs)
	}
	return l
}

func (m *Variables) Unmarshal(buf []byte, typ reflect.Type) (interface{}, []string, error) {
	return NewParse(m).Unmarshal(buf, typ)
}

func NewVariables(data map[string][]byte) IVariable {
	v := &Variables{data: make(map[string]map[string]string, len(data)), requireManager: NewRequireManager()}
	for namespace, value := range data {
		nvs := make(map[string]string)
		err := json.Unmarshal(value, &nvs)
		if err != nil {
			continue
		}
		v.data[namespace] = nvs
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

func (m *Variables) check(namespace string, variables map[string]string) ([]string, error) {
	old, has := m.getByNamespace(namespace)
	if !has {
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
func (m *Variables) Check(namespace string, variables map[string]string) ([]string, IVariable, error) {
	// variables的key为：{变量名}@{namespace}，如：v1@default
	m.lock.RLock()
	defer m.lock.RUnlock()
	vs, err := m.check(namespace, variables)
	if err != nil {
		return nil, nil, err
	}
	clone := m.clone()
	clone.data[namespace] = variables
	return vs, clone, nil
}

func (m *Variables) SetByNamespace(namespace string, variables map[string]string) error {
	m.lock.Lock()
	defer m.lock.Lock()
	_, err := m.check(namespace, variables)
	if err != nil {
		return err
	}
	m.data[namespace] = variables
	return nil
}
func (m *Variables) clone() *Variables {
	m.lock.RLock()
	defer m.lock.RUnlock()

	data := make(map[string]map[string]string)
	for namespace, vs := range m.data {
		tmp := make(map[string]string)
		for k, v := range vs {
			tmp[k] = v
		}
		data[namespace] = tmp
	}
	return &Variables{
		lock:           sync.RWMutex{},
		data:           data,
		requireManager: m.requireManager,
	}
}
func (m *Variables) getByNamespace(namespace string) (map[string]string, bool) {

	variables, has := m.data[namespace]
	if !has {
		return nil, false
	}
	newMap := make(map[string]string)
	for key, value := range variables {
		newMap[key] = value
	}
	return newMap, true
}

func (m *Variables) GetByNamespace(namespace string) (map[string]string, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.getByNamespace(namespace)
}
