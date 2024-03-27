package workers

import (
	"fmt"
	"github.com/eolinker/eosc/professions"
	"github.com/eolinker/eosc/utils/config"
	"sync"

	//port_reqiure "github.com/eolinker/eosc/common/port-reqiure"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
)

type IWorkers interface {
	eosc.IWorkers
	Del(id string) error
	Set(id, profession, name, driverName string, body []byte, variable eosc.IVariable) error
	Update(id string, variable eosc.IVariable) error
	Reset(wdl []*eosc.WorkerConfig, variable eosc.IVariable) error
}

type ConfigCache struct {
	profession string
	name       string
	driver     string
	config     []byte
}

var _ IWorkers = (*Workers)(nil)

type Workers struct {
	locker      sync.Mutex
	professions professions.IProfessions
	variables   eosc.IVariable
	data        *WorkerDatas
	configs     map[string]*ConfigCache
}

func (wm *Workers) Update(id string, variable eosc.IVariable) error {
	wm.locker.Lock()
	defer wm.locker.Unlock()
	con, has := wm.configs[id]
	if !has {
		return nil
	}
	return wm.set(id, con.profession, con.name, con.driver, con.config, variable)
}

func (wm *Workers) Del(id string) error {
	wm.locker.Lock()
	defer wm.locker.Unlock()

	worker, has := wm.data.Get(id)
	if !has {
		return eosc.ErrorWorkerNotExits
	}

	err := worker.Stop()
	if err != nil {
		return err
	}
	wm.data.Del(id)
	destroy, ok := worker.(eosc.IWorkerDestroy)
	if ok {
		destroy.Destroy()
	}
	delete(wm.configs, id)
	wm.variables.RemoveRequire(id)
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
		professions: profession,
		locker:      sync.Mutex{},
		data:        NewTypedWorkers(),
		configs:     make(map[string]*ConfigCache),
	}
}

func (wm *Workers) Reset(wdl []*eosc.WorkerConfig, variable eosc.IVariable) error {

	ps := wm.professions.Sort()

	pm := make(map[string][]*eosc.WorkerConfig)
	for _, wd := range wdl {
		pm[wd.Profession] = append(pm[wd.Profession], wd)
	}

	wm.locker.Lock()
	defer wm.locker.Unlock()
	wm.variables = variable
	wm.configs = make(map[string]*ConfigCache)
	olddata := wm.data
	wm.data = NewTypedWorkers()

	log.Debug("worker init... size is ", len(wdl))
	for _, p := range ps {
		log.Debug("init profession:", p.Name)
		for _, wd := range pm[p.Name] {
			old, has := olddata.Del(wd.Id)
			if has {
				wm.data.Set(wd.Id, old)
			}
			log.Debug("init set:", wd.Id, " ", wd.Profession, " ", wd.Name, " ", wd.Driver, " ", string(wd.Body))
			if err := wm.set(wd.Id, wd.Profession, wd.Name, wd.Driver, wd.Body, variable); err != nil {
				log.Error("init set worker: ", err)
				continue
			}
		}
	}
	for _, ov := range olddata.All() {
		variable.RemoveRequire(ov.Id())
		ov.Stop()
	}
	return nil
}

func (wm *Workers) Set(id, profession, name, driverName string, body []byte, variable eosc.IVariable) error {
	wm.locker.Lock()
	defer wm.locker.Unlock()

	return wm.set(id, profession, name, driverName, body, variable)
}

func (wm *Workers) set(id, profession, name, driverName string, body []byte, variable eosc.IVariable) error {
	log.Debug("set:", id, ",", profession, ",", name, ",", driverName)
	p, has := wm.professions.Get(profession)
	if !has {
		return fmt.Errorf("%s:%w", profession, eosc.ErrorProfessionNotExist)
	}
	driver, has := p.GetDriver(driverName)
	if !has {
		return fmt.Errorf("%s,%w", driverName, eosc.ErrorDriverNotExist)
	}
	conf, useVariables, err := variable.Unmarshal(body, driver.ConfigType())
	if err != nil {
		return fmt.Errorf("worker unmarshal error:%s", err)
	}

	requires, err := config.CheckConfig(conf, wm)
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
		wm.variables.SetRequire(id, useVariables)
		wm.configs[id] = &ConfigCache{
			profession: profession,
			name:       name,
			driver:     driverName,
			config:     body,
		}
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
	}

	// store
	wm.data.Set(id, worker)
	wm.variables.SetRequire(id, useVariables)
	wm.configs[id] = &ConfigCache{
		profession: profession,
		name:       name,
		driver:     driverName,
		config:     body,
	}
	log.Debug("worker-data set worker done:", id)
	return nil
}
