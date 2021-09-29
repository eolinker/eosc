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
	*eosc.WorkerData
	Data WorkerAttr
	Info *eosc.WorkerInfo
}

func (w *Worker) MarshalJSON() ([]byte, error) {

	if w.Data == nil {
		return nil, ErrorInvalidWorkerData
	}

	return EncodeWorkerData(w.WorkerData)
}
func EncodeWorkerData(wd *eosc.WorkerData) ([]byte, error) {
	return json.Marshal(wd)
}
func DecodeWorkerData(data []byte) (*eosc.WorkerData, error) {
	wd := new(eosc.WorkerData)
	err := json.Unmarshal(data, wd)
	if err != nil {
		return nil, err
	}
	return wd, nil
}
func DecodeWorker(data []byte) (*Worker, error) {

	wd, err := DecodeWorkerData(data)
	if err != nil {
		return nil, err
	}
	return ToWorker(wd)
}
func ReadTWorker(obj interface{}) (map[string]interface{}, error) {
	switch v := obj.(type) {
	case []byte:
		wk, err := DecodeWorker(v)
		if err != nil {
			return nil, err
		}
		return wk.Format(nil), nil
	case *eosc.WorkerData:
		wk, err := ToWorker(v)
		if err != nil {
			return nil, err
		}
		return wk.Format(nil), nil
	}
	return nil, errors.New("unknown type")
}
func ToWorker(wd *eosc.WorkerData) (*Worker, error) {
	wa := make(WorkerAttr)
	err := json.Unmarshal(wd.Body, &wa)
	if err != nil {
		return nil, err
	}
	return &Worker{
		WorkerData: wd,
		Data:       wa,
		Info: &eosc.WorkerInfo{
			Id:         wd.Id,
			Profession: wd.Profession,
			Name:       wd.Name,
			Driver:     wd.Driver,
			Create:     wd.Create,
			Update:     wd.Update,
		},
	}, nil
}

func (w *Worker) Format(attrs []string) map[string]interface{} {
	m := make(map[string]interface{})
	m["id"] = w.Id
	m["profession"] = w.Profession
	m["name"] = w.Name
	m["driver"] = w.Driver
	m["create"] = w.Create
	m["update"] = w.Update
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
