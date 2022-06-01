package process_admin

import (
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/professions"
	"github.com/eolinker/eosc/workers/require"
)

type ProfessionsRequire struct {
	professions.IProfessions
	requires require.IRequires
}

func (p *ProfessionsRequire) Delete(name string) error {
	if err := p.IProfessions.Delete(name); err != nil {
		return err
	}
	p.requires.Del(name)
	return nil
}

func (p *ProfessionsRequire) Set(name string, profession *eosc.ProfessionConfig) error {

	if err := p.IProfessions.Set(name, profession); err != nil {
		return err
	}
	drivers := make([]string, 0, len(profession.Drivers))
	for _, d := range profession.Drivers {
		drivers = append(drivers, d.Id)
	}
	p.requires.Set(name, drivers)
	return nil
}

func (p *ProfessionsRequire) Reset(configs []*eosc.ProfessionConfig) {
	p.IProfessions.Reset(configs)
	for _, c := range configs {
		drivers := make([]string, 0, len(c.Drivers))
		for _, d := range c.Drivers {
			drivers = append(drivers, d.Id)
		}
		p.requires.Set(c.Name, drivers)
	}
}

func NewProfessionsRequire(professions professions.IProfessions, requires require.IRequires) *ProfessionsRequire {
	return &ProfessionsRequire{IProfessions: professions, requires: requires}
}
