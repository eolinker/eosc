package eosc

var _ IWorkers = (*WorkManager)(nil)
var _ iWorkData = (*Workers)(nil)

type iWorkData interface {
	Set(id string, w IWorker)
	Del(id string) (*tWorker, bool)
	Get(id string) (*tWorker, bool)
}

type IWorkers interface {
	Set(id string, w IWorker)
	Del(id string) (IWorker, bool)
	Get(id string) (IWorker, bool)
}
type WorkManager struct {
	data        Workers
}

func (wm *WorkManager) Set(id string, w IWorker) {
 	wm.data.Set(id,w)
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



func NewWorkers() IWorkers{

	ws := &WorkManager{
		//store:       store,
 		data:        Workers{data: NewUntyped()},
	}

	return ws
}



type Workers struct {
	data IUntyped
}

func (ws *Workers) Set(id string, w IWorker) {
	wk:=newTWorker(w)
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
