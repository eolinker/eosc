package process_worker

import (
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"sync"

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
}

var _ IWorkers = (*WorkerManager)(nil)

func ReadWorkers(r io.Reader) []*eosc.WorkerData {
	frame, err := utils.ReadFrame(r)
	if err != nil {
		log.Warn("read  workers frame :", err)

		return nil
	}

	wd := new(eosc.WorkersData)
	if e := proto.Unmarshal(frame, wd); e != nil {
		log.Warn("unmarshal workers data :", e)
		return nil
	}
	return wd.Data
}

type WorkerManager struct {
	locker         sync.Mutex
	professions    IProfessions
	data           ITypedWorkers
	requireManager IWorkerRequireManager
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
		return errors.New("not exist")
	}
	if wm.requireManager.RequireByCount(id) > 0 {
		return errors.New("required")
	}

	err := worker.Stop()
	if err != nil {
		return err
	}
	wm.data.Del(id)
	wm.requireManager.Del(id)
	return nil
}

func (wm *WorkerManager) Get(id string) (eosc.IWorker, bool) {
	return wm.data.Get(id)
}

func NewWorkerManager(professions IProfessions) *WorkerManager {
	return &WorkerManager{professions: professions, data: NewTypedWorkers()}
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
	wm.locker.Lock()
	defer wm.locker.Unlock()

	// if update
	o, has := wm.data.Get(id)
	if has {
		err := o.Reset(conf, requires)
		if err != nil {
			return err
		}
		wm.requireManager.Set(id, requires)
		return nil
	}
	// create
	w, err := driver.Create(id, name, conf, requires)
	if err != nil {
		return err
	}
	// start
	e := w.Start()
	if e != nil {
		return e
	}
	// store
	wm.data.Set(id, NewWorker(id, profession, name, driverName, body, w, p, driver))
	wm.requireManager.Set(id, requires)

	return nil
}
