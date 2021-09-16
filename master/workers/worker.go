package workers

import (
	"encoding/json"
	"errors"

	"github.com/eolinker/eosc"
	raft_service "github.com/eolinker/eosc/raft/raft-service"
	"github.com/eolinker/eosc/store"
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

type Worker struct {
	store eosc.IStore
}

func (w *Worker) Snapshot() []byte {
	values := w.store.All()
	data, _ := json.Marshal(values)
	return data
}

func (w *Worker) ResetHandler(data []byte) error {
	values := make([]eosc.StoreValue, 0, 10)
	err := json.Unmarshal(data, &values)
	if err != nil {
		return err
	}

	return w.store.Reset(values)
}

func NewWorker() *Worker {
	return &Worker{store: store.NewStore()}
}

func (w *Worker) ProcessHandler(cmd string, propose interface{}) ([]byte, error) {
	return json.Marshal(propose)
}

func (w *Worker) CommitHandler(cmd string, data []byte) error {
	if w.store == nil {
		return errors.New("no valid store")
	}

	switch cmd {
	case CommandSet:
		{
			iData := eosc.BytesData(data)
			value := new(eosc.StoreValue)
			err := iData.UnMarshal(value)
			if err != nil {
				return err
			}
			err = w.store.Set(*value)
			if err != nil {
				return err
			}
			return nil
		}
	case CommandDel:
		{
			id := string(data)

			err := w.store.Del(id)
			if err != nil {
				return err
			}
			return nil
		}
	default:
		return raft_service.ErrInvalidCommand
	}
}
