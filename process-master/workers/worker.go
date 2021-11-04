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
	*eosc.WorkerConfig
	Data WorkerAttr
}

func (w *Worker) MarshalJSON() ([]byte, error) {

	if w.Data == nil {
		return nil, ErrorInvalidWorkerData
	}

	return EncodeWorkerData(w.WorkerConfig)
}
func EncodeWorkerData(wd *eosc.WorkerConfig) ([]byte, error) {
	return json.Marshal(wd)
}
func DecodeWorkerData(data []byte) (*eosc.WorkerConfig, error) {
	wd := new(eosc.WorkerConfig)
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
		return wk.Info(), nil
	case *eosc.WorkerConfig:
		wk, err := ToWorker(v)
		if err != nil {
			return nil, err
		}
		return wk.Info(), nil
	}
	return nil, errors.New("unknown type")
}
func ToWorker(wd *eosc.WorkerConfig) (*Worker, error) {
	wa := make(WorkerAttr)
	if len(wd.Body) > 0 {
		err := json.Unmarshal(wd.Body, &wa)
		if err != nil {
			return nil, err
		}
	}

	return &Worker{
		WorkerConfig: wd,
		Data:         wa,
	}, nil
}
func (w *Worker) Info() eosc.TWorker {
	m := make(map[string]interface{})
	m["id"] = w.Id
	m["profession"] = w.Profession
	m["name"] = w.Name
	m["driver"] = w.Driver
	m["create"] = w.Create
	m["update"] = w.Update
	return m
}
func (w *Worker) Detail() eosc.TWorker {
	m := w.Info()
	if w.Data != nil {

		for k, v := range w.Data {
			m[k] = v
		}

	}
	return m
}

func (w *Worker) Format(attrs []string) eosc.TWorker {
	m := w.Info()
	if w.Data != nil {
		if attrs != nil {
			for _, f := range attrs {
				m[f] = w.Data[f]
			}
		}
	}

	return m
}
