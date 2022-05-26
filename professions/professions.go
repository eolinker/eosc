package professions

import (
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
)

var _ IProfessions = (*Professions)(nil)

type IProfessions interface {
	//eosc.IProfessions
	Get(name string) (*Profession, bool)
	Sort() []*Profession
	List() []*Profession
	Delete(name string) error
	Set(name string, profession *eosc.ProfessionConfig) error
	Reset(configs []*eosc.ProfessionConfig)
}

type Professions struct {
	//eosc.IProfessions
	data    eosc.IUntyped
	extends eosc.IExtenderDrivers
}

func (ps *Professions) Delete(name string) error {
	// todo 校验
	_, has := ps.data.Del(name)
	if !has {
		return eosc.ErrorProfessionNotExist
	}
	return nil
}

func (ps *Professions) List() []*Profession {
	list := ps.data.List()
	rs := make([]*Profession, len(list))
	for i, v := range list {
		rs[i] = v.(*Profession)
	}
	return rs
}

func (ps *Professions) Sort() []*Profession {
	list := ps.List()

	sl := make([]*Profession, 0, len(list))
	sm := make(map[string]int)
	index := 0
	for i, p := range list {
		if p.Mod == eosc.ProfessionConfig_Singleton {
			sl = append(sl, p)
			sm[p.Name] = index
			index++
			list[i] = nil
		}
	}
	for len(list) > 0 {
		sc := 0
		for i, v := range list {
			if v == nil {
				sc++
				continue
			}
			dependenciesHas := 0
			for _, dep := range v.Dependencies {
				if _, has := sm[dep]; !has {
					break
				}

				dependenciesHas++
			}
			if dependenciesHas == len(v.Dependencies) {
				sl = append(sl, v)
				sm[v.Name] = index
				index++
				sc++
				list[i] = nil
			}
		}
		if sc == 0 {
			// todo profession dependencies error
			break
		}
		tmp := list[:0]
		for _, v := range list {
			if v != nil {
				tmp = append(tmp, v)
			}
		}
		list = tmp
	}
	return sl
}

func NewProfessions(extends eosc.IExtenderDrivers) *Professions {
	ps := &Professions{
		//IProfessions: professions.NewProfessions(),
		extends: extends,
	}
	return ps
}

func (ps *Professions) Set(name string, c *eosc.ProfessionConfig) error {

	p := NewProfession(c, ps.extends)
	ps.data.Set(name, p)

	// todo refresh worker
	return nil
}
func (ps *Professions) Reset(configs []*eosc.ProfessionConfig) {
	data := eosc.NewUntyped()

	for _, c := range configs {
		log.Debug("add profession config:", c)
		p := NewProfession(c, ps.extends)
		data.Set(c.Name, p)
	}
	ps.data = data
}
func (ps *Professions) Get(name string) (*Profession, bool) {
	p, b := ps.data.Get(name)
	log.Debug("get profession:", name, ":", b, "->", p)
	if !b {
		log.Debug("professions data:", ps.data)
	}
	if b {

		return p.(*Profession), true
	}
	return nil, false
}
