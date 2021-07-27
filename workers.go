package eosc

import (
	"fmt"

	"github.com/eolinker/eosc/log"
)

var _ IWorkers = (*WorkManager)(nil)
var _ iWorkData = (*Workers)(nil)

type iWorkData interface {
	Set(id string, w *tWorker)
	Del(id string) (*tWorker, bool)
	Get(id string) (*tWorker, bool)
}

type IWorkers interface {
	//Set(id string, w IWorker)
	Del(id string) (IWorker, bool)
	Get(id string) (IWorker, bool)
	//SetWorker(v StoreValue) error
}
type WorkManager struct {
	data        Workers
	professions IProfessions
	store       IStore
}

func (wm *WorkManager) Del(id string) (IWorker, bool) {

	if w, has := wm.data.Del(id); has {
		return w.worker, true
	}
	return nil, false
}

func (wm *WorkManager) Get(id string) (IWorker, bool) {
	if w, has := wm.data.Get(id); has {
		return w.worker, true
	}
	return nil, false
}

func (wm *WorkManager) SetWorker(v *StoreValue) error {
	log.Debugf("set worker:%s", v.Id)
	p, has := wm.professions.get(v.Profession)
	if !has {
		return fmt.Errorf("%s:%w", v.Profession, ErrorProfessionNotExist)
	}
	d, has := p.getDriver(v.Driver)
	if !has {
		return fmt.Errorf("%s in %s:%w", v.Driver, v.Profession, ErrorDriverNotExist)
	}

	config := newConfig(d.ConfigType())

	err := v.IData.UnMarshal(&config)
	if err != nil {
		return err
	}

	requires, err := CheckConfig(config, wm)
	if err != nil {
		return err
	}
	worker, err := d.Create(v.Id, v.Name, config, requires)

	if err != nil {
		return nil
	}

	w, has := wm.data.Get(v.Id)
	if has {
		err := w.worker.Reset(config, requires)
		if err != nil {
			return err
		}
	} else {
		err := worker.Start()
		if err != nil {
			return err
		}
		w = newTWorker(worker)
		wm.data.Set(v.Id, w)
	}

	return nil
}

func (wm *WorkManager) OnChange(v StoreValue) error {

	return wm.SetWorker(&v)
}

func (wm *WorkManager) OnDel(v StoreValue) error {

	w, has := wm.data.Del(v.Id)
	if has {
		return w.worker.Stop()
	}
	return nil
}

func (wm *WorkManager) OnInit(vs []StoreValue) error {

	wdata := make(map[string][]*StoreValue)
	for _, v := range vs {
		vt := v
		wdata[v.Profession] = append(wdata[v.Profession], &vt)
	}
	for _, pi := range wm.professions.Infos() {

		_, has := wm.professions.get(pi.Name)
		if !has {
			return fmt.Errorf("%s:%w", pi.Name, ErrorProfessionNotExist)
		}

		for _, v := range wdata[pi.Name] {
			err := wm.SetWorker(v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NewWorkers(professions IProfessions, store IStore) (*WorkManager, error) {

	ws := &WorkManager{
		store:       store,
		professions: professions,
		data:        Workers{data: NewUntyped()},
	}
	err := ws.init()
	if err != nil {
		return nil, err
	}

	return ws, nil
}

func (wm *WorkManager) init() error {

	return wm.store.GetListener().AddListen(wm)
}

type Workers struct {
	data IUntyped
}

func (ws *Workers) Set(id string, w *tWorker) {

	ws.data.Set(id, w)
}

func (ws *Workers) Del(id string) (*tWorker, bool) {

	o, has := ws.data.Del(id)
	if has {
		w, ok := o.(*tWorker)
		return w, ok
	}
	return nil, false
}

func (ws *Workers) Get(id string) (*tWorker, bool) {
	o, has := ws.data.Get(id)
	if has {
		w, ok := o.(*tWorker)
		return w, ok
	}
	return nil, false
}
