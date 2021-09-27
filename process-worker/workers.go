package process_worker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sync"

	"github.com/eolinker/eosc/listener"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/utils"
	"github.com/golang/protobuf/proto"
)

type IWorkers interface {
	eosc.IWorkers
	Del(id string) error
	Check(id, profession, name, driverName string, body []byte) error
	Set(id, profession, name, driverName string, body []byte) error
	RequiredCount(id string) int
	ResourcesPort() []int32
}

var _ IWorkers = (*WorkerManager)(nil)

func ReadWorkers(r io.Reader) []*eosc.WorkerData {
	frame, err := utils.ReadFrame(r)
	if err != nil {
		log.Warn("read  workerIds frame :", err)

		return nil
	}

	wd := new(eosc.WorkersData)
	if e := proto.Unmarshal(frame, wd); e != nil {
		log.Warn("unmarshal workerIds data :", e)
		return nil
	}
	return wd.Data
}

type WorkerManager struct {
	locker         sync.Mutex
	professions    IProfessions
	data           ITypedWorkers
	requireManager IWorkerRequireManager
	portsRequire   listener.IPortsRequire
}

func (wm *WorkerManager) ResourcesPort() []int32 {
	return wm.portsRequire.All()
}

func (wm *WorkerManager) RequiredCount(id string) int {
	return wm.requireManager.RequireByCount(id)
}

func (wm *WorkerManager) Check(id, profession, name, driverName string, body []byte) error {
	p, has := wm.professions.Get(profession)
	if !has {
		return errors.New("profession not exist")
	}
	driver, has := p.GetDriver(driverName)
	if !has {
		return errors.New("driver not exist")
	}

	configType := driver.ConfigType()
	conf := reflect.New(configType).Interface()

	err := json.Unmarshal(body, conf)
	if err != nil {
		return err
	}
	requires, err := eosc.CheckConfig(conf, wm)
	if err != nil {
		return err
	}
	if dc, ok := driver.(eosc.IProfessionDriverCheckConfig); ok {
		if err := dc.Check(conf, requires); err != nil {
			return err
		}
	}
	return nil

}

func (wm *WorkerManager) Del(id string) error {
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
	wm.portsRequire.Del(id)

	return nil
}

func (wm *WorkerManager) Get(id string) (eosc.IWorker, bool) {
	return wm.data.Get(id)
}

func NewWorkerManager(professions IProfessions) *WorkerManager {
	return &WorkerManager{
		locker:         sync.Mutex{},
		professions:    professions,
		data:           NewTypedWorkers(),
		requireManager: NewWorkerRequireManager(),
		portsRequire:   listener.NewPortsRequire(),
	}
}

func (wm *WorkerManager) Init(wdl []*eosc.WorkerData) error {

	ps := wm.professions.Sort()

	pm := make(map[string][]*eosc.WorkerData)
	for _, wd := range wdl {
		pm[wd.Profession] = append(pm[wd.Profession], wd)
	}

	for _, p := range ps {
		for _, wd := range pm[p.Name] {
			if err := wm.Set(wd.Id, wd.Profession, wd.Name, wd.Driver, wd.Body); err != nil {
				return err
			}
		}
	}
	return nil
}

func (wm *WorkerManager) Set(id, profession, name, driverName string, body []byte) error {

	p, has := wm.professions.Get(profession)
	if !has {
		return fmt.Errorf("%s:%w", profession, eosc.ErrorProfessionNotExist)
	}
	driver, has := p.GetDriver(driverName)
	if !has {
		return fmt.Errorf("%s:%w", driverName, eosc.ErrorDriverNotExist)
	}

	configType := driver.ConfigType()
	conf := reflect.New(configType).Interface()

	err := json.Unmarshal(body, conf)
	if err != nil {
		return err
	}
	requires, err := eosc.CheckConfig(conf, wm)
	if err != nil {
		return err
	}
	if dc, ok := driver.(eosc.IProfessionDriverCheckConfig); ok {
		if e := dc.Check(conf, requires); err != nil {

			return e
		}
	}
	wm.locker.Lock()
	defer wm.locker.Unlock()

	// if update
	o, has := wm.data.Get(id)
	if has {
		e := o.Reset(conf, requires)
		if e != nil {
			return e
		}
		wm.requireManager.Set(id, getIds(requires))
		if res, ok := o.target.(eosc.IWorkerResources); ok {
			wm.portsRequire.Set(id, res.Ports())
		}
		return nil
	}
	// create
	worker, err := driver.Create(id, name, conf, requires)
	if err != nil {
		return err
	}
	// start
	e := worker.Start()
	if e != nil {
		return e
	}

	// store
	wm.data.Set(id, NewWorker(id, profession, name, driverName, body, worker, p, driver))
	wm.requireManager.Set(id, getIds(requires))
	if res, ok := worker.(eosc.IWorkerResources); ok {
		wm.portsRequire.Set(id, res.Ports())
	}
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
