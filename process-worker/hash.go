package process_worker

import (
	"github.com/eolinker/eosc"
	"sync"
)

type iCustomerVar interface {
	eosc.ICustomerVar
	Set(key string, value map[string]string)
}

type imlCustomerVar struct {
	rwLock sync.RWMutex
	data   map[string]map[string]string
}

func (i *imlCustomerVar) GetAll(key string) (map[string]string, bool) {
	i.rwLock.RLock()
	defer i.rwLock.RUnlock()
	m, has := i.data[key]
	return m, has
}

func (i *imlCustomerVar) Get(key string, field string) (string, bool) {
	i.rwLock.RLock()
	defer i.rwLock.RUnlock()
	m, has := i.data[key]
	if has {
		if v, has := m[field]; has {
			return v, true
		}
	}
	return "", false
}

func (i *imlCustomerVar) Exists(key string, field string) bool {
	_, has := i.Get(key, field)
	return has
}

func (i *imlCustomerVar) Set(key string, value map[string]string) {
	i.rwLock.Lock()
	defer i.rwLock.Unlock()
	if len(value) == 0 {
		delete(i.data, key)
		return
	}
	i.data[key] = value
}

func newImlCustomerVar() iCustomerVar {
	return &imlCustomerVar{
		rwLock: sync.RWMutex{},
		data:   make(map[string]map[string]string),
	}
}
