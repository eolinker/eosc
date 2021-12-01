package workers

import (
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
)

type ITypedWorkers interface {
	All() []*Worker
	Set(id string, w *Worker)
	Get(id string) (*Worker, bool)
	Del(id string) (*Worker, bool)
	Reset(ds []*eosc.WorkerConfig)
}

type TypedWorkers struct {
	data eosc.IUntyped
}

func (t *TypedWorkers) Reset(ds []*eosc.WorkerConfig) {
	nw := eosc.NewUntyped()
	log.Debug("reset worker data len: ", len(ds))
	for _, v := range ds {
		log.Debug("reset worker data detail: ", *v)
		wv, err := ToWorker(v)
		if err != nil {
			continue
		}
		nw.Set(v.Id, wv)
	}
	t.data = nw
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
