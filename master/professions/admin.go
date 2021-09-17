package professions

import (
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/admin"
)

type Profession struct {
	drivers    eosc.IUntyped
	appendAttr []string
}

func (p *Profession) Drivers() []*eosc.DriverInfo {
	drivers := p.drivers.List()
	ds := make([]*eosc.DriverInfo, 0, len(drivers))
	for _, d := range drivers {
		v, ok := d.(*eosc.DriverInfo)
		if !ok {
			continue
		}
		ds = append(ds, v)
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
	return p.appendAttr
}

func (p *Profession) Render(driver string) (*admin.Render, bool) {
	panic("implement me")
}

func (p *Profession) Renders() map[string]*admin.Render {
	panic("implement me")
}

func (p *Profession) DriversItem() []admin.Item {
	panic("implement me")
}

func (p *Profession) List() []*eosc.ProfessionInfo {
	panic("implement me")
}

//func (p *Professions) Render(profession, driver string) (*eosc.Render, error) {
//	return nil, nil
//}
//
//func (p *Professions) Renders(profession string) (map[string]*eosc.Render, error) {
//	panic("implement me")
//}
//
//func (p *Professions) Drivers(profession string) ([]eosc.DriverInfo, error) {
//	if value, has := p.drivers.Get(profession); has {
//		vl, ok := value.(eosc.IUntyped)
//		if !ok {
//			return nil, errors.New("invalid type")
//		}
//		drivers := vl.List()
//		ds := make([]eosc.DriverInfo, 0, len(drivers))
//		for _, d := range drivers {
//			v, ok := d.(eosc.DriverInfo)
//			if !ok {
//				continue
//			}
//			ds = append(ds, v)
//		}
//		return ds, nil
//	}
//	return nil, errors.New("invalid profession")
//}
//
//func (p *Professions) DriverInfo(profession, driver string) (eosc.DriverDetail, error) {
//	panic("implement me")
//}
//
//func (p *Professions) DriversItem(profession string) ([]eosc.Item, error) {
//	panic("implement me")
//}
//
//func (p *Professions) ListProfessions() []eosc.ProfessionInfo {
//}

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
