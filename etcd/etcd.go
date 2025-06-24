package etcd

import (
	"go.etcd.io/etcd/api/v3/version"
	"go.etcd.io/etcd/server/v3/etcdserver/api/membership"
)

type Versions version.Versions
type Etcd interface {
	IsLeader() (bool, []string)
	KV
	Watch(prefix string, handler ServiceHandler)
	HandlerLeader(h ...ILeaderStateHandler)
	Join(target string) error
	Leave() error
	Close() error
	Info() *Node
	Nodes() []*Node
	Status() ClusterInfo
	Version() Versions
	Remove(name string) error
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
	LeaderChange(isLeader bool, isRestart bool)
}
type Info = *membership.Member
