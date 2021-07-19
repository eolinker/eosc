package drivers

import (
	"sync"

	"github.com/eolinker/eosc/log"
)

type TFactory interface {
	Get(name string, config string, f log.Formatter) (log.EntryTransporter, error)
	Driver() string
	Destroy(name string)
}
type TransportEx interface {
	log.EntryTransporter
	Driver() string
}
type iTransportEx struct {
	log.EntryTransporter
	driverName string
}

func NewTransportEx(entryTransporter log.EntryTransporter, driverName string) TransportEx {
	return &iTransportEx{EntryTransporter: entryTransporter, driverName: driverName}
}

func (i *iTransportEx) Driver() string {
	return i.driverName
}

type Drivers struct {
	lock    sync.Mutex
	drivers map[string]TFactory
	caches  map[string]TransportEx
}

func NewDrivers(drivers map[string]TFactory) *Drivers {
	return &Drivers{
		lock:    sync.Mutex{},
		drivers: drivers,
		caches:  nil,
	}
}
func (d *Drivers) GetDriver(name string) (TFactory, bool) {
	d.lock.Lock()
	f, has := d.drivers[name]
	d.lock.Unlock()
	return f, has
}
func (d *Drivers) Cache(caches map[string]TransportEx) {
	if caches == nil {
		caches = make(map[string]TransportEx)
	}
	// 释放已关闭的内容
	d.lock.Lock()
	defer d.lock.Unlock()

	if d.caches != nil {
		for name, told := range d.caches {
			if tnew, has := caches[name]; !has || tnew.Driver() != told.Driver() {
				driver, has := d.drivers[told.Driver()]
				if has {
					driver.Destroy(name)
				} else {
					told.Close()
				}
			}
		}
	}
	d.caches = caches
}
