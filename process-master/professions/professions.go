package professions

import (
	"github.com/eolinker/eosc"
)

const (
	SpaceProfession = "profession"
)

type Professions struct {
	data ITypeProfessionData
}

func (p *Professions) Set(name string, profession *eosc.ProfessionConfig) error {
	adminProfession := NewProfession(profession)

	p.data.Set(name, adminProfession)
	return nil
}
func (p *Professions) All() []*eosc.ProfessionConfig {
	return p.data.Data()
}
func (p *Professions) Delete(name string) error {
	p.data.Del(name)
	return nil
}

func (p *Professions) Infos() []*eosc.ProfessionInfo {
	professions := p.data.List()
	ps := make([]*eosc.ProfessionInfo, 0, len(professions))
	for _, pv := range professions {
		ps = append(ps, pv.info)
	}
	return ps
}

func (p *Professions) GetProfession(name string) (eosc.IProfession, bool) {
	vl, has := p.data.Get(name)

	return vl, has
}

func (p *Professions) Reset(professions []*eosc.ProfessionConfig) {
	pfs := NewProfessionData()
	for _, pf := range professions {
		adminProfession := NewProfession(pf)
		pfs.Set(pf.Name, adminProfession)
	}
	p.data = pfs
}

func NewProfessions() *Professions {
	return &Professions{
		data: NewProfessionData(),
	}
}
