package workers

import (
	"github.com/eolinker/eosc"
)

type WorkerDatas struct {
	data eosc.Untyped[string, eosc.IWorker]
}

func (wd *WorkerDatas) All() []eosc.IWorker {
	all := wd.data.List()
	return all
}

func NewTypedWorkers() *WorkerDatas {
	return &WorkerDatas{
		data: eosc.BuildUntyped[string, eosc.IWorker](),
	}
}

func (wd *WorkerDatas) Set(id string, w eosc.IWorker) {
	wd.data.Set(id, w)
}

func (wd *WorkerDatas) Del(id string) (eosc.IWorker, bool) {
	return wd.data.Del(id)
}

func (wd *WorkerDatas) Get(id string) (eosc.IWorker, bool) {
	return wd.data.Get(id)

}
