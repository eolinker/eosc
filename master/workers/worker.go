package workers

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrorInvalidWorkerData = errors.New("invalid worker data")
	ErrorNotExist          = errors.New("not exist")
	ErrorUnknown           = errors.New("unknown error")
	ErrorChangeDriver      = errors.New("try change driver")
	ErrorInvalidProfession = errors.New("invalid profession")
	ErrorInvalidDriver     = errors.New("invalid driver")
)

type WorkerData struct {
	Id         string `json:"id"`
	Profession string `json:"profession"`
	Name       string `json:"name"`
	Driver     string `json:"driver"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
	//Sing       string `json:"sing"`
	Data []byte `json:"data"`
}

type WorkerAttr map[string]interface{}

func (wa WorkerAttr) Get(field string) string {
	v, has := wa[field]
	if !has {
		return ""
	}
	return fmt.Sprint(v)
}

type Worker struct {
	Id         string
	Profession string
	Name       string
	Driver     string
	CreateTime string
	UpdateTime string
	Data       WorkerAttr
}

func (w *Worker) MarshalJSON() ([]byte, error) {
	if w.Data == nil {
		return nil, ErrorInvalidWorkerData
	}

	data, err := json.Marshal(w.Data)
	if err != nil {
		return nil, err
	}

	wd := &WorkerData{
		Id:         w.Id,
		Profession: w.Profession,
		Name:       w.Name,
		Driver:     w.Driver,
		CreateTime: w.UpdateTime,
		UpdateTime: w.UpdateTime,
		Data:       data,
	}
	return json.Marshal(wd)
}
func encodeWorker(w *Worker) ([]byte, error) {
	return w.MarshalJSON()
}
func decodeWorker(data []byte) (*Worker, error) {
	w := new(WorkerData)
	err := json.Unmarshal(data, w)
	if err != nil {
		return nil, err
	}
	wa := make(WorkerAttr)
	err = json.Unmarshal(w.Data, &wa)
	if err != nil {
		return nil, err
	}
	return &Worker{
		Id:         w.Id,
		Profession: w.Profession,
		Name:       w.Name,
		Driver:     w.Driver,
		CreateTime: w.CreateTime,
		UpdateTime: w.UpdateTime,
		Data:       wa,
	}, nil
}

func (w *Worker) Format(attrs []string) map[string]interface{} {
	m := make(map[string]interface{})
	m["id"] = w.Id
	m["profession"] = w.Profession
	m["name"] = w.Name
	m["create"] = w.CreateTime
	m["update"] = w.UpdateTime
	if w.Data != nil {
		for _, n := range attrs {
			m[n] = w.Data[n]
		}
	}
	return m
}
