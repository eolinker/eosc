package workers

import (
	"encoding/json"
	"fmt"
	"github.com/eolinker/eosc/professions"
	"github.com/eolinker/eosc/workers/require"
	"reflect"
	"sync"

	//port_reqiure "github.com/eolinker/eosc/common/port-reqiure"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
)

type IWorkers interface {
	eosc.IWorkers
	Del(id string) error
	//Check(id, profession, name, driverName string, body []byte) error
	Set(id, profession, name, driverName string, body []byte) error

	//RequiredCount(id string) int
	Reset(wdl []*eosc.WorkerConfig) error
	//All() []*Worker
}

var _ IWorkers = (*Workers)(nil)

type Workers struct {
	locker         sync.Mutex
	professions    professions.IProfessions
	data           *WorkerDatas
	requireManager require.IWorkerRequireManager
	//portsRequire   port_reqiure.IPortsRequire
}

func (wm *Workers) RequiredCount(id string) int {
	return wm.requireManager.RequireByCount(id)
}

func (wm *Workers) Check(id, profession, name, driverName string, body []byte) error {
	wm.locker.Lock()
	defer wm.locker.Unlock()

	p, has := wm.professions.Get(profession)
	if !has {
		return fmt.Errorf("%s:%w", profession, eosc.ErrorProfessionNotExist)
	}
	driver, has := p.GetDriver(driverName)
	if !has {
		return fmt.Errorf("%s,%w", driverName, eosc.ErrorDriverNotExist)
	}

	conf := newConfig(driver.ConfigType())
	err := json.Unmarshal(body, conf)
	if err != nil {
		return err
	}
	requires, err := eosc.CheckConfig(conf, wm)
	if err != nil {
		return err
	}
	if dc, ok := driver.(eosc.IExtenderConfigChecker); ok {
		if err := dc.Check(conf, requires); err != nil {
			return err
		}
	}
	return nil

}

func (wm *Workers) Del(id string) error {
	wm.locker.Lock()
	defer wm.locker.Unlock()

	worker, has := wm.data.Get(id)
	if !has {
		return eosc.ErrorWorkerNotExits
	}
	if wm.requireManager.RequireByCount(id) > 0 {
		return eosc.ErrorRequire
	}

	err := worker.Stop()
	if err != nil {
		return err
	}
	wm.data.Del(id)
	wm.requireManager.Del(id)
	//wm.portsRequire.Del(id)

	return nil
}

func (wm *Workers) Get(id string) (eosc.IWorker, bool) {
	w, has := wm.data.Get(id)
	if has {
		return w.(eosc.IWorker), true
	}
	return nil, false
}

func NewWorkerManager(profession professions.IProfessions) *Workers {
	return &Workers{
		professions:    profession,
		locker:         sync.Mutex{},
		data:           NewTypedWorkers(),
		requireManager: require.NewWorkerRequireManager(),
		//portsRequire:   port_reqiure.NewPortsRequire(),
	}
}

func (wm *Workers) Reset(wdl []*eosc.WorkerConfig) error {
	ps := wm.professions.Sort()

	pm := make(map[string][]*eosc.WorkerConfig)
	for _, wd := range wdl {
		pm[wd.Profession] = append(pm[wd.Profession], wd)
	}

	wm.locker.Lock()
	defer wm.locker.Unlock()

	olddata := wm.data
	wm.data = NewTypedWorkers()

	log.Debug("worker init... size is ", len(wdl))
	for _, p := range ps {
		for _, wd := range pm[p.Name] {
			old, has := olddata.Del(wd.Id)
			if has {
				wm.data.Set(wd.Id, old)
			}
			log.Debug("init set:", wd.Id, " ", wd.Profession, " ", wd.Name, " ", wd.Driver, " ", string(wd.Body))
			if err := wm.set(wd.Id, wd.Profession, wd.Name, wd.Driver, wd.Body); err != nil {
				log.Error("init set worker: ", err)
				continue
			}
		}
	}
	for _, ov := range olddata.All() {
		ov.Stop()
	}
	return nil
}

func (wm *Workers) Set(id, profession, name, driverName string, body []byte) error {
	wm.locker.Lock()
	defer wm.locker.Unlock()
	return wm.set(id, profession, name, driverName, body)
}

func (wm *Workers) set(id, profession, name, driverName string, body []byte) error {

	log.Debug("set:", id, ",", profession, ",", name, ",", driverName)
	p, has := wm.professions.Get(profession)
	if !has {
		return fmt.Errorf("%s:%w", profession, eosc.ErrorProfessionNotExist)
	}
	driver, has := p.GetDriver(driverName)
	if !has {
		return fmt.Errorf("%s,%w", driverName, eosc.ErrorDriverNotExist)
	}

	conf := newConfig(driver.ConfigType())

	err := json.Unmarshal(body, conf)
	if err != nil {
		log.Debug(string(body))
		return err
	}
	requires, err := eosc.CheckConfig(conf, wm)
	if err != nil {
		return err
	}
	if dc, ok := driver.(eosc.IExtenderConfigChecker); ok {
		if e := dc.Check(conf, requires); err != nil {
			return e
		}
	}
	o, has := wm.data.Get(id)
	if has {

		e := o.Reset(conf, requires)
		if e != nil {
			return e
		}
		wm.requireManager.Set(id, getIds(requires))
		//if res, ok := o.IWorker.(eosc.IWorkerResources); ok {
		//	wm.portsRequire.Set(id, res.Ports())
		//}
		return nil
	}
	// create
	worker, err := driver.Create(id, name, conf, requires)
	if err != nil {
		log.Warn("worker-data set worker create:", err)
		return err
	}
	// start
	e := worker.Start()
	if e != nil {
		log.Warn("worker-data set worker start:", e)
		return e
	}

	// store
	wm.data.Set(id, worker)
	wm.requireManager.Set(id, getIds(requires))
	//if res, ok := worker.(eosc.IWorkerResources); ok {
	//	wm.portsRequire.Set(id, res.Ports())
	//}
	log.Debug("worker-data set worker done:", id)
	return nil
}

func getIds(m map[eosc.RequireId]interface{}) []string {
	if len(m) == 0 {
		return nil
	}
	rs := make([]string, 0, len(m))
	for k := range m {
		rs = append(rs, string(k))
	}
	return rs
}

func newConfig(t reflect.Type) interface{} {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return reflect.New(t).Interface()
}
