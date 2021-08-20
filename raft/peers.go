package raft

import (
	"sync"
)

type Peers struct {
	currentPeers map[uint64]string
	configCount  int
	mu           sync.RWMutex
}

func NewPeers(peers map[uint64]string, count int) *Peers {
	return &Peers{
		currentPeers: peers,
		configCount:  count,
		mu:           sync.RWMutex{},
	}
}

func (p *Peers) GetPeerNum() int {
	return len(p.currentPeers)
}

func (p *Peers) GetConfigCount() int {
	return p.configCount
}

// CheckExist 判断host对应的ID是否存在
func (p *Peers) CheckExist(host string) (uint64, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for k, v := range p.currentPeers {
		if v == host {
			return k, true
		}
	}
	return 0, false
}

func (p *Peers) GetPeerByID(id uint64) (string, bool) {
	p.mu.RLock()
	v, ok := p.currentPeers[id]
	p.mu.RUnlock()
	return v, ok
}

func (p *Peers) SetPeer(id uint64, value string) {
	p.mu.Lock()
	p.currentPeers[id] = value
	p.mu.Unlock()
	p.configCount++
}

//GetAllPeers 获取所有节点列表
func (p *Peers) GetAllPeers() map[uint64]string {
	res := make(map[uint64]string)
	p.mu.RLock()
	for k, v := range p.currentPeers {
		res[k] = v
	}
	p.mu.RUnlock()
	return res
}

//UpdatePeerByID 通过ID更新节点信息
func (p *Peers) UpdatePeerByID(id uint64, value string) {
	p.mu.Lock()
	p.currentPeers[id] = value
	p.mu.Unlock()
}

//DeletePeerByID 通过ID删除节点
func (p *Peers) DeletePeerByID(id uint64) {
	p.mu.Lock()
	delete(p.currentPeers, id)
	p.mu.Unlock()
	p.configCount++
}
