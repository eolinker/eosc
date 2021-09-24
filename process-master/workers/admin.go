package workers

import (
	"encoding/json"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/admin"
)

func (w *Workers) GetWork(id string) (admin.TWorker, error) {
	if o, b := w.data.Get(id); b {
		ow, ok := o.(*Worker)
		if !ok {
			return nil, ErrorUnknown
		}
		p, has := w.professions.GetProfession(ow.Profession)
		if !has {
			return nil, ErrorUnknown
		}
		return ow.Format(p.AppendAttr()), nil
	}
	return nil, ErrorNotExist
}

func (w *Workers) GetList(profession string) ([]admin.TWorker, error) {
	p, has := w.professions.GetProfession(profession)
	if !has {
		return nil, ErrorInvalidProfession
	}
	attrs := p.AppendAttr()

	all := w.data.List()

	result := make([]admin.TWorker, 0, len(all))
	for _, o := range all {
		ow, ok := o.(*Worker)
		if !ok {
			continue
		}
		if ow.Profession == profession {
			result = append(result, ow.Format(attrs))
		}
	}
	return result, nil
}
func (w *Workers) Delete(id string) (*eosc.WorkerInfo, error) {
	o, has := w.data.Get(id)
	if !has {
		return nil, admin.ErrorWorkerNotExist
	}

	err := w.service.Send(SpaceWorker, CommandDel, []byte(id))
	if err != nil {
		return nil, err
	}
	worker := o.(*Worker)
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
	if o, has := w.data.Get(id); has {
		ow, ok := o.(*Worker)
		if !ok {
			return ErrorUnknown
		}
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
