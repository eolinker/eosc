package etcdRaft

import "net/http"

type Etcd interface {
	http.Handler
	IsLeader() (bool, []string, error)
	KV
	//Watch(prefix string,handler ServiceHandler)
}
type KValue struct {
	Key   []byte
	Value []byte
}
type KV interface {
	Put(key, value string) error
	Delete(key string) error
}

type ServiceHandler interface {
	KV
	Reset([]*KValue)
}
