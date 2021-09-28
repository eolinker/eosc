package process_master

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/process-master/workers"
	raft_service "github.com/eolinker/eosc/raft/raft-service"
	"github.com/eolinker/eosc/service"
	"github.com/eolinker/eosc/utils"
	"github.com/golang/protobuf/proto"
)

type WorkersData struct {
	workers.ITypedWorkers
}

func NewWorkersData(ITypedWorkers workers.ITypedWorkers) *WorkersData {
	return &WorkersData{ITypedWorkers: ITypedWorkers}
}
func (w *WorkersData) Encode(startIndex int) ([]byte, []*os.File, error) {
	data, err := w.encode()
	if err != nil {
		return nil, nil, err
	}
	return utils.EncodeFrame(data), nil, nil
}
func (w *WorkersData) encode() ([]byte, error) {
	values := w.ITypedWorkers.All()

	wd := &eosc.WorkersData{
		Data: make([]*eosc.WorkerData, len(values)),
	}
	for i, v := range values {

		wd.Data[i] = v.Org
	}
	data, err := proto.Marshal(wd)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (w *WorkersData) decode(data []byte) ([]*eosc.WorkerData, error) {
	wd := new(eosc.WorkersData)
	err := proto.Unmarshal(data, wd)
	if err != nil {
		return nil, err
	}
	return wd.Data, nil
}
func (w *WorkersData) reset(vs []*eosc.WorkerData) {
	w.ITypedWorkers.Reset(vs)
}

type WorkersRaft struct {
	data                    *WorkersData
	professions             eosc.IProfessionsData
	workerServiceClient     service.WorkerServiceClient
	service                 raft_service.IService
	workerProcessController WorkerProcessController
}

func NewWorkersRaft(workerData *WorkersData, professions eosc.IProfessionsData, workerServiceClient service.WorkerServiceClient, service raft_service.IService) *WorkersRaft {
	return &WorkersRaft{data: workerData, professions: professions, workerServiceClient: workerServiceClient, service: service}
}

func (w *WorkersRaft) Delete(id string) error {

	err := w.service.Send(workers.SpaceWorker, workers.CommandSet, []byte(id))
	if err != nil {
		return err
	}
	return nil
}

func (w *WorkersRaft) Set(profession, name, driver string, data []byte) error {
	id := eosc.ToWorkerId(name, profession)

	createTime := eosc.Now()
	if ow, has := w.data.Get(id); has {
		if ow.Driver != driver {
			return workers.ErrorChangeDriver
		}
		createTime = ow.CreateTime
	} else {
		pf, b := w.professions.GetProfession(profession)
		if !b {
			return workers.ErrorInvalidProfession
		}
		if !pf.HasDriver(driver) {
			return workers.ErrorInvalidDriver
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
	return w.service.Send(workers.SpaceWorker, name, body)
}

func (w *WorkersRaft) ProcessHandler(cmd string, body []byte) ([]byte, error) {

	switch cmd {
	case workers.CommandSet:
		request := &service.WorkerSetRequest{
			Body: body,
		}

		response, err := w.workerServiceClient.SetCheck(context.TODO(), request)
		if err != nil {
			return nil, err
		}
		if response.Status != service.WorkerStatusCode_SUCCESS {
			return nil, errors.New(response.Message)
		}
		return body, nil
	case workers.CommandDel:
		request := &service.WorkerDeleteRequest{
			Id: string(body),
		}

		response, err := w.workerServiceClient.DeleteCheck(context.TODO(), request)
		if err != nil {
			return nil, err
		}
		if response.Status != service.WorkerStatusCode_SUCCESS {
			return nil, errors.New(response.Message)
		}
		return body, nil
	}
	return nil, workers.ErrorInvalidCommand

}

func (w *WorkersRaft) CommitHandler(cmd string, data []byte) error {
	switch cmd {
	case workers.CommandSet:
		{
			worker, err := workers.DecodeWorker(data)
			if err != nil {
				return err
			}
			w.data.Set(worker.Id, worker)
			req := &service.WorkerSetRequest{
				Id:         worker.Id,
				Profession: worker.Profession,
				Name:       worker.Name,
				Driver:     worker.Driver,
				Body:       worker.Org.Body,
			}
			response, err := w.workerServiceClient.Set(context.TODO(), req)
			if err != nil {
				return err
			}
			if response.Status != service.WorkerStatusCode_SUCCESS {
				return errors.New(response.Message)
			}

			return nil
		}
	case workers.CommandDel:
		{
			id := string(data)
			w.data.Del(id)
			return nil
		}
	default:
		return raft_service.ErrInvalidCommand
	}
}

func (w *WorkersRaft) Snapshot() []byte {

	data, err := w.data.encode()
	if err != nil {
		return nil
	}
	return data
}

func (w *WorkersRaft) ResetHandler(data []byte) error {

	vs, err := w.data.decode(data)
	if err != nil {
		return err
	}

	w.data.reset(vs)

	return nil
}

func (w *WorkersRaft) GetWork(id string) (eosc.TWorker, error) {
	if ow, b := w.data.Get(id); b {
		return ow.Format(nil), nil
	}
	return nil, workers.ErrorNotExist
}

func (w *WorkersRaft) GetList(profession string) ([]eosc.TWorker, error) {
	p, has := w.professions.GetProfession(profession)
	if !has {
		return nil, workers.ErrorInvalidProfession
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
