package extenders

import (
	"errors"
	"fmt"

	"github.com/eolinker/eosc/log"

	"github.com/eolinker/eosc/extends"

	"github.com/eolinker/eosc/service"

	"github.com/eolinker/eosc"
)

var _ ITypedExtenderSetting = (*ExtenderSetting)(nil)

var (
	errExtenderNotExist = errors.New("the extender does not exist")
)

type ITypedExtenderSetting interface {
	Set(group, project, version string)
	Del(group, project string)
	Get(group, project string) (version string, has bool)
	All() map[string]string
	Reset(map[string]string)
	ITypedPlugin
}

type ITypedPlugin interface {
	GetPluginsByExtenderID(extenderID string) ([]*service.Plugin, bool)
	SetPluginsByExtenderID(extenderID string, plugins []*service.Plugin)
	GetPlugins() []*service.Plugin
	GetPluginByID(id string) (*service.Plugin, bool)
	SetPluginByID(id string, plugin *service.Plugin)
}

type ExtenderSetting struct {
	data    eosc.IUntyped
	plugins eosc.IUntyped
}

func (i *ExtenderSetting) GetPluginsByExtenderID(extenderID string) ([]*service.Plugin, bool) {
	extender, has := i.plugins.Get(extenderID)
	if !has {
		return nil, false
	}
	es, ok := extender.(eosc.IUntyped)
	if !ok {
		return nil, false
	}
	list := es.List()
	plugins := make([]*service.Plugin, 0, len(list))
	for _, p := range list {
		v, ok := p.(*service.Plugin)
		if !ok {
			continue
		}
		plugins = append(plugins, v)
	}

	return plugins, true
}

func (i *ExtenderSetting) SetPluginsByExtenderID(extenderID string, plugins []*service.Plugin) {
	extender, has := i.plugins.Get(extenderID)
	if !has {
		log.Error(errExtenderNotExist)
		return
	}
	es, ok := extender.(eosc.IUntyped)
	if !ok {
		log.Error(errExtenderNotExist)
		return
	}
	for _, p := range plugins {
		es.Set(p.Id, p)
	}
}

func (i *ExtenderSetting) GetPlugins() []*service.Plugin {
	all := i.plugins.All()
	plugins := make([]*service.Plugin, 0, len(all)*5)
	for _, a := range all {
		v, ok := a.(eosc.IUntyped)
		if !ok {
			continue
		}
		ps := v.List()
		for _, p := range ps {
			pl, ok := p.(*service.Plugin)
			if !ok {
				continue
			}
			plugins = append(plugins, pl)
		}
	}
	return plugins
}

func (i *ExtenderSetting) getPlugin(id string) (eosc.IUntyped, bool) {
	group, project, _, err := extends.DecodeExtenderId(id)
	if err != nil {
		return nil, false
	}
	plugins, has := i.plugins.Get(extends.FormatProject(group, project))
	if !has {
		return nil, false
	}
	ps, ok := plugins.(eosc.IUntyped)

	return ps, ok
}

func (i *ExtenderSetting) GetPluginByID(id string) (*service.Plugin, bool) {
	plugins, has := i.getPlugin(id)
	if !has {
		return nil, false
	}
	v, has := plugins.Get(id)
	if !has {
		return nil, false
	}
	p, ok := v.(*service.Plugin)

	return p, ok
}

func (i *ExtenderSetting) SetPluginByID(id string, plugin *service.Plugin) {
	plugins, has := i.getPlugin(id)
	if !has {
		return
	}
	plugins.Set(id, plugin)
}

func (i *ExtenderSetting) Reset(m map[string]string) {
	data := eosc.NewUntyped()
	if m != nil {
		for k, v := range m {
			data.Set(k, v)
		}
	}
	i.data = data
}

func NewInstallData() *ExtenderSetting {
	return &ExtenderSetting{data: eosc.NewUntyped(), plugins: eosc.NewUntyped()}
}

func (i *ExtenderSetting) Set(group, project, version string) {
	id := toId(group, project)
	i.data.Set(id, version)
}

func (i *ExtenderSetting) Del(group, project string) {
	id := toId(group, project)
	i.data.Del(id)
}

func (i *ExtenderSetting) Get(group, project string) (string, bool) {
	id := toId(group, project)
	v, has := i.data.Get(id)
	if has {
		return v.(string), true
	}
	return "", false
}

func (i *ExtenderSetting) All() map[string]string {
	data := i.data.All()
	mk := make(map[string]string)
	for k, v := range data {
		mk[k] = v.(string)
	}
	return mk
}

func toId(group, project string) string {
	return fmt.Sprint(group, ":", project)
}
