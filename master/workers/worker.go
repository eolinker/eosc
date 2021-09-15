package workers

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"github.com/eolinker/eosc/store"
	"time"

	raft_service "github.com/eolinker/eosc/raft/raft-service"

	"github.com/eolinker/eosc"
)

var (
	commandSet = "set"
	commandDel = "delete"
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
	return &Worker{store: 	store.NewStore()}
}

func (w *Worker) ProcessHandler(propose []byte) (string, []byte, error) {
	return eosc.SpaceWorker, propose, nil
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
	case commandSet:
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
	case commandDel:
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

// Cmd 用于传输的结构
type Cmd struct {
	Key    string
	Config *baseConfig
}

func (kv *Cmd) Encode() ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(kv); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (kv *Cmd) Decode(data []byte) error {
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	if err := dec.Decode(kv); err != nil {
		return err
	}
	return nil
}
