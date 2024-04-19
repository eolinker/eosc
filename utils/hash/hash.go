package hash

import (
	"encoding/json"
	"sync"
)

var (
	_ Hash[string, string] = (*imlHash[string, string])(nil)
)

type Hash[K comparable, T any] interface {
	Clone() Hash[K, T]
	Set(key K, value T)
	Get(key K) (T, bool)
	Delete(key K)
	Keys() []K
	Map() map[K]T
	Len() int
}

type imlHash[K comparable, T any] struct {
	lock sync.RWMutex
	data map[K]T
}

func (h *imlHash[K, T]) MarshalBinary() (data []byte, err error) {
	return h.MarshalJSON()
}

func (h *imlHash[K, T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.data)
}

func (h *imlHash[K, T]) Len() int {
	h.lock.RLock()
	defer h.lock.RUnlock()
	return len(h.data)
}

func NewHash[K comparable, T any](data map[K]T) Hash[K, T] {

	h := &imlHash[K, T]{
		data: make(map[K]T),
	}
	for k, v := range data {
		h.data[k] = v
	}
	return h
}

func (h *imlHash[K, T]) cloneData() map[K]T {
	data := make(map[K]T, len(h.data))
	for k, v := range h.data {
		data[k] = v
	}
	return data
}
func (h *imlHash[K, T]) Clone() Hash[K, T] {

	h.lock.RLock()
	defer h.lock.RUnlock()
	return &imlHash[K, T]{data: h.cloneData()}
}

func (h *imlHash[K, T]) Set(key K, value T) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.data == nil {
		h.data = make(map[K]T)
	}
	h.data[key] = value
}

func (h *imlHash[K, T]) Get(key K) (T, bool) {
	h.lock.RLock()
	defer h.lock.RUnlock()
	value, ok := h.data[key]
	return value, ok
}

func (h *imlHash[K, T]) Delete(key K) {
	h.lock.Lock()
	defer h.lock.Unlock()
	delete(h.data, key)
}

func (h *imlHash[K, T]) Keys() []K {
	h.lock.RLock()
	defer h.lock.RUnlock()
	keys := make([]K, 0, len(h.data))
	for k := range h.data {
		keys = append(keys, k)
	}
	return keys
}

func (h *imlHash[K, T]) Map() map[K]T {
	h.lock.RLock()
	defer h.lock.RUnlock()

	return h.cloneData()
}
