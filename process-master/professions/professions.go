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
	fileName    string
	professions eosc.IUntyped
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

func (p *Professions) Infos() []*admin.ProfessionInfo {
	professions := p.professions.List()
	ps := make([]*admin.ProfessionInfo, 0, len(professions))
	for _, p := range professions {
		v, ok := p.(admin.IProfession)
		if !ok {
			continue
		}
		ps = append(ps, v.Info())
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

func (p *Professions) ResetHandler(data []byte) error {
	professions, err := readProfessionConfig(p.fileName)
	if err != nil {
		return err
	}
	for _, pf := range professions {
		adminProfession := NewProfession(
			&admin.ProfessionInfo{
				Name:         pf.Name,
				LocalName:    pf.Name,
				Desc:         pf.Desc,
				Dependencies: pf.Dependencies,
				AppendLabels: pf.AppendLabel,
			})
		for _, d := range pf.Drivers {
			adminProfession.SetDriver(d.Name, &eosc.DriverInfo{
				Id:         d.ID,
				Name:       d.Name,
				Label:      d.Label,
				Desc:       d.Desc,
				Profession: pf.Name,
			})
		}
		p.professions.Set(pf.Name, adminProfession)
	}
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

func NewProfessions(fileName string) *Professions {
	return &Professions{
		fileName:    fileName,
		professions: eosc.NewUntyped(),
	}
}
