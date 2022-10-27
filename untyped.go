// Package eosc SPDX-License-Identifier: Apache-2.0
package eosc

import (
	"sync"
)

type Untyped[K comparable, T any] interface {
	Set(k K, v T)
	Get(k K) (T, bool)
	Del(k K) (T, bool)
	List() []T
	Keys() []K
	All() map[K]T
	Clone() Untyped[K, T]
	Count() int
}

func BuildUntyped[K comparable, T any]() Untyped[K, T] {
	return &tUntyped[K, T]{
		data:  map[K]T{},
		mutex: &sync.RWMutex{},
		sort:  nil,
	}
}

type tUntyped[K comparable, T any] struct {
	data  map[K]T
	sort  []K
	mutex *sync.RWMutex
}

func (u *tUntyped[K, T]) Count() int {
	return len(u.sort)
}

func cloneUntyped[K comparable, T any](data map[K]T, sort []K) *tUntyped[K, T] {
	return &tUntyped[K, T]{
		data:  data,
		sort:  sort,
		mutex: &sync.RWMutex{},
	}
}

func (u *tUntyped[K, T]) Del(name K) (T, bool) {

	u.mutex.Lock()
	v, ok := u.data[name]
	if ok {
		u.sort = remove(u.sort, name)
		delete(u.data, name)
	}

	u.mutex.Unlock()

	return v, ok
}
func (u *tUntyped[K, T]) Set(name K, v T) {

	u.mutex.Lock()
	_, has := u.data[name]
	if !has {
		u.sort = append(u.sort, name)
	}
	u.data[name] = v
	u.mutex.Unlock()
}

func (u *tUntyped[K, T]) Get(name K) (T, bool) {

	u.mutex.RLock()
	v, ok := u.data[name]
	u.mutex.RUnlock()
	return v, ok
}

func (u *tUntyped[K, T]) Clone() Untyped[K, T] {

	u.mutex.RLock()
	res := make(map[K]T, len(u.data))
	for k, v := range u.data {
		res[k] = v
	}
	sort := make([]K, len(u.sort))
	copy(sort, u.sort)
	u.mutex.RUnlock()
	return cloneUntyped(res, sort)
}
func (u *tUntyped[K, T]) All() map[K]T {
	u.mutex.RLock()
	res := make(map[K]T, len(u.data))
	for k, v := range u.data {
		res[k] = v
	}
	u.mutex.RUnlock()
	return res
}

func (u *tUntyped[K, T]) Keys() []K {
	u.mutex.RLock()
	res := make([]K, len(u.data))
	copy(res, u.sort)
	u.mutex.RUnlock()
	return res
}

func (u *tUntyped[K, T]) List() []T {
	u.mutex.RLock()
	res := make([]T, len(u.data))
	for i, k := range u.sort {
		res[i] = u.data[k]
	}
	u.mutex.RUnlock()
	return res
}

func remove[K comparable](src []K, t K) []K {

	for i, v := range src {
		if v == t {
			copy(src[i:], src[i+1:])
			return src[:len(src)-1]
		}
	}
	return src

}
