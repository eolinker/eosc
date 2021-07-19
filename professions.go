package eosc

import (
	"github.com/eolinker/eosc/internal"
)

var (
	_ iProfessionsData = (*tProfessionUntyped)(nil)
	_ IProfessions     = (*Professions)(nil)
)

type iProfessionsData interface {
	get(name string) (IProfession, bool)
	add(name string, p IProfession)
	list() []IProfession
}

type IProfessions interface {
	iProfessionsData
	Infos() []ProfessionInfo
}

type Professions struct {
	iProfessionsData
	store IStore
	infos []ProfessionInfo
}

func (ps *Professions) Infos() []ProfessionInfo {
	return ps.infos
}

func checkProfessions(infos []ProfessionInfo) ([]ProfessionInfo, error) {

	less := make([]ProfessionInfo, len(infos))
	copy(less, infos)
	plist := make([]ProfessionInfo, 0, len(less))
	exist := make(map[string]int)
	do := 1
	for do > 0 && len(less) > 0 {
		do = 0
		ls := less
		less = make([]ProfessionInfo, 0, len(ls))
	FIND:
		for _, v := range ls {

			for _, d := range v.Dependencies {
				if _, has := exist[d]; !has {
					less = append(less, v)
					continue FIND
				}
			}
			plist = append(plist, v)
			exist[v.Name] = 1
			do++

		}
	}
	if len(less) > 0 {
		return nil, ErrorProfessionDependencies
	}
	return plist, nil

}

type tProfessionUntyped struct {
	data internal.IUntyped
}

func newTProfessionUntyped() *tProfessionUntyped {
	return &tProfessionUntyped{
		data: internal.NewUntyped(),
	}
}

func (ps *tProfessionUntyped) list() []IProfession {
	ls := ps.data.List()
	rs := make([]IProfession, 0, len(ls))
	for _, i := range ls {
		rs = append(rs, i.(IProfession))
	}
	return rs
}

func (ps *tProfessionUntyped) get(name string) (IProfession, bool) {
	if o, h := ps.data.Get(name); h {
		return o.(IProfession), true
	}
	return nil, false
}

func (ps *tProfessionUntyped) add(name string, p IProfession) {
	ps.data.Set(name, p)
}
