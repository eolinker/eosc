package workers

import "github.com/eolinker/eosc"

type ITypedWorkers interface {
	All() []*Worker
	Set(id string, w *Worker)
	Get(id string) (*Worker, bool)
	Del(id string) (*Worker, bool)
	Reset(ds []*eosc.WorkerData)
}

type TypedWorkers struct {
	data eosc.IUntyped
}

func (t *TypedWorkers) Reset(ds []*eosc.WorkerData) {
	for _, key := range t.data.Keys() {
		t.data.Del(key)
	}
	for _, v := range ds {
		wv, err := ToWorker(v)
		if err != nil {
			continue
		}
		t.data.Set(v.Id, wv)
	}
}

func NewTypedWorkers() ITypedWorkers {
	return &TypedWorkers{
		data: eosc.NewUntyped(),
	}
}

func (t *TypedWorkers) All() []*Worker {
	vs := t.data.List()
	rs := make([]*Worker, len(vs))
	for i, v := range vs {
		rs[i] = v.(*Worker)
	}
	return rs
}

func (t *TypedWorkers) Set(id string, w *Worker) {
	t.data.Set(id, w)
}

func (t *TypedWorkers) Get(id string) (*Worker, bool) {
	w, h := t.data.Get(id)
	if h {
		return w.(*Worker), true
	}
	return nil, false
}

func (t *TypedWorkers) Del(id string) (*Worker, bool) {
	w, h := t.data.Del(id)
	if h {
		return w.(*Worker), true
	}
	return nil, false
}
