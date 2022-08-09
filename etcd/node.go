package etcd

import "go.etcd.io/etcd/client/pkg/v3/types"

type Node struct {
	Id       string   `json:"id,omitempty"`
	Name     string   `json:"name,omitempty"`
	Admin    []string `json:"admin" json:"admin,omitempty"`
	Server   []string `json:"server" json:"server,omitempty"`
	IsLeader bool     `json:"leader" json:"is_leader,omitempty"`
}

func parseMember(info Info, leader types.ID) *Node {
	return &Node{
		Id:       info.ID.String(),
		Name:     info.Name,
		Admin:    info.PeerURLs,
		Server:   info.ClientURLs,
		IsLeader: leader == info.ID,
	}
}
