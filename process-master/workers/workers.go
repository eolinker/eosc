package workers

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/eolinker/eosc/service"

	"github.com/eolinker/eosc"
	raft_service "github.com/eolinker/eosc/raft/raft-service"
)

const (
	SpaceWorker = "process-worker"
)

var (
	CommandSet = "set"
	CommandDel = "delete"
)

type Workers struct {
	professions         eosc.IProfessionsData
	data                ITypedWorkers
	workerServiceClient service.WorkerServiceClient
	service             raft_service.IService
}

func NewWorkers(professions eosc.IProfessionsData, service raft_service.IService) *Workers {
	return &Workers{professions: professions, data: NewTypedWorkers(), service: service}
}

func (w *Workers) Snapshot() []byte {

	data, err := w.encode()
	if err != nil {
		return nil
	}
	return data
}
func (w *Workers) encode() ([]byte, error) {
	values := w.data.All()

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
func (w *Workers) decode(data []byte) ([]*eosc.WorkerData, error) {
	wd := new(eosc.WorkersData)
	err := proto.Unmarshal(data, wd)
	if err != nil {
		return nil, err
	}
	return wd.Data, nil
}
func (w *Workers) reset(vs []*eosc.WorkerData) {
	nw := NewTypedWorkers()
	for _, v := range vs {
		wv, err := toWorker(v)
		if err != nil {
			continue
		}
		nw.Set(v.Id, wv)
	}
	w.data = nw
}
func (w *Workers) ResetHandler(data []byte) error {
	values := make([]*Worker, 0, 10)
	err := json.Unmarshal(data, &values)
	if err != nil {
		return err
	}

	return nil
}

func (w *Workers) ProcessHandler(cmd string, body []byte) ([]byte, error) {
	err := w.checkClient()
	if err != nil {
		return nil, err
	}

	switch cmd {
	case CommandSet:
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
	case CommandDel:
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
	return nil, ErrorInvalidCommand

}

func (w *Workers) CommitHandler(cmd string, data []byte) error {
	switch cmd {
	case CommandSet:
		{
			worker, err := decodeWorker(data)
			if err != nil {
				return err
			}
			w.data.Set(worker.Id, worker)
			return nil
		}
	case CommandDel:
		{
			id := string(data)
			w.data.Del(id)
			return nil
		}
	default:
		return raft_service.ErrInvalidCommand
	}
}
func (w *Workers) checkClient() (errOut error) {

	for i := 0; i < 2; i++ {
		if w.workerServiceClient == nil {
			client, err := createClient()
			if err != nil {
				return
			}
			w.workerServiceClient = client
		}
		hello := strconv.FormatInt(time.Now().UnixNano(), 10)
		response, err := w.workerServiceClient.Ping(context.TODO(), &service.WorkerHelloRequest{
			Hello: hello,
		})
		if err != nil {
			w.workerServiceClient = nil
			errOut = err
			continue
		}
		if response.Hello != hello {
			w.workerServiceClient = nil
			continue
		}
		errOut = nil
		return
	}
	w.workerServiceClient = nil
	return
}
