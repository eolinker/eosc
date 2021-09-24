package workers

import "github.com/eolinker/eosc"

type ITypedWorkers interface {
	All() []*Worker
	Set(id string, w *Worker)
	Get(id string) (*Worker, bool)
	Del(id string) (*Worker, bool)
}

type TypedWorkers struct {
	data eosc.IUntyped
}

func NewTypedWorkers() *TypedWorkers {
	return &TypedWorkers{
		data: eosc.NewUntyped(),
	}
}

func (t *TypedWorkers) All() []*Worker {
	panic("implement me")
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
