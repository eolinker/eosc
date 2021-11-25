package extenders

import (
	"sync"

	"github.com/eolinker/eosc"
)

type ITypedExtenderData interface {
	Reset([]*Extender)
	Set(extender *Extender)
	Del(id string)
	All() []*Extender
}

type ExtenderData struct {
	locker sync.RWMutex
	data   eosc.IUntyped
}

func (ed *ExtenderData) Reset(extenders []*Extender) {
	data := eosc.NewUntyped()

	for _, e := range extenders {
		data.Set(e.Id, e)
	}
	ed.locker.Lock()
	ed.data = data
	ed.locker.Unlock()
}

func (ed *ExtenderData) Set(extender *Extender) {
	ed.locker.Lock()
	ed.data.Set(extender.Id, extender)
	ed.locker.Unlock()
}

func (ed *ExtenderData) Del(id string) {
	ed.locker.Lock()
	ed.data.Del(id)
	ed.locker.Unlock()
}

func (ed *ExtenderData) All() []*Extender {
	list := ed.data.List()
	extenders := make([]*Extender, 0, len(list))
	for _, extender := range list {
		v, ok := extender.(*Extender)
		if !ok {
			continue
		}
		extenders = append(extenders, v)
	}
	return extenders
}
