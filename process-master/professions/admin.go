package professions

import (
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/admin"
)

type Profession struct {
	drivers    eosc.IUntyped
	profession *admin.ProfessionInfo
}

func NewProfession(profession *admin.ProfessionInfo) *Profession {
	return &Profession{profession: profession, drivers: eosc.NewUntyped()}
}

func (p *Profession) Drivers() []*eosc.DriverInfo {
	drivers := p.drivers.List()
	ds := make([]*eosc.DriverInfo, 0, len(drivers))
	for _, d := range drivers {
		v, ok := d.(*eosc.DriverDetail)
		if !ok {
			continue
		}
		ds = append(ds, &v.DriverInfo)
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

func (p *Profession) Render(driver string) (*admin.Render, bool) {
	return nil, false
}

func (p *Profession) Renders() map[string]*admin.Render {
	return nil
}

func (p *Profession) DriversItem() []admin.Item {
	drivers := p.drivers.List()
	ds := make([]admin.Item, 0, len(drivers))
	for _, d := range drivers {
		v, ok := d.(*eosc.DriverInfo)
		if !ok {
			continue
		}
		ds = append(ds, admin.Item{
			Value: v.Name,
			Label: v.Label,
		})
	}
	return ds
}

func (p *Profession) Info() *admin.ProfessionInfo {
	return p.profession
}

func (p *Profession) SetDriver(name string, detail *eosc.DriverDetail) error {
	p.drivers.Set(name, detail)
	return nil
}

func (p *Profession) DeleteDriver(name string) error {
	p.drivers.Del(name)
	return nil
}
