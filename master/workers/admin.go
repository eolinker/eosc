package workers

import (
	"encoding/json"
	"fmt"

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

func (w *Workers) Delete(id string) (*admin.WorkerInfo, error) {
	o, has := w.data.Get(id)
	if !has {
		return nil, admin.ErrorWorkerNotExist
	}

	err := w.service.Send(SpaceWorker, CommandDel, []byte(id))
	if err != nil {
		return nil, err
	}
	worker := o.(*Worker)
	return &admin.WorkerInfo{
		Id:     worker.Id,
		Name:   worker.Name,
		Driver: worker.Driver,
		Create: worker.CreateTime,
		Update: worker.UpdateTime,
	}, nil
}

func (w *Workers) Set(profession, name, driver string, data []byte) error {
	d := &WorkerData{
		Id:         fmt.Sprintf("%s@%s", name, profession),
		Profession: profession,
		Name:       name,
		Driver:     driver,
		CreateTime: "",
		UpdateTime: "",
		Sing:       "",
		Data:       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return err
	}
	return w.service.Send(SpaceWorker, name, body)
}
