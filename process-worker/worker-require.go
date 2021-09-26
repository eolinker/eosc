package process_worker

import (
	"sync"

	"github.com/eolinker/eosc"
)

var (
	_ IWorkerRequireManager = (*WorkerRequireManager)(nil)
)

type IWorkerRequireManager interface {
	Set(id string, requires map[eosc.RequireId]interface{})
	Del(id string)
	RequireByCount(requireId string) int
}

type WorkerRequireManager struct {
	locker    sync.Mutex
	requireBy eosc.IUntyped
	requires  eosc.IUntyped
}

func (w *WorkerRequireManager) Set(id string, requires map[eosc.RequireId]interface{}) {
	w.locker.Lock()
	defer w.locker.Unlock()
	w.del(id)
	if len(requires) > 0 {
		ids := make([]string, len(requires))
		for rid := range requires {
			ridStr := string(rid)
			ids = append(ids, ridStr)
			d, has := w.requireBy.Get(ridStr)
			if !has {
				w.requireBy.Set(ridStr, []string{id})
				continue
			} else {
				w.requireBy.Set(ridStr, append(d.([]string), id))
			}

		}
		w.requires.Set(id, ids)
	}
}

func (w *WorkerRequireManager) Del(id string) {
	w.locker.Lock()
	defer w.locker.Unlock()

	w.del(id)

}
func (w *WorkerRequireManager) del(id string) {
	if r, has := w.requires.Del(id); has {
		rs := r.([]string)
		for _, rid := range rs {
			w.removeBy(id, rid)
		}
	}
}

func (w *WorkerRequireManager) removeBy(id string, requireId string) {
	if d, has := w.requireBy.Get(requireId); has {
		rs := d.([]string)
		for i, rid := range rs {
			if rid == id {
				rs = append(rs[:i], rs[i+1])
				break
			}
		}
		if len(rs) == 0 {
			w.requireBy.Del(requireId)
		} else {
			w.requireBy.Set(requireId, rs)
		}
	}
}

func (w *WorkerRequireManager) RequireByCount(requireId string) int {
	if rs, has := w.requireBy.Get(requireId); has {
		return len(rs.([]string))
	}
	return 0
}
