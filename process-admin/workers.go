package process_admin

import (
	"fmt"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/require"

	"github.com/eolinker/eosc/professions"
	"github.com/eolinker/eosc/utils/config"
	"reflect"
)

type Workers struct {
	professions    professions.IProfessions
	data           *WorkerDatas
	requireManager eosc.IRequires
	variables      eosc.IVariable
}

func NewWorkers(professions professions.IProfessions, data *WorkerDatas, variables eosc.IVariable) *Workers {

	ws := &Workers{professions: professions, data: data, requireManager: require.NewRequireManager(), variables: variables}
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
			_, err := oe.set(v.config.Id, v.config.Profession, v.config.Name, v.config.Driver, v.config.Description, JsonData(v.config.Body))
			if err != nil {
				log.Errorf("init %s:%s", v.config.Id, err.Error())
			}
		}
	}
}
func (oe *Workers) ListEmployees(profession string) ([]interface{}, error) {
	p, has := oe.professions.Get(profession)
	if !has {
		return nil, eosc.ErrorProfessionNotExist
	}
	appendLabels := p.AppendLabels
	all := oe.data.All()
	vs := make([]interface{}, 0, len(all))
	for _, w := range all {
		if w.config.Profession == p.Name {
			vs = append(vs, w.Info(appendLabels...))
		}
	}
	return vs, nil

}

func (oe *Workers) Update(profession, name, driver, desc string, data IData) (*WorkerInfo, error) {
	id, ok := eosc.ToWorkerId(name, profession)
	if !ok {
		return nil, fmt.Errorf("%s@%s:invalid id", name, profession)
	}
	log.Debug("update:", id, " ", profession, ",", name, ",", driver, ",", data)
	if driver == "" {

		employee, err := oe.GetEmployee(profession, name)
		if err != nil {
			return nil, err
		}
		driver = employee.config.Driver
	}
	body, _ := data.Encode()
	w, err := oe.set(id, profession, name, driver, desc, body)
	if err != nil {
		return nil, err
	}

	return w, nil

}
func (oe *Workers) rebuild(id string) error {
	info, has := oe.data.GetInfo(id)
	if has {
		_, err := oe.set(id, info.config.Profession, info.config.Name, info.config.Driver, info.config.Description, info.config.Body)
		return err
	}
	return nil

}
func (oe *Workers) Export() map[string][]*WorkerInfo {
	all := make(map[string][]*WorkerInfo)
	for _, w := range oe.data.All() {
		all[w.config.Profession] = append(all[w.config.Profession], w)
	}
	return all
}
func (oe *Workers) DeleteTest(ids ...string) (requires []string) {
	for _, id := range ids {
		if oe.requireManager.RequireByCount(id) > 0 {
			requires = append(requires, id)
		}
	}
	return requires
}
func (oe *Workers) Delete(id string) (*WorkerInfo, error) {

	worker, has := oe.data.GetInfo(id)
	if !has {
		return nil, eosc.ErrorWorkerNotExits
	}

	if oe.requireManager.RequireByCount(id) > 0 {
		return nil, eosc.ErrorRequire
	}

	oe.data.Del(id)
	destroy, ok := worker.worker.(eosc.IWorkerDestroy)
	if ok {
		destroy.Destroy()
	}
	oe.requireManager.Del(id)
	oe.variables.RemoveRequire(id)
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

func (oe *Workers) set(id, profession, name, driverName, desc string, body []byte) (*WorkerInfo, error) {

	log.Debug("set:", id, ",", profession, ",", name, ",", driverName)
	p, has := oe.professions.Get(profession)
	if !has {
		return nil, fmt.Errorf("%s:%w", profession, eosc.ErrorProfessionNotExist)
	}
	if p.Mod == eosc.ProfessionConfig_Singleton {
		driverName = name
	}
	driver, has := p.GetDriver(driverName)
	if !has {
		return nil, fmt.Errorf("%s,%w", driverName, eosc.ErrorDriverNotExist)
	}

	conf, usedVariables, err := oe.variables.Unmarshal(body, driver.ConfigType())
	if err != nil {
		return nil, err
	}

	requires, err := config.CheckConfig(conf, oe.data)
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
		wInfo.reset(driverName, desc, body, wInfo.worker, driver.ConfigType())
		oe.variables.SetVariablesById(id, usedVariables)
		return wInfo, nil
	}
	// create
	worker, err := driver.Create(id, name, conf, requires)
	if err != nil {
		log.Warn("worker-data set worker create:", err)
		return nil, err
	}

	if !hasInfo {
		wInfo = NewWorkerInfo(worker, id, profession, name, driverName, desc, eosc.Now(), eosc.Now(), body, driver.ConfigType())
	} else {
		wInfo.reset(driverName, desc, body, worker, driver.ConfigType())
	}

	// store
	oe.data.Set(id, wInfo)
	oe.requireManager.Set(id, getIds(requires))
	oe.variables.SetVariablesById(id, usedVariables)
	log.Debug("worker-data set worker done:", id)

	return wInfo, nil
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
