package etcd

import "go.etcd.io/etcd/server/v3/etcdserver/api/membership"

type Etcd interface {
	IsLeader() (bool, []string)
	KV
	Watch(prefix string, handler ServiceHandler)
	HandlerLeader(h ...ILeaderStateHandler)
	Join(target string) error
	Leave() error
	Close() error
	Info() Info
}

type KValue struct {
	Key   []byte
	Value []byte
}
type KV interface {
	Put(key string, value []byte) error
	Delete(key string) error
}

type ServiceHandler interface {
	KV
	Reset([]*KValue)
}

type ILeaderStateHandler interface {
	LeaderChange(isLeader bool)
}
type Info *membership.Member
