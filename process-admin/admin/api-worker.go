package admin

import (
	"context"
	"fmt"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/process-admin/marshal"
	"github.com/eolinker/eosc/service"
)

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

func (oe *imlAdminApi) SetWorker(ctx context.Context, profession, name, driver, version, desc string, data marshal.IData) (*WorkerInfo, error) {
	body, err := data.Encode()
	if err != nil {
		return nil, err
	}
	id, ok := eosc.ToWorkerId(name, profession)
	if !ok {
		return nil, fmt.Errorf("%s@%s:invalid id", name, profession)
	}

	log.Debug("update:", id, " ", profession, ",", name, ",", driver, ",", body)
	old, exits := oe.imlAdminData.workers.Get(id)
	if exits {
		// update
		if driver == "" {
			driver = old.Driver()
		}

		info, err := oe.imlAdminData.setWorker(id, profession, name, driver, version, desc, body, eosc.Now(), old.config.Create)
		if err != nil {
			return nil, err
		}
		oe.actions = append(oe.actions, newRollBackForSet(old.config))
		oe.events = append(oe.events, &service.Event{
			Command:   eosc.EventSet,
			Namespace: eosc.NamespaceWorker,
			Key:       id,
			Data:      info.Body(),
		})
		return info, nil
	}
	// create
	if driver == "" {
		return nil, fmt.Errorf("require driver")
	}
	info, err := oe.imlAdminData.setWorker(id, profession, name, driver, version, desc, body, eosc.Now(), eosc.Now())
	if err != nil {
		return nil, err
	}
	oe.actions = append(oe.actions, newRollBackForCreate(id))
	oe.events = append(oe.events, &service.Event{
		Command:   eosc.EventSet,
		Namespace: eosc.NamespaceWorker,
		Key:       id,
		Data:      info.Body(),
	})
	return info, nil
}
