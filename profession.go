package eosc

import (
	"github.com/eolinker/eosc/internal"
	"strings"
)

var _ IProfession = (*Profession)(nil)

type IProfession interface {
	getId(name string) (string, bool)
	setId(name, id string) error
	del(name string) error
	ids() []string
	getDriver(name string) (IProfessionDriverInfo, bool)
	getDrivers() []IProfessionDriverInfo
	genInfo(v *StoreValue) interface{}

	Name() string
	Label() string
	Desc() string
	Dependencies() []string
	AppendLabels() []string
}

type Profession struct {
	name         string
	label        string
	desc         string
	dependencies []string
	appendLabels []string
	drivers      IProfessionDrivers

	data internal.IUntyped
}

func (p *Profession) Name() string {
	return p.name

}

func (p *Profession) Label() string {
	return p.label
}

func (p *Profession) Desc() string {
	return p.desc
}

func (p *Profession) Dependencies() []string {
	return p.dependencies
}

func (p *Profession) AppendLabels() []string {
	return p.appendLabels
}

func (p *Profession) getDrivers() []IProfessionDriverInfo {
	return p.drivers.List()
}

func (p *Profession) getId(name string) (string, bool) {
	if o, has := p.data.Get(strings.ToLower(name)); has {
		return o.(string), true
	}
	return "", false
}

func (p *Profession) genInfo(v *StoreValue) interface{} {

	r := make(map[string]interface{})
	r["name"] = v.Name
	r["id"] = v.Id
	r["driver"] = v.Driver
	r["create_time"] = v.CreateTime
	r["update_time"] = v.UpdateTime

	item := make(map[string]interface{})
	e := v.IData.UnMarshal(&item)
	if e != nil {
		return r
	}
	for _, l := range p.appendLabels {
		r[l] = item[l]
	}
	return r
}
func (p *Profession) setId(name, id string) error {

	p.data.Set(strings.ToLower(name), id)
	return nil
}
func (p *Profession) ids() []string {
	list := p.data.List()
	res := make([]string, len(list))
	for i, v := range list {
		res[i] = v.(string)
	}
	return res
}
func (p *Profession) del(name string) error {
	p.data.Del(name)
	return nil
}

func (p *Profession) getDriver(name string) (IProfessionDriverInfo, bool) {
	return p.drivers.Get(name)
}
