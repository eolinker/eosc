package port_reqiure

import (
	"strconv"
	"sync"

	"github.com/eolinker/eosc"
)

var _ IPortsRequire = (*PortsRequire)(nil)

type IPortsRequire interface {
	Set(id string, ports []int)
	Del(id string)
	All() []int32
}

type PortsRequire struct {
	locker  sync.Mutex
	workers eosc.IUntyped
	ports   eosc.IUntyped
}

func NewPortsRequire() IPortsRequire {
	return &PortsRequire{
		locker:  sync.Mutex{},
		workers: eosc.NewUntyped(),
		ports:   eosc.NewUntyped(),
	}
}

func (p *PortsRequire) Set(id string, ports []int) {
	p.locker.Lock()
	defer p.locker.Unlock()

	p.del(id)

	if len(ports) == 0 {
		return
	}

	p.workers.Set(id, ports)

	for _, port := range ports {
		p.add(id, port)
	}
}

func (p *PortsRequire) Del(id string) {
	p.locker.Lock()
	defer p.locker.Unlock()
	p.del(id)
}
func (p *PortsRequire) del(id string) {
	portList, has := p.workers.Del(id)
	if has {
		pv := portList.([]int)
		for _, v := range pv {
			p.remove(id, v)
		}
	}
}
func (p *PortsRequire) add(id string, port int) {
	pv := strconv.Itoa(port)

	ids, has := p.ports.Get(pv)
	if !has {
		p.ports.Set(pv, []string{id})
		return
	}

	idsv := ids.([]string)

	for _, idv := range idsv {
		if idv == id {
			return
		}
	}

	idsv = append(idsv, id)
	p.ports.Set(pv, idsv)

}
func (p *PortsRequire) remove(id string, port int) {
	pv := strconv.Itoa(port)

	ids, has := p.ports.Get(pv)
	if !has {
		return
	}

	idsv := ids.([]string)
	for i, idv := range idsv {
		if idv == id {
			idsv = append(idsv[:i], idsv[i+1:]...)
			if len(idsv) > 0 {
				p.ports.Set(pv, idsv)
			} else {
				p.ports.Del(pv)
			}
			return
		}
	}
}

func (p *PortsRequire) All() []int32 {
	p.locker.Lock()
	list := p.ports.Keys()
	p.locker.Unlock()

	rs := make([]int32, len(list))
	for i, pv := range list {
		port, _ := strconv.Atoi(pv)
		rs[i] = int32(port)
	}
	return rs
}
