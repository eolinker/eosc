package process_admin

import (
	"encoding/json"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
)

type WorkerDatas struct {
	data eosc.Untyped[string, *WorkerInfo]
}

func (w *WorkerDatas) Get(id string) (eosc.IWorker, bool) {
	info, has := w.GetInfo(id)
	if has {
		return info.worker, true
	}
	return nil, false
}

func NewWorkerDatas(initData map[string][]byte) *WorkerDatas {
	data := &WorkerDatas{data: eosc.BuildUntyped[string, *WorkerInfo]()}
	for id, d := range initData {

		cf := new(eosc.WorkerConfig)
		e := json.Unmarshal(d, cf)
		if e != nil {
			continue
		}
		data.Set(id, &WorkerInfo{
			worker: nil,
			config: cf,
			attr:   nil,
			info:   nil,
		})
	}
	return data
}

func (w *WorkerDatas) Set(id string, v *WorkerInfo) {
	log.DebugF("worker set:%s==>%v", id, v.config)
	w.data.Set(id, v)
}

func (w *WorkerDatas) GetInfo(id string) (*WorkerInfo, bool) {
	return w.data.Get(id)
}

func (w *WorkerDatas) Del(id string) (*WorkerInfo, bool) {
	return w.data.Del(id)
}

func (w *WorkerDatas) List() []*WorkerInfo {
	return w.data.List()
}

func (w *WorkerDatas) Keys() []string {
	return w.data.Keys()
}

func (w *WorkerDatas) All() map[string]*WorkerInfo {
	ds := w.data.All()
	return ds
}
