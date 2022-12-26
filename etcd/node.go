package etcd

import "go.etcd.io/etcd/client/pkg/v3/types"

type Node struct {
	Id       types.ID `json:"-"`
	ID       string   `json:"id,omitempty"`
	Name     string   `json:"name,omitempty"`
	Peer     []string `json:"peer,omitempty"`
	Admin    []string `json:"admin,omitempty"`
	Server   []string `json:"server,omitempty"`
	IsLeader bool     `json:"leader,omitempty"`
}

type ClusterInfo struct {
	Cluster string  `json:"cluster"`
	Nodes   []*Node `json:"nodes,omitempty"`
}
