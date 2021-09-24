package workers

import (
	"encoding/json"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/admin"
)

func (w *Workers) GetWork(id string) (eosc.TWorker, error) {
	if ow, b := w.data.Get(id); b {
		return ow.Format(nil), nil
	}
	return nil, ErrorNotExist
}

func (w *Workers) GetList(profession string) ([]eosc.TWorker, error) {
	p, has := w.professions.GetProfession(profession)
	if !has {
		return nil, ErrorInvalidProfession
	}
	attrs := p.AppendAttr()

	all := w.data.All()

	result := make([]eosc.TWorker, 0, len(all))
	for _, ow := range all {

		if ow.Profession == profession {
			result = append(result, ow.Format(attrs))
		}
	}
	return result, nil
}
func (w *Workers) Delete(id string) (*eosc.WorkerInfo, error) {
	worker, has := w.data.Get(id)
	if !has {
		return nil, admin.ErrorWorkerNotExist
	}

	err := w.service.Send(SpaceWorker, CommandDel, []byte(id))
	if err != nil {
		return nil, err
	}
	return &eosc.WorkerInfo{
		Id:         worker.Id,
		Profession: worker.Profession,
		Name:       worker.Name,
		Driver:     worker.Driver,
		Create:     worker.CreateTime,
		Update:     worker.UpdateTime,
	}, nil
}

func (w *Workers) Set(profession, name, driver string, data []byte) error {

	id := eosc.ToWorkerId(name, profession)

	createTime := eosc.Now()
	if ow, has := w.data.Get(id); has {
		if ow.Driver != driver {
			return ErrorChangeDriver
		}
		createTime = ow.CreateTime
	} else {
		pf, b := w.professions.GetProfession(profession)
		if !b {
			return ErrorInvalidProfession
		}
		if !pf.HasDriver(driver) {
			return ErrorInvalidDriver
		}
	}

	d := &eosc.WorkerData{
		Id:         id,
		Profession: profession,
		Name:       name,
		Driver:     driver,
		Create:     createTime,
		Update:     eosc.Now(),
		Body:       data,
	}

	body, err := json.Marshal(d)
	if err != nil {
		return err
	}
	return w.service.Send(SpaceWorker, name, body)
}
