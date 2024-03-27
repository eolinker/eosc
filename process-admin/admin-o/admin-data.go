package admin_o

import (
	"context"
	"encoding/json"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	open_api "github.com/eolinker/eosc/open-api"
	"github.com/eolinker/eosc/professions"
	"github.com/eolinker/eosc/require"
	"sync"
)

var (
	_ AdminController = (*imlAdminData)(nil)
	_ eosc.IWorkers   = (*imlAdminData)(nil)
)

type imlAdminData struct {
	transactionLocker sync.Mutex
	workers           eosc.Untyped[string, *WorkerInfo]
	professionData    professions.Professions
	variable          eosc.IVariable
	settings          eosc.ISettings
	requireManager    eosc.IRequires
}

func NewImlAdminData(initData map[string][]byte, professionData professions.Professions, variable eosc.IVariable, settings eosc.ISettings) AdminController {
	data := &imlAdminData{
		professionData: professionData,
		variable:       variable,
		settings:       settings,
		requireManager: require.NewRequireManager(),
	}
	for id, d := range initData {
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

func (d *imlAdminData) Transaction(ctx context.Context, f func(ctx context.Context, api AdminApi) error) ([]*open_api.EventResponse, error) {

	adminTransaction, err := d.Begin(ctx)
	if err != nil {
		return nil, err
	}
	err = f(ctx, adminTransaction)
	if err != nil {
		rollbackError := adminTransaction.Rollback()
		if rollbackError != nil {
			log.Error("rollback error:", rollbackError)
		}
		return nil, err
	}
	return adminTransaction.Commit()
}

func (d *imlAdminData) Begin(ctx context.Context) (AdminTransaction, error) {
	d.transactionLocker.Lock()
	return newImlAdminApi(d), nil
}
