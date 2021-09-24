package professions

import (
	"encoding/json"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/admin"
)

const (
	SpaceProfession = "profession"
)

type Professions struct {
	professions eosc.IUntyped
}

func (p *Professions) Set(name string, profession *eosc.ProfessionConfig) error {
	adminProfession := NewProfession(profession)

	p.professions.Set(name, adminProfession)
	return nil
}

func (p *Professions) Delete(name string) error {
	p.professions.Del(name)
	return nil
}

func (p *Professions) List() []admin.IProfession {
	professions := p.professions.List()
	ps := make([]admin.IProfession, 0, len(professions))
	for _, p := range professions {
		v, ok := p.(admin.IProfession)
		if !ok {
			continue
		}
		ps = append(ps, v)
	}
	return ps
}

func (p *Professions) Infos() []*eosc.ProfessionInfo {
	professions := p.professions.List()
	ps := make([]*eosc.ProfessionInfo, 0, len(professions))
	for _, p := range professions {
		v, ok := p.(*Profession)
		if !ok {
			continue
		}
		ps = append(ps, v.info)
	}
	return ps
}

func (p *Professions) GetProfession(name string) (admin.IProfession, bool) {
	vl, has := p.professions.Get(name)
	if !has {
		return nil, false
	}
	v, ok := vl.(admin.IProfession)
	if ok {
		return v, ok
	}
	return nil, false
}

func (p *Professions) Reset(professions []*eosc.ProfessionConfig) {
	pfs := eosc.NewUntyped()
	for _, pf := range professions {
		adminProfession := NewProfession(pf)
		pfs.Set(pf.Name, adminProfession)
	}
	p.professions = pfs
}

func (p *Professions) ResetHandler(data []byte) error {
	var professions []*eosc.ProfessionConfig
	err := json.Unmarshal(data, professions)
	if err != nil {
		return err
	}
	p.Reset(professions)
	return nil
}

func (p *Professions) CommitHandler(cmd string, data []byte) error {
	return nil
}

func (p *Professions) Snapshot() []byte {
	professions := p.List()
	data, _ := json.Marshal(professions)
	return data
}

func (p *Professions) ProcessHandler(cmd string, body []byte) ([]byte, error) {
	return nil, nil
}

func NewProfessions() *Professions {
	return &Professions{
		professions: eosc.NewUntyped(),
	}
}
