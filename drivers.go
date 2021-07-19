package eosc

import (
	"github.com/eolinker/eosc/internal"

)

type IProfessionDrivers interface {
	Get(name string)(driver IProfessionDriverInfo,has bool)
	List()[]IProfessionDriverInfo
	Keys()[]string
}
type _TProfessionDriver struct {
	IProfessionDriver
	extendInfo ExtendInfo
	driverInfo DriverInfo
}

func newProfessionDriver(IProfessionDriver IProfessionDriver, extendInfo ExtendInfo, driverInfo DriverInfo) *_TProfessionDriver {
	return &_TProfessionDriver{IProfessionDriver: IProfessionDriver, extendInfo: extendInfo, driverInfo: driverInfo}
}


func (p *_TProfessionDriver) ExtendInfo() ExtendInfo {
	return p.extendInfo
}

func (p *_TProfessionDriver) DriverInfo() DriverInfo {
	return p.driverInfo
}

type ProfessionDrivers struct {
	data internal.IUntyped
}

func (d *ProfessionDrivers) Get(name string) (  IProfessionDriverInfo,   bool) {

	 if o,h := d.data.Get(name);h{
	 	driver,has := o.(IProfessionDriverInfo)
	 	return driver,has
	 }
	 return nil,false
}

func (d *ProfessionDrivers) List() []IProfessionDriverInfo {
	list:=d.data.List()
	res:=make([]IProfessionDriverInfo,len(list))
	for  i,v:=range list{
		res[i] = v.(IProfessionDriverInfo)
	}
	return res
}

func (d *ProfessionDrivers) Keys() []string {
	return d.Keys()
}

func NewDrivers(drivers []*_TProfessionDriver)IProfessionDrivers {
	data := internal.NewUntyped()
	for _,d:=range drivers{
		data.Set(d.driverInfo.Name,d)
	}
	return &ProfessionDrivers{
		data:data,
	}
}