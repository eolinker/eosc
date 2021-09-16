package workers

import (
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/admin"
)

func (w *Workers) GetWork(id string) (admin.TWorker, error) {
	panic("implement me")
}

func (w *Workers) GetList(profession string) ([]admin.TWorker, error) {
	panic("implement me")
}

func (w *Workers) CheckerSkill(id string, skill string) (bool, error) {
	panic("implement me")
}

func (w *Workers) Delete(id string) (*admin.WorkerInfo, bool) {
	panic("implement me")
}

func (w *Workers) Set(profession, name, driver string, data eosc.IData) error {
	panic("implement me")
}
