package dispatcher

import (
	"strings"
	"sync"

	"github.com/eolinker/eosc"
)

type Data struct {
	lock sync.RWMutex
	data map[string]map[string][]byte
}

func NewMyData(data map[string]map[string][]byte) *Data {
	if data == nil {
		data = make(map[string]map[string][]byte)
	}
	return &Data{
		data: data,
	}
}

func (d *Data) DoEvent(event IEvent) {
	d.lock.Lock()
	defer d.lock.Unlock()

	switch strings.ToLower(event.Event()) {
	case eosc.EventDel:
		d.delete(event.Namespace(), event.Key())
	case eosc.EventSet:
		d.set(event.Namespace(), event.Key(), event.Data())
	case eosc.EventInit, eosc.EventReset:
		d.reset(event.All())
	}
}
func (d *Data) getAll() map[string]map[string][]byte {
	m := make(map[string]map[string][]byte)

	for n, ns := range d.data {
		nsn := make(map[string][]byte)
		for k, v := range ns {
			nsn[k] = v
		}
		m[n] = nsn
	}
	return m
}
func (d *Data) reset(all map[string]map[string][]byte) {
	d.data = all

}
func (d *Data) delete(namespace, key string) {
	sub, has := d.data[namespace]
	if !has {
		return
	}

	delete(sub, key)
}

func (d *Data) set(namespace, key string, data []byte) {
	sub, has := d.data[namespace]
	if !has {
		sub = make(map[string][]byte)
		d.data[namespace] = sub
	}
	sub[key] = data
}
func (d *Data) GetNamespace(namespace string)(map[string][]byte,bool) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	data,has := d.data[namespace]
	if has{
		tmp:=make(map[string][]byte)
		for k,v:=range data{
			tmp[k] = v
		}
		return tmp,has
	}
	return nil,false
}
func (d *Data) GET() map[string]map[string][]byte {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.getAll()
}

type InitEvent map[string]map[string][]byte

func (r InitEvent) Namespace() string {
	return ""
}

func (r InitEvent) Event() string {
	return eosc.EventInit
}

func (r InitEvent) Key() string {
	return ""
}

func (r InitEvent) Data() []byte {
	return nil
}
func (r InitEvent) All() map[string]map[string][]byte {
	return r
}
