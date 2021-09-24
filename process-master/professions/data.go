package professions

import "github.com/eolinker/eosc"

type untypeProfessionData interface {
	Set(name string, profession *Profession)
	List() []*Profession
	Del(name string) (*Profession, bool)
	Get(name string) (*Profession, bool)
	Data() []*eosc.ProfessionConfig
}
type TProfessionData struct {
	data eosc.IUntyped
}

func (d *TProfessionData) Data() []*eosc.ProfessionConfig {
	ls := d.List()
	rs := make([]*eosc.ProfessionConfig, len(ls))
	for i, v := range ls {
		rs[i] = v.config
	}
	return rs
}

func NewProfessionData() untypeProfessionData {
	return &TProfessionData{
		data: eosc.NewUntyped(),
	}
}

func (d *TProfessionData) Set(name string, profession *Profession) {
	d.data.Set(name, profession)
}

func (d *TProfessionData) List() []*Profession {
	list := d.data.List()
	result := make([]*Profession, len(list))
	for i, v := range list {
		result[i] = v.(*Profession)
	}
	return result
}

func (d *TProfessionData) Del(name string) (*Profession, bool) {
	v, has := d.data.Del(name)
	if has {
		vp, ok := v.(*Profession)
		return vp, ok
	}
	return nil, false
}

func (d *TProfessionData) Get(name string) (*Profession, bool) {
	v, has := d.data.Get(name)
	if has {
		vp, ok := v.(*Profession)
		return vp, ok
	}
	return nil, false
}
