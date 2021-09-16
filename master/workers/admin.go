package workers

import (
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/admin"
)

func (w *Worker) GetWork(id string) (admin.TWorker, error) {
	panic("implement me")
}

func (w *Worker) GetList(profession string) ([]admin.TWorker, error) {
	panic("implement me")
}

func (w *Worker) CheckerSkill(id string, skill string) (bool, error) {
	panic("implement me")
}

func (w *Worker) Delete(id string) (*admin.WorkerInfo, bool) {
	panic("implement me")
}

func (w *Worker) Set(profession, name, driver string, data eosc.IData) error {
	panic("implement me")
}
