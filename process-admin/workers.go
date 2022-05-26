package process_admin

import (
	"fmt"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/professions"
	require "github.com/eolinker/eosc/workers/require"
	"reflect"
)

type Workers struct {
	professions    professions.IProfessions
	data           *WorkerDatas
	requireManager require.IRequires
}

func NewWorkers(professions professions.IProfessions, data *WorkerDatas) *Workers {

	ws := &Workers{professions: professions, data: data, requireManager: require.NewRequireManager()}
	ws.init()
	return ws
}
func (oe *Workers) init() {
	ps := oe.professions.Sort()

	pm := make(map[string][]*WorkerInfo)
	for _, wd := range oe.data.List() {
		pm[wd.config.Profession] = append(pm[wd.config.Profession], wd)
	}

	for _, pw := range ps {
		for _, v := range pm[pw.Name] {
			oe.set(v.config.Id, v.config.Profession, v.config.Name, v.config.Driver, JsonData(v.config.Body))
		}
	}
}
func (oe *Workers) ListEmployees(profession string) ([]*WorkerInfo, error) {
	p, has := oe.professions.Get(profession)
	if !has {
		return nil, eosc.ErrorProfessionNotExist
	}
	all := oe.data.All()
	vs := make([]*WorkerInfo, len(all))
	for _, w := range all {
		if w.config.Profession == p.Name {
			vs = append(vs, w)
		}
	}
	return vs, nil

}

func (oe *Workers) Update(profession, name string, driver string, data IData) (*WorkerInfo, error) {
	id, ok := eosc.ToWorkerId(name, profession)
	if !ok {
		return nil, fmt.Errorf("%s@%s:invalid id", name, profession)
	}
	w, err := oe.set(id, profession, name, driver, data)
	if err != nil {
		return nil, err
	}

	return w, nil

}

func (oe *Workers) Export() map[string][]*WorkerInfo {
	all := make(map[string][]*WorkerInfo)
	for _, w := range oe.data.All() {
		all[w.config.Profession] = append(all[w.config.Profession], w)
	}
	return all
}
func (oe *Workers) Delete(profession, name string) (*WorkerInfo, error) {
	id, ok := eosc.ToWorkerId(name, profession)
	if !ok {
		return nil, fmt.Errorf("%s %w", profession, ErrorNotMatch)
	}
	worker, has := oe.data.GetInfo(id)
	if !has {
		return nil, eosc.ErrorWorkerNotExits
	}
	if oe.requireManager.RequireByCount(id) > 0 {
		return nil, eosc.ErrorRequire
	}

	err := worker.worker.Stop()
	if err != nil {
		return nil, err
	}
	oe.data.Del(id)
	oe.requireManager.Del(id)

	return worker, nil
}

func (oe *Workers) GetEmployee(profession, name string) (*WorkerInfo, error) {

	id, ok := eosc.ToWorkerId(name, profession)
	if !ok {
		return nil, fmt.Errorf("%s %w", id, ErrorNotMatch)
	}
	d, has := oe.data.GetInfo(id)
	if !has {
		return nil, fmt.Errorf("%s %w", id, ErrorNotExist)
	}
	return d, nil
}

func (oe *Workers) set(id, profession, name, driverName string, data IData) (*WorkerInfo, error) {

	log.Debug("set:", id, ",", profession, ",", name, ",", driverName)
	p, has := oe.professions.Get(profession)
	if !has {
		return nil, fmt.Errorf("%s:%w", profession, eosc.ErrorProfessionNotExist)
	}
	driver, has := p.GetDriver(driverName)
	if !has {
		return nil, fmt.Errorf("%s,%w", driverName, eosc.ErrorDriverNotExist)
	}

	conf := newConfig(driver.ConfigType())

	err := data.UnMarshal(conf)
	if err != nil {
		return nil, err
	}
	requires, err := eosc.CheckConfig(conf, oe.data)
	if err != nil {
		return nil, err
	}
	if dc, ok := driver.(eosc.IExtenderConfigChecker); ok {
		if e := dc.Check(conf, requires); err != nil {
			return nil, e
		}
	}
	wInfo, hasInfo := oe.data.GetInfo(id)
	if hasInfo && wInfo.worker != nil {

		e := wInfo.worker.Reset(conf, requires)
		if e != nil {
			return nil, e
		}
		oe.requireManager.Set(id, getIds(requires))
		wInfo.reset(driverName, conf, wInfo.worker)
		return wInfo, nil
	}
	// create
	worker, err := driver.Create(id, name, conf, requires)
	if err != nil {
		log.Warn("worker-data set worker create:", err)
		return nil, err
	}
	// start
	e := worker.Start()
	if e != nil {
		log.Warn("worker-data set worker start:", e)
		return nil, e
	}
	if !hasInfo {
		wInfo = NewWorkerInfo(worker, id, profession, name, driverName, eosc.Now(), eosc.Now(), conf)
	} else {
		wInfo.reset(driverName, conf, worker)
	}
	// store
	oe.data.Set(id, wInfo)
	oe.requireManager.Set(id, getIds(requires))

	log.Debug("worker-data set worker done:", id)
	return wInfo, nil
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
