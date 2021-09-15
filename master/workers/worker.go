package workers

import (
	"encoding/json"
	"errors"
	"time"

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

func (w *Worker) ProcessHandler(propose []byte) (string, []byte, error) {
	return SpaceWorker, propose, nil
}

func (w *Worker) CommitHandler(data []byte) error {
	if w.store == nil {
		return errors.New("no valid store")
	}
	kv := &Cmd{}
	err := kv.Decode(data)
	if err != nil {
		return err
	}
	switch kv.Key {
	case CommandSet:
		{
			if kv.Config.CreateTime == "" {
				kv.Config.CreateTime = time.Now().Format("2006-01-02 15:04:05")
			}
			if kv.Config.UpdateTime == "" {
				kv.Config.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
			}
			b, err := json.Marshal(kv.Config)
			if err != nil {
				return err
			}
			storeValue := eosc.StoreValue{
				Id:         kv.Config.Id,
				Profession: kv.Config.Profession,
				Name:       kv.Config.Name,
				Driver:     kv.Config.Driver,
				CreateTime: kv.Config.CreateTime,
				UpdateTime: kv.Config.UpdateTime,
				IData:      eosc.JsonData(b),
				Sing:       "",
			}
			err = w.store.Set(storeValue)
			if err != nil {
				return err
			}
			return nil
		}
	case CommandDel:
		{
			err = w.store.Del(kv.Config.Id)
			if err != nil {
				return err
			}
			return nil
		}
	default:
		return raft_service.ErrInvalidKey
	}
}
