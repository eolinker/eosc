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
	// 初始化 客户自定义hash数据
	if len(hashInitData) > 0 {
		for key, d := range hashInitData {
			v := make(map[string]string)
			err := json.Unmarshal(d, &v)
			if err != nil {
				log.Debug("init hash data error:", err)
				continue
			}
			data.setHashValue(key, hash.NewHash(v))
		}
	}
	type DataContext struct {
		Id   string
		Data []byte
	}
	wdc := utils.MapType(workerInitData, func(id string, d []byte) (*DataContext, bool) {
		return &DataContext{
			Id:   id,
			Data: d,
		}, true
	})
	workerInitDataByProfession := utils.MapReGroup(wdc, func(id string, v *DataContext) string {
		p, _, _ := eosc.SplitWorkerId(id)
		return p
	})
	// 初始化setting
	settingData := workerInitDataByProfession[Setting]
	for _, conf := range settingData {

		_, name, _ := eosc.SplitWorkerId(conf.Id)
		settingDriver, has := data.settings.GetDriver(name)
		log.Debug("init setting id: ", conf.Id, " conf: ", string(conf.Data), " ", has)
		if has && settingDriver.Mode() == eosc.SettingModeSingleton {
			config := new(eosc.WorkerConfig)
			err := json.Unmarshal(conf.Data, config)
			if err != nil {
				log.Warn("init setting Unmarshal WorkerConfig:", err)
				continue
			}
			log.Debug("init setting id body: ", conf.Id, " conf: ", string(config.Body), " ", has)
			err = data.settings.SettingWorker(name, config.Body)
			if err != nil {
				log.Warn("init setting:", err)
			}
		}
	}

	for _, profession := range professionData.Sort() {
		for _, conf := range workerInitDataByProfession[profession.Name] {
			cf := new(eosc.WorkerConfig)
			e := json.Unmarshal(conf.Data, cf)
			if e != nil {
				log.Debug("init setting Unmarshal WorkerConfig:", e)
				continue
			}
			_, err := data.setWorker(cf)
			if err != nil {
				log.Debug("init setting worker:", err)
				continue
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
