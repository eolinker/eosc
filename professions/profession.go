package professions

import (
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
)

var _ IProfession = (*Profession)(nil)

type IProfession interface {
	GetDriver(name string) (eosc.IExtenderDriver, bool)
}
type Profession struct {
	*eosc.ProfessionConfig
	driverData map[string]*eosc.DriverConfig
	drivers    ITypedProfessionDrivers
}

func (p *Profession) DriverConfig(name string) (*eosc.DriverConfig, bool) {
	di, has := p.driverData[name]
	return di, has
}
func (p *Profession) GetDriver(name string) (eosc.IExtenderDriver, bool) {
	return p.drivers.Get(name)
}

func NewProfession(professionConfig *eosc.ProfessionConfig, extends eosc.IExtenderDrivers) *Profession {
	ds := NewProfessionDrivers()
	driverData := make(map[string]*eosc.DriverConfig)
	for _, driverConfig := range professionConfig.Drivers {
		driverData[driverConfig.Name] = driverConfig
		df, b := extends.GetDriver(driverConfig.Id)
		if !b {
			log.Warn("driver not exist:", driverConfig.Id)
			continue
		}

		var params map[string]interface{}
		if driverConfig.Params != nil {
			params = make(map[string]interface{})
			for k, v := range driverConfig.Params {
				params[k] = v
			}
		}

		driver, err := df.Create(professionConfig.Name, driverConfig.Name, driverConfig.Label, driverConfig.Desc, params)
		if err != nil {
			log.Warnf("create driver %s of %s :%v", driverConfig.Id, professionConfig.Name, err)
			continue
		}

		ds.data.Set(driverConfig.Name, driver)
	}
	return &Profession{
		ProfessionConfig: professionConfig,
		drivers:          ds,
		driverData:       driverData,
	}
}
