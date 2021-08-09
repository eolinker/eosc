package main

import (
	"go.etcd.io/etcd/client/pkg/v3/types"
	"sync"
)

type Peers struct {
	currentPeers map[types.ID]string
	configCount  int
	mu           sync.RWMutex
}

func NewPeers(peers map[types.ID]string, count int) *Peers {
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
	p.mu.RLock()
	count := p.configCount
	p.mu.RUnlock()
	return count
}

// CheckExist 判断host对应的ID是否存在
func (p *Peers) CheckExist(host string) (types.ID, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for k, v := range p.currentPeers {
		if v == host {
			return k, true
		}
	}
	return 0, false
}

func (p *Peers) GetPeerByID(id types.ID) (string, bool) {
	p.mu.RLock()
	v, ok := p.currentPeers[id]
	p.mu.RUnlock()
	return v, ok
}
func (p *Peers) SetPeer(id types.ID, value string) {
	p.mu.Lock()
	p.currentPeers[id] = value
	p.configCount++
	p.mu.Unlock()
}
func (p *Peers) GetAllPeers() map[types.ID]string {
	p.mu.RLock()
	res := make(map[types.ID]string, len(p.currentPeers))
	for k, v := range p.currentPeers {
		res[k] = v
	}
	p.mu.RUnlock()
	return res
}

func (p *Peers) UpdatePeerByID(id types.ID, value string) {
	p.mu.Lock()
	p.currentPeers[id] = value
	p.mu.Unlock()
}
func (p *Peers) DeletePeerByID(id types.ID) {
	p.mu.Lock()
	delete(p.currentPeers, id)
	p.configCount++
	p.mu.Unlock()
}
