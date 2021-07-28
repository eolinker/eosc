package eosc

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

	Infos() []ProfessionInfo
}

type Professions struct {
	data    iProfessionsData
	store   IStore
	infos   []ProfessionInfo
	workers IWorkers
}

func (ps *Professions) Infos() []ProfessionInfo {
	return ps.infos
}


type tProfessionUntyped struct {
	data IUntyped
}

func newTProfessionUntyped() iProfessionsData {
	return &tProfessionUntyped{
		data: NewUntyped(),
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
