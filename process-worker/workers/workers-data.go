package workers

import (
	"github.com/eolinker/eosc"
)

type WorkerDatas struct {
	data eosc.IUntyped
}

func (wd *WorkerDatas) All() []eosc.IWorker {
	all := wd.data.List()
	list := make([]eosc.IWorker, 0, len(all))
	for _, v := range all {
		list = append(list, v.(eosc.IWorker))
	}
	return list
}

func NewTypedWorkers() *WorkerDatas {
	return &WorkerDatas{
		data: eosc.NewUntyped(),
	}
}

func (wd *WorkerDatas) Set(id string, w eosc.IWorker) {
	wd.data.Set(id, w)
}

func (wd *WorkerDatas) Del(id string) (eosc.IWorker, bool) {
	worker, has := wd.data.Del(id)
	if !has {
		return nil, false
	}

	return worker.(eosc.IWorker), true
}

func (wd *WorkerDatas) Get(id string) (eosc.IWorker, bool) {
	worker, has := wd.data.Get(id)
	if !has {
		return nil, false
	}
	return worker.(eosc.IWorker), true
}
