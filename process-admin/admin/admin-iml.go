/*
 * Copyright (c) 2024. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package admin

import (
	"context"
	"fmt"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/professions"
	"github.com/eolinker/eosc/require"
	"github.com/eolinker/eosc/utils/config"
	"sync"
)

type imlAdmin struct {
	professions    professions.IProfessions
	data           *WorkerDatas
	requireManager eosc.IRequires
	variables      eosc.IVariable

	lockTransaction sync.Mutex
}

func (oe *imlAdmin) lock() {
	oe.lockTransaction.Lock()
}

func (oe *imlAdmin) unLock() {
	oe.lockTransaction.Unlock()
}

func (oe *imlAdmin) Get(id string) (eosc.IWorker, bool) {
	return oe.data.Get(id)
}

func (oe *imlAdmin) GetInfo(id string) (*WorkerInfo, bool) {
	return oe.data.GetInfo(id)
}

func (oe *imlAdmin) GetProfession(profession string) (*professions.Profession, bool) {
	return oe.professions.Get(profession)
}

func (oe *imlAdmin) Begin(ctx context.Context) ITransactionCtx {
	oe.lockTransaction.Lock()
	return newImlTransaction(ctx, oe)
}

func NewWorkers(professions professions.IProfessions, data *WorkerDatas, variables eosc.IVariable) IAdmin {

	ws := &imlAdmin{requireManager: require.NewRequireManager()}
	ws.init(professions, data, variables)
	return ws
}
func (oe *imlAdmin) init(professions professions.IProfessions, data *WorkerDatas, variables eosc.IVariable) {
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
			_, err := oe.set(v.config.Id, v.config.Profession, v.config.Name, v.config.Driver, v.config.Version, v.config.Description, v.config.Body, v.config.Update, v.config.Create)
			if err != nil {
				log.Errorf("init %s:%s", v.config.Id, err.Error())
			}
		}
	}
}
func (oe *imlAdmin) ListEmployees(profession string) ([]interface{}, error) {
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

func (oe *imlAdmin) Rebuild(id string) error {
	info, has := oe.data.GetInfo(id)
	if has {
		_, err := oe.set(id, info.config.Profession, info.config.Name, info.config.Driver, info.config.Version, info.config.Description, info.config.Body, info.config.Update, info.config.Create)
		return err
	}
	return nil

}
func (oe *imlAdmin) Export() map[string][]*WorkerInfo {
	all := make(map[string][]*WorkerInfo)
	for _, w := range oe.data.All() {
		all[w.config.Profession] = append(all[w.config.Profession], w)
	}
	return all
}
func (oe *imlAdmin) CheckDelete(ids ...string) (requires []string) {
	for _, id := range ids {
		if oe.requireManager.RequireByCount(id) > 0 {
			requires = append(requires, id)
		}
	}
	return requires
}
func (oe *imlAdmin) Delete(id string) (*WorkerInfo, error) {

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

func (oe *imlAdmin) GetEmployee(profession, name string) (*WorkerInfo, error) {

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

func (oe *imlAdmin) set(id, profession, name, driverName, version, desc string, body []byte, updateAt, createAt string) (*WorkerInfo, error) {

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
		wInfo.reset(driverName, version, desc, body, wInfo.worker, driver.ConfigType(), updateAt, createAt)
		oe.variables.SetRequire(id, usedVariables)
		return wInfo, nil
	}
	// create
	worker, err := driver.Create(id, name, conf, requires)
	if err != nil {
		log.Warn("worker-data set worker create:", err)
		return nil, err
	}

	if !hasInfo {
		if updateAt == "" {
			updateAt = eosc.Now()
		}
		if createAt == "" {
			createAt = eosc.Now()
		}
		wInfo = NewWorkerInfo(worker, id, profession, name, driverName, version, desc, createAt, updateAt, body, driver.ConfigType())
	} else {
		if updateAt == "" {
			updateAt = eosc.Now()
		}

		createAt = wInfo.config.Create

		wInfo.reset(driverName, version, desc, body, worker, driver.ConfigType(), updateAt, createAt)
	}

	// store
	oe.data.Set(id, wInfo)
	oe.requireManager.Set(id, getIds(requires))
	oe.variables.SetRequire(id, usedVariables)
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
