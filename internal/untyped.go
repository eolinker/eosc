// SPDX-License-Identifier: Apache-2.0
package internal

import (
	"strings"
	"sync"
)

type IUntyped interface {
	Set(name string, v interface{})
	Get(name string) (interface{}, bool)
	Del(name string) (interface{}, bool)
	List() []interface{}
	Keys() []string
	All() map[string]interface{}
	Clone() IUntyped
	Count() int
}

func NewUntyped() IUntyped {
	return &tUntyped{
		data:  map[string]interface{}{},
		mutex: &sync.RWMutex{},
		sort:  nil,
	}
}

type tUntyped struct {
	data  map[string]interface{}
	sort  []string
	mutex *sync.RWMutex
}

func (u *tUntyped) Count() int {
	return len(u.sort)
}

func cloneUntyped(data map[string]interface{}, sort []string) *tUntyped {
	return &tUntyped{
		data:  data,
		sort:  sort,
		mutex: &sync.RWMutex{},
	}
}

func (u *tUntyped) Del(name string) (interface{}, bool) {
	name = strings.ToLower(name)
	u.mutex.Lock()
	v, ok := u.data[name]
	if ok {
		u.sort = remove(u.sort, name)
		delete(u.data, name)
	}

	u.mutex.Unlock()

	return v, ok
}
func (u *tUntyped) Set(name string, v interface{}) {
	name = strings.ToLower(name)
	u.mutex.Lock()
	_, has := u.data[name]
	if !has {
		u.sort = append(u.sort, name)
	}
	u.data[name] = v
	u.mutex.Unlock()
}

func (u *tUntyped) Get(name string) (interface{}, bool) {
	name = strings.ToLower(name)

	u.mutex.RLock()
	v, ok := u.data[name]
	u.mutex.RUnlock()
	return v, ok
}

func (u *tUntyped) Clone() IUntyped {

	u.mutex.RLock()
	res := make(map[string]interface{}, len(u.data))
	for k, v := range u.data {
		res[k] = v
	}
	sort := make([]string, len(u.sort))
	copy(sort, u.sort)
	u.mutex.RUnlock()
	return cloneUntyped(res, sort)
}
func (u *tUntyped) All() map[string]interface{} {
	u.mutex.RLock()
	res := make(map[string]interface{}, len(u.data))
	for k, v := range u.data {
		res[k] = v
	}
	u.mutex.RUnlock()
	return res
}
func (u *tUntyped) Keys() []string {
	u.mutex.RLock()
	res := make([]string, len(u.data))
	copy(res, u.sort)
	u.mutex.RUnlock()
	return res
}

func (u *tUntyped) List() []interface{} {
	u.mutex.RLock()
	res := make([]interface{}, len(u.data))
	for i, k := range u.sort {
		res[i] = u.data[k]
	}
	u.mutex.RUnlock()
	return res
}

func remove(src []string, t string) []string {

	for i, v := range src {
		if v == t {
			copy(src[i:], src[i+1:])
			return src[:len(src)-1]
		}
	}
	return src

}
