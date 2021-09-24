package worker

import "github.com/eolinker/eosc"

type IWorkers = eosc.IWorkers
type IWorker = eosc.IWorker

var _ IWorkers = (*WorkersData)(nil)
var _ iWorkData = (*Workers)(nil)

type iWorkData interface {
	Set(id string, w eosc.IWorker)
	Del(id string) (*tWorker, bool)
	Get(id string) (*tWorker, bool)
}

type WorkersData struct {
	data Workers
}

func (wm *WorkersData) Set(id string, w IWorker) {
	wm.data.Set(id, w)
}

func (wm *WorkersData) Del(id string) (IWorker, bool) {

	if w, has := wm.data.Del(id); has {
		return w.worker, true
	}
	return nil, false
}

func (wm *WorkersData) Get(id string) (IWorker, bool) {
	if w, has := wm.data.Get(id); has {
		return w.worker, true
	}
	return nil, false
}

func NewWorkers() IWorkers {

	ws := &WorkersData{
		//store:       store,
		data: Workers{data: eosc.NewUntyped()},
	}

	return ws
}

type Workers struct {
	data eosc.IUntyped
}

func (ws *Workers) Set(id string, w eosc.IWorker) {
	wk := newTWorker(w)
	ws.data.Set(id, wk)
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
