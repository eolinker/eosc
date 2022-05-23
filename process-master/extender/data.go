package extender

import (
	"sync"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/service"
)

type Data struct {
	data       eosc.IUntyped
	pluginData *PluginData
	locker     sync.Mutex
}

func NewData() *Data {
	return &Data{data: eosc.NewUntyped(), pluginData: NewPluginData()}
}

func (e *Data) Reset(data []*service.ExtendsInfo) error {
	es := eosc.NewUntyped()
	plugins := make([]*service.Plugin, 0, len(data)*10)
	for _, d := range data {
		es.Set(d.Name, d)
		plugins = append(plugins, d.Plugins...)
	}
	e.locker.Lock()
	e.data = es
	e.pluginData.Reset(plugins)
	e.locker.Unlock()
	return nil
}

func (e *Data) Save(name string, data *service.ExtendsInfo) error {
	e.locker.Lock()
	e.data.Set(name, data)
	e.locker.Unlock()
	return nil
}

func (e *Data) Del(name string) (*service.ExtendsInfo, error) {
	info, ok := e.data.Del(name)
	if !ok {
		return nil, errNotExist
	}
	v, _ := info.(*service.ExtendsInfo)
	return v, nil
}

func (e *Data) List(sortItem string, desc bool) ([]*service.ExtendsInfo, error) {
	// TODO: 补充排序相关内容
	list := e.data.List()
	es := make([]*service.ExtendsInfo, 0, len(list))
	for _, l := range list {
		v, _ := l.(*service.ExtendsInfo)
		es = append(es, v)
	}
	return es, nil
}

func (e *Data) GetPlugin(id string) (*service.Plugin, error) {
	return e.pluginData.Get(id)
}

func (e *Data) Plugins(name string) ([]*service.Plugin, error) {
	ps, ok := e.data.Get(name)
	if !ok {
		return nil, errNotExist
	}
	v, _ := ps.(*service.ExtendsInfo)
	return v.Plugins, nil
}

type PluginData struct {
	locker sync.Mutex
	data   eosc.IUntyped
}

func NewPluginData() *PluginData {
	return &PluginData{data: eosc.NewUntyped()}
}

func (ed *PluginData) Reset(extenders []*service.Plugin) {
	data := eosc.NewUntyped()

	for _, e := range extenders {
		data.Set(e.Id, e)
	}
	ed.locker.Lock()
	ed.data = data
	ed.locker.Unlock()
}

func (ed *PluginData) Set(extender *service.Plugin) {
	ed.locker.Lock()
	ed.data.Set(extender.Id, extender)
	ed.locker.Unlock()
}

func (ed *PluginData) Get(id string) (*service.Plugin, error) {
	ps, ok := ed.data.Get(id)
	if !ok {
		return nil, errNotExist
	}
	v, _ := ps.(*service.Plugin)
	return v, nil
}

func (ed *PluginData) Del(id string) {
	ed.locker.Lock()
	ed.data.Del(id)
	ed.locker.Unlock()
}

func (ed *PluginData) All() []*service.Plugin {
	list := ed.data.List()
	extenders := make([]*service.Plugin, 0, len(list))
	for _, extender := range list {
		v, ok := extender.(*service.Plugin)
		if !ok {
			continue
		}
		extenders = append(extenders, v)
	}
	return extenders
}
