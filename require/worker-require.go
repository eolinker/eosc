package require

import (
	"sync"

	"github.com/eolinker/eosc"
)

var (
	_ eosc.IRequires = (*Manager)(nil)
)

type Manager struct {
	locker    sync.Mutex
	requireBy eosc.IUntyped
	workerIds eosc.IUntyped
}

func NewRequireManager() eosc.IRequires {
	return &Manager{
		locker:    sync.Mutex{},
		requireBy: eosc.NewUntyped(),
		workerIds: eosc.NewUntyped(),
	}
}

func (w *Manager) Set(id string, requiresIds []string) {
	w.locker.Lock()
	defer w.locker.Unlock()
	w.del(id)
	if len(requiresIds) > 0 {

		for _, rid := range requiresIds {
			d, has := w.requireBy.Get(rid)
			if !has {
				w.requireBy.Set(rid, []string{id})
			} else {
				w.requireBy.Set(rid, append(d.([]string), id))
			}
		}
		w.workerIds.Set(id, requiresIds)
	}
}
func (w *Manager) RequireBy(requireId string) []string {
	// 获取依赖requireID的所有id列表
	w.locker.Lock()
	ids, has := w.requireBy.Get(requireId)
	w.locker.Unlock()
	if has {
		return ids.([]string)
	}
	return nil
}
func (w *Manager) Del(id string) {
	w.locker.Lock()
	w.del(id)
	w.locker.Unlock()

}
func (w *Manager) del(id string) {
	if r, has := w.workerIds.Del(id); has {
		rs := r.([]string)
		for _, rid := range rs {
			w.removeBy(id, rid)
		}
	}
}

func (w *Manager) removeBy(id string, requireId string) {
	if d, has := w.requireBy.Get(requireId); has {
		rs := d.([]string)
		for i, rid := range rs {
			if rid == id {
				rs = append(rs[:i], rs[i+1:]...)
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
func (w *Manager) Requires(id string) []string {
	// 根据id获取所有依赖的id列表
	w.locker.Lock()
	ids, has := w.workerIds.Get(id)
	w.locker.Unlock()
	if has {
		return ids.([]string)
	}
	return nil
}

func (w *Manager) RequireByCount(requireId string) int {
	w.locker.Lock()
	rs, has := w.requireBy.Get(requireId)
	w.locker.Unlock()
	if has {
		return len(rs.([]string))
	}
	return 0
}
