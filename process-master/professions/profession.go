package professions

import (
	"github.com/eolinker/eosc"
)

type Profession struct {
	drivers    eosc.IUntyped
	profession *eosc.ProfessionDetail
	config     *eosc.ProfessionConfig
	info       *eosc.ProfessionInfo
}

func (p *Profession) Mod() eosc.ProfessionConfig_ProfessionMod {
	return p.config.Mod
}

func NewProfession(config *eosc.ProfessionConfig) *Profession {
	dm := eosc.NewUntyped()

	for _, d := range config.Drivers {
		dm.Set(d.Name, eosc.ToDriverDetail(d))
	}
	return &Profession{
		config: config,
		profession: &eosc.ProfessionDetail{
			Name:         config.Name,
			Label:        config.Label,
			Desc:         config.Desc,
			Dependencies: config.Dependencies,
			AppendLabels: config.AppendLabels,
			Drivers:      eosc.ToDriverDetails(config.Drivers),
		},
		drivers: dm,
		info: &eosc.ProfessionInfo{
			Name:  config.Name,
			Label: config.Label,
			Desc:  config.Desc,
		},
	}
}

func (p *Profession) Drivers() []*eosc.DriverInfo {
	drivers := p.drivers.List()
	ds := make([]*eosc.DriverInfo, 0, len(drivers))
	for _, d := range drivers {
		v, ok := d.(*eosc.DriverDetail)
		if !ok {
			continue
		}
		ds = append(ds, &eosc.DriverInfo{
			Id:    v.Id,
			Name:  v.Driver,
			Label: v.Label,
			Desc:  v.Desc,
		})
	}
	return ds
}

func (p *Profession) GetDriver(name string) (*eosc.DriverDetail, bool) {
	vl, has := p.drivers.Get(name)
	if !has {
		return nil, false
	}
	v, ok := vl.(*eosc.DriverDetail)
	return v, ok
}

func (p *Profession) HasDriver(name string) bool {
	_, has := p.drivers.Get(name)
	return has
}

func (p *Profession) AppendAttr() []string {
	return p.profession.AppendLabels
}

func (p *Profession) DriversItem() []*eosc.Item {
	drivers := p.drivers.List()
	ds := make([]*eosc.Item, 0, len(drivers))
	for _, d := range drivers {
		v, ok := d.(*eosc.DriverInfo)
		if !ok {
			continue
		}
		ds = append(ds, &eosc.Item{
			Value: v.Name,
			Label: v.Label,
		})
	}
	return ds
}

func (p *Profession) Detail() *eosc.ProfessionDetail {
	return p.profession
}

func (p *Profession) SetDriver(name string, detail *eosc.DriverConfig) error {

	p.drivers.Set(name, detail)
	return nil
}

func (p *Profession) DeleteDriver(name string) error {
	p.drivers.Del(name)
	return nil
}
