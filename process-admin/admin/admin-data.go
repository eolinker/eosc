package admin

import (
	"context"
	"encoding/json"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/professions"
	"github.com/eolinker/eosc/require"
	"github.com/eolinker/eosc/setting"
	"github.com/eolinker/eosc/utils"
	"github.com/eolinker/eosc/utils/hash"
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

	customerHash eosc.Untyped[string, stringHash]
}

func NewImlAdminData(workerInitData map[string][]byte, professionData professions.IProfessions, variable eosc.IVariable, hashInitData map[string][]byte) AdminController {

	workerData := utils.MapFilter(workerInitData, func(k string, v []byte) bool {
		profession, _, _ := eosc.SplitWorkerId(k)
		return profession != Setting
	})

	data := &imlAdminData{
		professionData: professionData,
		variable:       variable,
		workers:        eosc.BuildUntyped[string, *WorkerInfo](),
		settings:       setting.GetSettings(),
		requireManager: require.NewRequireManager(),
		customerHash:   eosc.BuildUntyped[string, stringHash](),
	}
	var i eosc.ICustomerVar = newImlCustomerHash(data.customerHash)
	bean.Injection(&i)
	if len(hashInitData) > 0 {
		for key, d := range hashInitData {
			v := make(map[string]string)
			err := json.Unmarshal(d, &v)
			if err != nil {
				continue
			}
			data.setHashValue(key, hash.NewHash(v))
		}
	}

	for _, d := range workerData {
		cf := new(eosc.WorkerConfig)
		e := json.Unmarshal(d, cf)
		if e != nil {
			continue
		}
		_, err := data.setWorker(cf)
		if err != nil {
			continue
		}
	}
	settingData := utils.MapFilter(workerInitData, func(k string, v []byte) bool {
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
			err = data.settings.SettingWorker(name, config.Body)
			if err != nil {
				log.Warn("init setting:", err)
			}
		}
	}
	return data
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

func (d *imlAdminData) Transaction(ctx context.Context, f func(ctx context.Context, api AdminApiWrite) error) error {

	adminTransaction := d.Begin(ctx)

	err := f(ctx, adminTransaction)
	if err != nil {
		rollbackError := adminTransaction.Rollback()
		if rollbackError != nil {
			log.Error("rollback error:", rollbackError)
		}
		return err
	}
	return adminTransaction.Commit()
}

func (d *imlAdminData) Begin(ctx context.Context) AdminTransaction {
	d.transactionLocker.Lock()
	return newImlAdminApi(d)
}
