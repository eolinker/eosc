package extenders

import (
	"fmt"

	"github.com/eolinker/eosc"
)

type ITypedExtenderSetting interface {
	Set(group, project, version string)
	Del(group, project string)
	Get(group, project string) (version string, has bool)
	All() map[string]string
	Reset(map[string]string)
}

type ExtenderSetting struct {
	data eosc.IUntyped
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
	return &ExtenderSetting{data: eosc.NewUntyped()}
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
