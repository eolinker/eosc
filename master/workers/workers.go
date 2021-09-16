package workers

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/eolinker/eosc/service"

	"github.com/eolinker/eosc/admin"

	"github.com/eolinker/eosc"
	raft_service "github.com/eolinker/eosc/raft/raft-service"
)

const (
	SpaceWorker = "worker"
)

var (
	CommandSet = "set"
	CommandDel = "delete"
)

type baseConfig struct {
	Id         string `json:"id" yaml:"id"`
	Name       string `json:"name" yaml:"name"`
	Profession string `json:"profession" yaml:"profession"`
	Driver     string `json:"driver" yaml:"driver"`
	CreateTime string `json:"create_time" yaml:"create_time"`
	UpdateTime string `json:"update_time" yaml:"update_time"`
}

type Workers struct {
	professions         admin.IProfessions
	data                eosc.IUntyped
	workerServiceClient service.WorkerServiceClient
}

func NewWorkers(professions admin.IProfessions) *Workers {
	return &Workers{professions: professions, data: eosc.NewUntyped()}
}

func (w *Workers) Snapshot() []byte {
	values := w.data.All()
	data, _ := json.Marshal(values)
	return data
}

func (w *Workers) ResetHandler(data []byte) error {
	values := make([]*Worker, 0, 10)
	err := json.Unmarshal(data, &values)
	if err != nil {
		return err
	}
	buf := eosc.NewUntyped()
	for _, v := range values {
		buf.Set(v.Id, v)
	}
	w.data = buf
	return nil
}

func (w *Workers) ProcessHandler(cmd string, body []byte) ([]byte, error) {
	err := w.checkClient()
	if err != nil {
		return nil, err
	}
	request := &service.WorkerCheckRequest{
		Cmd:  cmd,
		Body: body,
	}

	response, err := w.workerServiceClient.Check(context.TODO(), request)
	if err != nil {
		return nil, err
	}
	if response.Status != 0 {
		return nil, errors.New(response.Msg)
	}
	return response.Body, nil
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
