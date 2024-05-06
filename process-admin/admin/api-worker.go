package admin

import (
	"context"
	"fmt"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/service"
	"github.com/eolinker/eosc/utils"
)

func (d *imlAdminData) CheckDeleteWorker(ids ...string) (requires []string) {
	for _, id := range ids {
		if d.requireManager.RequireByCount(id) > 0 {
			requires = append(requires, id)
		}
	}
	return requires
}

func (d *imlAdminData) AllWorkers(ctx context.Context) []*WorkerInfo {
	return d.workers.List()
}
func (d *imlAdminData) ListWorker(ctx context.Context, profession string) ([]*WorkerInfo, error) {
	list := d.workers.List()
	return utils.ArrayFilter(list, func(i int, v *WorkerInfo) bool {
		return v.config.Profession == profession
	}), nil
}

func (d *imlAdminData) GetWorker(ctx context.Context, id string) (*WorkerInfo, bool) {
	return d.workers.Get(id)
}
func (oe *imlAdminApi) DeleteWorker(ctx context.Context, id string) (*WorkerInfo, error) {

	worker, err := oe.imlAdminData.delWorker(id)
	if err != nil {
		return nil, err
	}
	oe.actions = append(oe.actions, newRollbackForDelete(worker))
	oe.events = append(oe.events, &service.Event{
		Command:   eosc.EventDel,
		Namespace: eosc.NamespaceWorker,
		Key:       id,
		Data:      nil,
	})
	return worker, nil

}

func (oe *imlAdminApi) SetWorker(ctx context.Context, wc *eosc.WorkerConfig) (*WorkerInfo, error) {

	id, ok := eosc.ToWorkerId(wc.Name, wc.Profession)
	if !ok {
		return nil, fmt.Errorf("%s@%s:invalid id", wc.Name, wc.Profession)
	}
	wc.Id = id
	log.Debug("update:", wc.Id, " ", wc.Profession, ",", wc.Name, ",", wc.Driver, ",", string(wc.Body))
	old, exits := oe.imlAdminData.workers.Get(id)
	if exits {
		// update
		if wc.Driver == "" {
			wc.Driver = old.Driver()
		}

		info, err := oe.imlAdminData.setWorker(wc)
		if err != nil {
			return nil, err
		}
		oe.actions = append(oe.actions, newRollBackForSet(old.config))
		oe.events = append(oe.events, &service.Event{
			Command:   eosc.EventSet,
			Namespace: eosc.NamespaceWorker,
			Key:       id,
			Data:      info.ConfigData(),
		})
		return info, nil
	}
	// create
	if wc.Driver == "" {
		return nil, fmt.Errorf("require driver")
	}

	info, err := oe.imlAdminData.setWorker(wc)
	if err != nil {
		return nil, err
	}
	oe.actions = append(oe.actions, newRollBackForCreate(id))
	oe.events = append(oe.events, &service.Event{
		Command:   eosc.EventSet,
		Namespace: eosc.NamespaceWorker,
		Key:       id,
		Data:      info.ConfigData(),
	})
	return info, nil
}
