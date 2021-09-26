package workers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/eolinker/eosc"
)

var (
	ErrorInvalidWorkerData = errors.New("invalid process-worker data")
	ErrorNotExist          = errors.New("not exist")
	ErrorUnknown           = errors.New("unknown error")
	ErrorChangeDriver      = errors.New("try change driver")
	ErrorInvalidProfession = errors.New("invalid profession")
	ErrorInvalidDriver     = errors.New("invalid driver")
	ErrorInvalidCommand    = errors.New("invalid command")
)

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
	Org        *eosc.WorkerData
	Info       *eosc.WorkerInfo
}

func (w *Worker) MarshalJSON() ([]byte, error) {

	if w.Data == nil {
		return nil, ErrorInvalidWorkerData
	}

	//data, err := json.Marshal(w.Data)
	//if err != nil {
	//	return nil, err
	//}

	//wd := &eosc.WorkerData{
	//	Id:         w.Id,
	//	Profession: w.Profession,
	//	Name:       w.Name,
	//	Driver:     w.Driver,
	//	Create:     w.UpdateTime,
	//	Update:     w.UpdateTime,
	//	Body:       data,
	//
	//}
	return json.Marshal(w.Org)
}

func DecodeWorker(data []byte) (*Worker, error) {
	w := new(eosc.WorkerData)
	err := json.Unmarshal(data, w)
	if err != nil {
		return nil, err
	}
	wa := make(WorkerAttr)
	err = json.Unmarshal(w.Body, &wa)
	if err != nil {
		return nil, err
	}
	return &Worker{
		Id:         w.Id,
		Profession: w.Profession,
		Name:       w.Name,
		Driver:     w.Driver,
		CreateTime: w.Create,
		UpdateTime: w.Update,
		Data:       wa,
		Org:        w,
		Info: &eosc.WorkerInfo{
			Id:         w.Id,
			Profession: w.Profession,
			Name:       w.Name,
			Driver:     w.Driver,
			Create:     w.Create,
			Update:     w.Update,
		},
	}, nil
}
func ToWorker(wd *eosc.WorkerData) (*Worker, error) {
	wa := make(WorkerAttr)
	err := json.Unmarshal(wd.Body, &wa)
	if err != nil {
		return nil, err
	}
	return &Worker{
		Id:         wd.Id,
		Profession: wd.Profession,
		Name:       wd.Name,
		Driver:     wd.Driver,
		CreateTime: wd.Create,
		UpdateTime: wd.Update,
		Data:       wa,
		Org:        wd,
	}, nil
}
func (w *Worker) Format(attrs []string) map[string]interface{} {
	m := make(map[string]interface{})
	m["id"] = w.Id
	m["profession"] = w.Profession
	m["name"] = w.Name
	m["driver"] = w.Driver
	m["create"] = w.CreateTime
	m["update"] = w.UpdateTime
	if w.Data != nil {
		if attrs != nil {
			for _, f := range attrs {
				m[f] = w.Data[f]
			}
		} else {
			for k, v := range w.Data {
				m[k] = v
			}
		}
	}
	return m
}
