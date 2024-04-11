package admin

import (
	"context"
	"fmt"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
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

func (oe *imlAdminApi) SetWorker(ctx context.Context, wc *eosc.WorkerConfig) (*WorkerInfo, error) {

	id, ok := eosc.ToWorkerId(wc.Name, wc.Profession)
	if !ok {
		return nil, fmt.Errorf("%s@%s:invalid id", wc.Name, wc.Profession)
	}

	log.Debug("update:", id, " ", wc.Profession, ",", wc.Name, ",", wc.Driver, ",", wc.Body)
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
