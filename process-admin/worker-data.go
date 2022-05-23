package process_admin

import (
	"github.com/eolinker/eosc"
)

type WorkerDatas struct {
	data eosc.IUntyped
}

func (w *WorkerDatas) Get(id string) (eosc.IWorker, bool) {
	info, has := w.GetInfo(id)
	if has {
		return info.worker, true
	}
	return nil, false
}

func NewWorkerDatas() *WorkerDatas {
	return &WorkerDatas{data: eosc.NewUntyped()}
}

func (w *WorkerDatas) Set(name string, v *WorkerInfo) {
	w.data.Set(name, v)
}

func (w *WorkerDatas) GetInfo(name string) (*WorkerInfo, bool) {
	v, has := w.data.Get(name)
	if has {
		return v.(*WorkerInfo), true
	}
	return nil, false
}

func (w *WorkerDatas) Del(name string) (*WorkerInfo, bool) {
	v, has := w.data.Del(name)
	if has {
		return v.(*WorkerInfo), true
	}
	return nil, false
}

func (w *WorkerDatas) List() []*WorkerInfo {
	list := w.data.List()
	rs := make([]*WorkerInfo, 0, len(list))
	for _, v := range list {
		rs = append(rs, v.(*WorkerInfo))
	}
	return rs
}

func (w *WorkerDatas) Keys() []string {
	return w.data.Keys()
}

func (w *WorkerDatas) All() map[string]*WorkerInfo {
	ds := w.data.All()
	m := make(map[string]*WorkerInfo)
	for k, v := range ds {
		m[k] = v.(*WorkerInfo)
	}
	return m
}

func (w *WorkerDatas) Clone() *WorkerDatas {
	return &WorkerDatas{data: w.data.Clone()}
}

func (w *WorkerDatas) Count() int {
	return w.data.Count()
}
