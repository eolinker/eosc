package raft

import (
	"fmt"
	"sync"
)

type Peers struct {
	peers       map[uint64]*NodeInfo
	peersByAddr map[string]*NodeInfo
	index       uint64
	mu          sync.RWMutex
}

func NewPeers() *Peers {
	return &Peers{
		peers:       make(map[uint64]*NodeInfo),
		peersByAddr: make(map[string]*NodeInfo),
		index:       1,
		mu:          sync.RWMutex{},
	}
}

func (p *Peers) GetPeerNum() int {
	return len(p.peers)
}

func (p *Peers) Index() uint64 {
	return p.index
}

// CheckExist 判断host对应的ID是否存在
func (p *Peers) CheckExist(host string) (uint64, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if v, ok := p.peersByAddr[host]; ok {
		return v.ID, true
	}
	return 0, false
}

func (p *Peers) GetPeerByID(id uint64) (*NodeInfo, bool) {
	p.mu.RLock()
	v, ok := p.peers[id]
	p.mu.RUnlock()
	return v, ok
}

func (p *Peers) SetPeer(id uint64, value *NodeInfo) {
	p.mu.Lock()

	addr := fmt.Sprintf("%s://%s", value.Protocol, value.BroadcastIP)
	if value.BroadcastPort > 0 {
		addr = fmt.Sprintf("%s:%d", addr, value.BroadcastPort)
	}
	if value.Addr != "" {
		addr = value.Addr
	}
	value.Addr = addr
	p.peers[id] = value
	p.peersByAddr[addr] = value
	p.mu.Unlock()
	if p.index < id {
		p.index = id
	}
}

//GetAllPeers 获取所有节点列表
func (p *Peers) GetAllPeers() map[uint64]*NodeInfo {
	res := make(map[uint64]*NodeInfo)
	p.mu.RLock()
	for k, v := range p.peers {
		res[k] = v
	}
	p.mu.RUnlock()
	return res
}

//DeletePeerByID 通过ID删除节点
func (p *Peers) DeletePeerByID(id uint64) {
	info, has := p.GetPeerByID(id)
	if !has {
		return
	}

	p.mu.Lock()
	delete(p.peers, id)
	delete(p.peersByAddr, fmt.Sprintf("%s:%d", info.BroadcastIP, info.BroadcastPort))
	p.mu.Unlock()
}
