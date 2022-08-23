package workers

import (
	"fmt"
	"github.com/eolinker/eosc/professions"
	"github.com/eolinker/eosc/utils/config"
	"github.com/eolinker/eosc/variable"
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
	Set(id, profession, name, driverName string, body []byte, variables map[string]string) error

	//RequiredCount(id string) int
	Reset(wdl []*eosc.WorkerConfig, variables map[string]string) error
	//All() []*Worker
}

var _ IWorkers = (*Workers)(nil)

type Workers struct {
	locker         sync.Mutex
	professions    professions.IProfessions
	data           *WorkerDatas
	requireManager require.IRequires
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
	destroy, ok := worker.(eosc.IWorkerDestroy)
	if ok {
		destroy.Destroy()
	}
	wm.requireManager.Del(id)

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
		requireManager: require.NewRequireManager(),
	}
}

func (wm *Workers) Reset(wdl []*eosc.WorkerConfig, variables map[string]string) error {
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
		log.Debug("init profession:", p.Name)
		for _, wd := range pm[p.Name] {
			old, has := olddata.Del(wd.Id)
			if has {
				wm.data.Set(wd.Id, old)
			}
			log.Debug("init set:", wd.Id, " ", wd.Profession, " ", wd.Name, " ", wd.Driver, " ", string(wd.Body))
			if err := wm.set(wd.Id, wd.Profession, wd.Name, wd.Driver, wd.Body, variables); err != nil {
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

func (wm *Workers) Set(id, profession, name, driverName string, body []byte, variables map[string]string) error {
	wm.locker.Lock()
	defer wm.locker.Unlock()

	return wm.set(id, profession, name, driverName, body, variables)
}

func (wm *Workers) set(id, profession, name, driverName string, body []byte, variables map[string]string) error {
	log.Debug("set:", id, ",", profession, ",", name, ",", driverName)
	p, has := wm.professions.Get(profession)
	if !has {
		return fmt.Errorf("%s:%w", profession, eosc.ErrorProfessionNotExist)
	}
	driver, has := p.GetDriver(driverName)
	if !has {
		return fmt.Errorf("%s,%w", driverName, eosc.ErrorDriverNotExist)
	}
	log.Debug("set body is ", string(body), ",variables is ", variables)
	conf, _, err := variable.NewParse(variables).Unmarshal(body, driver.ConfigType())
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
		wm.requireManager.Set(id, getIds(requires))
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
	wm.requireManager.Set(id, getIds(requires))

	log.Debug("worker-data set worker done:", id)
	return nil
}

func getIds(m map[eosc.RequireId]eosc.IWorker) []string {
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
