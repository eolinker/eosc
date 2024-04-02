package admin

import (
	"context"
	"encoding/json"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	open_api "github.com/eolinker/eosc/open-api"
	"github.com/eolinker/eosc/professions"
	"github.com/eolinker/eosc/require"
	"github.com/eolinker/eosc/setting"
	"github.com/eolinker/eosc/utils"
	"sync"
)

var (
	_ AdminController = (*imlAdminData)(nil)
	_ eosc.IWorkers   = (*imlAdminData)(nil)
	_ iAdminOperator  = (*imlAdminData)(nil)
)

type imlAdminData struct {
	transactionLocker sync.Mutex
	workers           eosc.Untyped[string, *WorkerInfo]
	professionData    professions.IProfessions
	variable          eosc.IVariable
	settings          eosc.ISettings
	requireManager    eosc.IRequires
}

func (d *imlAdminData) setProfession(name string, profession *eosc.ProfessionConfig) error {
	return d.professionData.Set(name, profession)
}

func (d *imlAdminData) delProfession(name string) error {
	return d.professionData.Delete(name)
}

func (d *imlAdminData) ListWorker(ctx context.Context, profession string) ([]*WorkerInfo, error) {
	list := d.workers.List()
	return utils.ArrayFilter(list, func(i int, v *WorkerInfo) bool {
		return v.config.Profession == profession
	}), nil
}

func (d *imlAdminData) GetWorker(ctx context.Context, id string) (*WorkerInfo, error) {
	workerInfo, has := d.workers.Get(id)
	if has {
		return workerInfo, nil
	}
	return nil, ErrorNotExist
}

func NewImlAdminData(initData map[string][]byte, professionData professions.IProfessions, variable eosc.IVariable) AdminController {

	workerData := utils.MapFilter(initData, func(k string, v []byte) bool {
		profession, _, _ := eosc.SplitWorkerId(k)
		return profession != Setting
	})

	data := &imlAdminData{
		professionData: professionData,
		variable:       variable,
		workers:        eosc.BuildUntyped[string, *WorkerInfo](),
		settings:       setting.GetSettings(),
		requireManager: require.NewRequireManager(),
	}
	for id, d := range workerData {
		cf := new(eosc.WorkerConfig)
		e := json.Unmarshal(d, cf)
		if e != nil {
			continue
		}
		data.workers.Set(id, &WorkerInfo{
			worker: nil,
			config: cf,
			attr:   nil,
			info:   nil,
		})
	}
	settingData := utils.MapFilter(initData, func(k string, v []byte) bool {
		profession, _, _ := eosc.SplitWorkerId(k)
		return profession == Setting
	})
	for id, conf := range settingData {

		_, name, _ := eosc.SplitWorkerId(id)
		_, has := data.settings.GetDriver(name)
		log.Debug("init setting id: ", id, " conf: ", string(conf), " ", has)
		if has {
			config := new(eosc.WorkerConfig)
			err := json.Unmarshal(conf, config)
			if err != nil {
				log.Warn("init setting Unmarshal WorkerConfig:", err)
				continue
			}
			log.Debug("init setting id body: ", id, " conf: ", string(config.Body), " ", has)
			err = data.settings.SettingWorker(name, config.Body, variable)
			if err != nil {
				log.Warn("init setting:", err)
			}
		}
	}
	return data
}
func (d *imlAdminData) GetProfession(ctx context.Context, profession string) (*professions.Profession, bool) {
	return d.professionData.Get(profession)
}

func (d *imlAdminData) ListProfession(ctx context.Context) []*professions.Profession {
	return d.professionData.List()
}

func (d *imlAdminData) GetSetting(ctx context.Context, name string) (any, bool) {
	_, has := d.settings.GetDriver(name)
	if !has {
		return nil, false
	}
	return d.settings.GetConfig(name), has
}
func (d *imlAdminData) Get(id string) (eosc.IWorker, bool) {
	w, has := d.workers.Get(id)
	if has {
		return w.worker, true
	}
	return nil, false
}

func (d *imlAdminData) unLock() {
	d.transactionLocker.Unlock()
}

func (d *imlAdminData) Transaction(ctx context.Context, f func(ctx context.Context, api AdminApiWrite) error) ([]*open_api.EventResponse, error) {

	adminTransaction := d.Begin(ctx)

	err := f(ctx, adminTransaction)
	if err != nil {
		rollbackError := adminTransaction.Rollback()
		if rollbackError != nil {
			log.Error("rollback error:", rollbackError)
		}
		return nil, err
	}
	return adminTransaction.Commit()
}

func (d *imlAdminData) AllVariables(ctx context.Context) map[string]map[string]string {
	return d.variable.All()
}

func (d *imlAdminData) GetVariables(ctx context.Context, namespace string) (map[string]string, bool) {
	values, has := d.variable.GetByNamespace(namespace)
	return values, has
}

func (d *imlAdminData) GetVariable(ctx context.Context, namespace, key string) (string, bool) {
	values, has := d.variable.GetByNamespace(namespace)
	if !has {
		return "", false
	}
	v, h := values[key]
	return v, h
}
func (d *imlAdminData) AllWorkers(ctx context.Context) []*WorkerInfo {
	//return utils.GroupBy(a.workers.List(), func(v *WorkerInfo) string {
	//	return v.config.Profession
	//})

	return d.workers.List()
}

func (d *imlAdminData) Begin(ctx context.Context) AdminTransaction {
	d.transactionLocker.Lock()
	return newImlAdminApi(d)
}
func (d *imlAdminData) CheckDelete(ids ...string) (requires []string) {
	for _, id := range ids {
		if d.requireManager.RequireByCount(id) > 0 {
			requires = append(requires, id)
		}
	}
	return requires
}
