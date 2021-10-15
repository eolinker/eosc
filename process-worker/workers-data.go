package process_worker

import (
	"github.com/eolinker/eosc"
)

var _ ITypedWorkers = (*TypedWorkers)(nil)

type ITypedWorkers interface {
	Set(id string, w *Worker)
	Del(id string) (*Worker, bool)
	Get(id string) (*Worker, bool)
}

type TypedWorkers struct {
	data eosc.IUntyped
}

func NewTypedWorkers() *TypedWorkers {
	return &TypedWorkers{
		data: eosc.NewUntyped(),
	}
}

func (wd *TypedWorkers) Set(id string, w *Worker) {
	wd.data.Set(id, w)
}

func (wd *TypedWorkers) Del(id string) (*Worker, bool) {
	worker, has := wd.data.Del(id)
	if !has {
		return nil, false
	}
	return worker.(*Worker), true
}

func (wd *TypedWorkers) Get(id string) (*Worker, bool) {
	worker, has := wd.data.Get(id)
	if !has {
		return nil, false
	}
	return worker.(*Worker), true
}
