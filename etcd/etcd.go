package etcd

type Etcd interface {
	IsLeader() (bool, []string, error)
	KV
	Watch(prefix string, handler ServiceHandler)
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
