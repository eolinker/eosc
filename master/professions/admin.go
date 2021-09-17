package professions

import (
	"errors"

	"github.com/eolinker/eosc"
)

func (p *Professions) Render(profession, driver string) (*eosc.Render, error) {
	return nil, nil
}

func (p *Professions) Renders(profession string) (map[string]*eosc.Render, error) {
	panic("implement me")
}

func (p *Professions) Drivers(profession string) ([]eosc.DriverInfo, error) {
	if value, has := p.drivers.Get(profession); has {
		vl, ok := value.(eosc.IUntyped)
		if !ok {
			return nil, errors.New("invalid type")
		}
		drivers := vl.List()
		ds := make([]eosc.DriverInfo, 0, len(drivers))
		for _, d := range drivers {
			v, ok := d.(eosc.DriverInfo)
			if !ok {
				continue
			}
			ds = append(ds, v)
		}
		return ds, nil
	}
	return nil, errors.New("invalid profession")
}

func (p *Professions) DriverInfo(profession, driver string) (eosc.DriverDetail, error) {
	panic("implement me")
}

func (p *Professions) DriversItem(profession string) ([]eosc.Item, error) {
	panic("implement me")
}

func (p *Professions) ListProfessions() []eosc.ProfessionInfo {
	panic("implement me")
}
