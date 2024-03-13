/*
 * Copyright (c) 2024. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package workers

import (
	"context"
	"fmt"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/require"

	"github.com/eolinker/eosc/professions"
	"github.com/eolinker/eosc/utils/config"
)

type IWorkers interface {
	Begin(ctx context.Context) ITransactionCtx
	GetEmployee(profession, name string) (*WorkerInfo, error)
	Export() map[string][]*WorkerInfo
	ListEmployees(profession string) ([]interface{}, error)
	GetProfession(profession string) (*professions.Profession, bool)
	CheckDelete(ids ...string) (requires []string)
	Rebuild(id string) error
}

type imlWorkers struct {
	professions    professions.IProfessions
	data           *WorkerDatas
	requireManager eosc.IRequires
	variables      eosc.IVariable
}

func (oe *imlWorkers) GetProfession(profession string) (*professions.Profession, bool) {
	return oe.professions.Get(profession)
}

func (oe *imlWorkers) Begin(ctx context.Context) ITransactionCtx {
	//TODO implement me
	panic("implement me")
}

func NewWorkers(professions professions.IProfessions, data *WorkerDatas, variables eosc.IVariable) IWorkers {

	ws := &imlWorkers{requireManager: require.NewRequireManager()}
	ws.init(professions, data, variables)
	return ws
}
func (oe *imlWorkers) init(professions professions.IProfessions, data *WorkerDatas, variables eosc.IVariable) {
	oe.professions = professions
	oe.data = data
	oe.variables = variables

	ps := oe.professions.Sort()

	pm := make(map[string][]*WorkerInfo)
	for _, wd := range oe.data.List() {
		pm[wd.config.Profession] = append(pm[wd.config.Profession], wd)
	}

	for _, pw := range ps {
		for _, v := range pm[pw.Name] {
			_, err := oe.set(v.config.Id, v.config.Profession, v.config.Name, v.config.Driver, v.config.Version, v.config.Description, v.config.Body)
			if err != nil {
				log.Errorf("init %s:%s", v.config.Id, err.Error())
			}
		}
	}
}
func (oe *imlWorkers) ListEmployees(profession string) ([]interface{}, error) {
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

func (oe *imlWorkers) update(profession, name, driver, version, desc string, data IData) (*WorkerInfo, error) {
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
	w, err := oe.set(id, profession, name, driver, version, desc, body)
	if err != nil {
		return nil, err
	}

	return w, nil

}
func (oe *imlWorkers) Rebuild(id string) error {
	info, has := oe.data.GetInfo(id)
	if has {
		_, err := oe.set(id, info.config.Profession, info.config.Name, info.config.Driver, info.config.Version, info.config.Description, info.config.Body)
		return err
	}
	return nil

}
func (oe *imlWorkers) Export() map[string][]*WorkerInfo {
	all := make(map[string][]*WorkerInfo)
	for _, w := range oe.data.All() {
		all[w.config.Profession] = append(all[w.config.Profession], w)
	}
	return all
}
func (oe *imlWorkers) CheckDelete(ids ...string) (requires []string) {
	for _, id := range ids {
		if oe.requireManager.RequireByCount(id) > 0 {
			requires = append(requires, id)
		}
	}
	return requires
}
func (oe *imlWorkers) delete(id string) (*WorkerInfo, error) {

	worker, has := oe.data.GetInfo(id)
	if !has {
		return nil, eosc.ErrorWorkerNotExits
	}

	if oe.requireManager.RequireByCount(id) > 0 {
		return nil, eosc.ErrorRequire
	}

	if destroy, ok := worker.worker.(eosc.IWorkerDestroy); ok {
		err := destroy.Destroy()
		if err != nil {
			return nil, err
		}
	}
	oe.data.Del(id)
	oe.requireManager.Del(id)
	oe.variables.RemoveRequire(id)
	return worker, nil
}

func (oe *imlWorkers) GetEmployee(profession, name string) (*WorkerInfo, error) {

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

func (oe *imlWorkers) set(id, profession, name, driverName, version, desc string, body []byte) (*WorkerInfo, error) {

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
		if err = dc.Check(conf, requires); err != nil {
			return nil, err
		}
	}
	wInfo, hasInfo := oe.data.GetInfo(id)
	if hasInfo && wInfo.worker != nil {

		e := wInfo.worker.Reset(conf, requires)
		if e != nil {
			return nil, e
		}
		oe.requireManager.Set(id, getIds(requires))
		wInfo.reset(driverName, version, desc, body, wInfo.worker, driver.ConfigType())
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
		wInfo = NewWorkerInfo(worker, id, profession, name, driverName, version, desc, eosc.Now(), eosc.Now(), body, driver.ConfigType())
	} else {
		wInfo.reset(driverName, version, desc, body, worker, driver.ConfigType())
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
