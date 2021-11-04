package process_master

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/eolinker/eosc/log"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/process-master/workers"
	raft_service "github.com/eolinker/eosc/raft/raft-service"
	"github.com/eolinker/eosc/service"
	"google.golang.org/protobuf/proto"
)

type WorkersRaft struct {
	data                *WorkerConfigs
	professions         eosc.IProfessions
	workerServiceClient service.WorkerServiceClient
	service             raft_service.IService
}

func (w *WorkersRaft) Append(cmd string, data []byte) error {

	switch cmd {
	case workers.CommandSet:
		{
			worker, err := workers.DecodeWorker(data)
			if err != nil {
				return err
			}
			log.Info("append worker set:", worker.Id)

			w.data.Set(worker.Id, worker)

			return nil
		}
	case workers.CommandDel:
		{

			id := string(data)
			log.Info("append worker delete:", id)
			w.data.Del(id)

			return nil
		}
	default:
		return raft_service.ErrInvalidCommand
	}
}

func (w *WorkersRaft) Complete() error {

	return nil
}

func NewWorkersRaft(WorkerConfig *WorkerConfigs, professions eosc.IProfessions, workerServiceClient service.WorkerServiceClient, service raft_service.IService) *WorkersRaft {
	// 初始化单例的worker
	for _, p := range professions.All() {
		if p.Mod == eosc.ProfessionConfig_Singleton {
			for _, d := range p.Drivers {
				id, _ := eosc.ToWorkerId(d.Name, p.Name)
				wr, _ := workers.ToWorker(&eosc.WorkerConfig{
					Id:         id,
					Profession: p.Name,
					Name:       d.Name,
					Driver:     d.Name,
					Create:     eosc.Now(),
					Update:     eosc.Now(),
					Body:       nil,
				})

				WorkerConfig.Set(id, wr)

			}
		}
	}
	return &WorkersRaft{data: WorkerConfig, professions: professions, workerServiceClient: workerServiceClient, service: service}
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
func (w *WorkersRaft) processHandlerWorkerSet(body []byte) (*eosc.WorkerConfig, error) {
	// reuest
	request, err := w.decodeWorkerSet(body)
	if err != nil {
		log.Warn("decode woker data:", err)
		return nil, fmt.Errorf("decode woker data:%w", err)
	}
	pf, b := w.professions.GetProfession(request.Profession)
	if !b {
		return nil, eosc.ErrorProfessionNotExist
	}
	log.Info("process worker set: ", request.Id)

	createTime := eosc.Now()

	if ow, has := w.data.Get(request.Id); has {
		// can not change driver

		if len(request.Driver) > 0 && ow.Driver != request.Driver {
			return nil, workers.ErrorChangeDriver
		}
		request.Driver = ow.Driver
		createTime = ow.WorkerConfig.Create
	} else {

		if pf.Mod() == eosc.ProfessionConfig_Singleton {
			return nil, eosc.ErrorNotAllowCreateForSingleton
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

	return &eosc.WorkerConfig{
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
		conf, err := w.processHandlerWorkerSet(body)
		if err != nil {
			log.Info("process command set error: ", err)
			return nil, nil, err
		}
		data, err := workers.EncodeWorkerData(conf)
		if err != nil {
			return nil, nil, err
		}
		return data, conf, err
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

		return body, ow.WorkerConfig, nil
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
	export := w.data.export()
	data, err := encode(export)
	if err != nil {
		return nil
	}
	return data
}

func (w *WorkersRaft) ResetHandler(data []byte) error {

	vs, err := decode(data)
	if err != nil {
		return err
	}

	w.data.reset(vs)
	//log.Debug("try restart...")
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
func decode(data []byte) ([]*eosc.WorkerConfig, error) {
	if len(data) == 0 {
		return nil, ErrClientNotInit
	}
	wd := new(eosc.WorkerConfigs)
	err := proto.Unmarshal(data, wd)
	if err != nil {
		return nil, err
	}
	return wd.Data, nil
}
func encode(cs []*eosc.WorkerConfig) ([]byte, error) {
	wd := new(eosc.WorkerConfigs)
	wd.Data = cs
	return proto.Marshal(wd)

}
