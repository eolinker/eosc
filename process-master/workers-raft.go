package process_master

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/eolinker/eosc/log"

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

		wd.Data[i] = v.WorkerData
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

func NewWorkersRaft(workerData *WorkersData, professions eosc.IProfessionsData, workerServiceClient service.WorkerServiceClient, service raft_service.IService, workerController WorkerProcessController) *WorkersRaft {
	return &WorkersRaft{data: workerData, professions: professions, workerServiceClient: workerServiceClient, service: service, workerProcessController: workerController}
}

func (w *WorkersRaft) Delete(id string) (eosc.TWorker, error) {

	obj, err := w.service.Send(workers.SpaceWorker, workers.CommandDel, []byte(id))
	if err != nil {
		return nil, err
	}

	worker, err := workers.ReadTWorker(obj)
	if err != nil {
		return nil, err
	}
	return worker, nil
}

func (w *WorkersRaft) Set(profession, name, driver string, data []byte) (eosc.TWorker, error) {
	id, ok := eosc.ToWorkerId(name, profession)
	if !ok {
		return nil, fmt.Errorf("%s %w", profession, errors.New("not match profession"))
	}
	obj, err := w.service.Send(workers.SpaceWorker, workers.CommandSet, w.encodeWorkerSet(&service.WorkerSetRequest{
		Id:         id,
		Profession: profession,
		Name:       name,
		Driver:     driver,
		Body:       data,
	}))
	if err != nil {
		return nil, err
	}
	worker, err := workers.ReadTWorker(obj)
	if err != nil {
		return nil, err
	}
	return worker, nil
}
func (w *WorkersRaft) encodeWorkerSet(body *service.WorkerSetRequest) []byte {
	data, _ := json.Marshal(body)
	return data
}
func (w *WorkersRaft) decodeWorkerSet(data []byte) (*service.WorkerSetRequest, error) {
	body := new(service.WorkerSetRequest)
	err := json.Unmarshal(data, body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
func (w *WorkersRaft) processHandlerWorkerSet(body []byte) (*eosc.WorkerData, error) {
	// reuest
	request, err := w.decodeWorkerSet(body)
	if err != nil {
		log.Warn("decode woker data:", err)
		return nil, fmt.Errorf("decode woker data:%w", err)
	}
	log.Info("process worker set: ", request.Id)
	ow, has := w.data.Get(request.Id)
	if has {
		// can not change driver
		if ow.Driver != request.Driver {
			return nil, workers.ErrorChangeDriver
		}
	} else {
		pf, b := w.professions.GetProfession(request.Profession)
		if !b {
			return nil, eosc.ErrorProfessionNotExist
		}
		if !pf.HasDriver(request.Driver) {
			return nil, eosc.ErrorDriverNotExist
		}
	}

	// check on process worker
	response, err := w.workerServiceClient.SetCheck(context.TODO(), request)
	if err != nil {
		return nil, err
	}

	if response.Status != service.WorkerStatusCode_SUCCESS {
		return nil, errors.New(response.Message)
	}

	createTime := eosc.Now()
	if has {
		createTime = ow.WorkerData.Create
	}
	return &eosc.WorkerData{
		Id:         request.Id,
		Name:       request.Name,
		Profession: request.Profession,
		Driver:     request.Driver,
		Create:     createTime,
		Update:     eosc.Now(),
		Body:       request.Body,
	}, nil
}
func (w *WorkersRaft) ProcessHandler(cmd string, body []byte) ([]byte, interface{}, error) {

	switch cmd {
	case workers.CommandSet:
		workerData, err := w.processHandlerWorkerSet(body)
		if err != nil {
			log.Info("process command set error: ", err)
			return nil, nil, err
		}
		data, err := workers.EncodeWorkerData(workerData)
		if err != nil {
			return nil, nil, err
		}
		return data, workerData, err
	case workers.CommandDel:
		request := &service.WorkerDeleteRequest{
			Id: string(body),
		}
		log.Info("process command delete: ", request.Id)
		ow, has := w.data.Get(request.Id)
		if !has {
			return nil, nil, workers.ErrorNotExist
		}
		response, err := w.workerServiceClient.DeleteCheck(context.TODO(), request)
		if err != nil {
			return nil, nil, err
		}
		if response.Status != service.WorkerStatusCode_SUCCESS {
			return nil, nil, errors.New(response.Message)
		}

		return body, ow.WorkerData, nil
	}
	return nil, nil, workers.ErrorInvalidCommand

}

func (w *WorkersRaft) CommitHandler(cmd string, data []byte) error {

	switch cmd {
	case workers.CommandSet:
		{
			worker, err := workers.DecodeWorker(data)
			if err != nil {
				return err
			}
			log.Info("commit worker set:", worker.Id)

			w.data.Set(worker.Id, worker)
			req := &service.WorkerSetRequest{
				Id:         worker.Id,
				Profession: worker.Profession,
				Name:       worker.Name,
				Driver:     worker.Driver,
				Body:       worker.Body,
			}
			response, err := w.workerServiceClient.Set(context.TODO(), req)
			if err != nil {
				log.Warn("set worker:", err)
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
			log.Info("commit worker delete:", id)
			w.data.Del(id)
			req := &service.WorkerDeleteRequest{
				Id: id,
			}
			response, err := w.workerServiceClient.Delete(context.TODO(), req)
			if err != nil {
				log.Warn("delete worker:", err)
				return err
			}
			if response.Status != service.WorkerStatusCode_SUCCESS {
				return errors.New(response.Message)
			}
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
	log.Debug("try restart...")
	w.workerProcessController.Restart()
	//if err != nil {
	//	log.Error("reset handler error: ", err)
	//}
	return nil
}

func (w *WorkersRaft) GetWork(id string) (eosc.TWorker, error) {
	if ow, b := w.data.Get(id); b {
		return ow.Detail(), nil
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
