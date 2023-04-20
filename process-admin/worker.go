package process_admin

import (
	"encoding/json"
	"reflect"

	"github.com/eolinker/eosc"
)

type WorkerInfo struct {
	worker       eosc.IWorker
	config       *eosc.WorkerConfig
	appendLabels []string
	attr         map[string]interface{}
	info         map[string]interface{}
	configType   reflect.Type
}

func NewWorkerInfo(worker eosc.IWorker, id string, profession string, name, driver, version, desc, create, update string, body []byte, configType reflect.Type) *WorkerInfo {

	return &WorkerInfo{

		worker: worker,
		config: &eosc.WorkerConfig{
			Id:          id,
			Profession:  profession,
			Name:        name,
			Driver:      driver,
			Create:      create,
			Update:      update,
			Description: desc,
			Body:        body,
			Version:     version,
		},
		configType: configType,
		attr:       nil,
	}
}

func (w *WorkerInfo) reset(driver, version, desc string, body []byte, worker eosc.IWorker, configType reflect.Type) {
	//w.config.Update = eosc.Now()
	if version != w.config.Version {
		w.config.Update = eosc.Now()
	}
	w.config.Driver = driver
	w.config.Description = desc
	w.config.Body = body
	w.config.Version = version
	w.configType = configType
	w.worker = worker
	w.info = nil
	w.attr = nil
}

func (w *WorkerInfo) Detail() interface{} {
	return w.toDetails()
}
func (w *WorkerInfo) toDetails() map[string]interface{} {
	if w.attr == nil {
		m := make(map[string]interface{})
		json.Unmarshal(w.config.Body, &m)
		m["id"] = w.config.Id
		m["profession"] = w.config.Profession
		m["name"] = w.config.Name
		m["driver"] = w.config.Driver
		m["version"] = w.config.Version
		m["description"] = w.config.Description
		m["update"] = w.config.Update
		m["create"] = w.config.Create
		w.attr = m
	}

	return w.attr
}
func (w *WorkerInfo) Info(appendLabels ...string) interface{} {
	if w.info == nil {
		detail := w.toDetails()
		w.info = make(map[string]interface{})
		for _, label := range appendLabels {
			w.info[label] = detail[label]
		}
		w.info["id"] = w.config.Id
		w.info["profession"] = w.config.Profession
		w.info["name"] = w.config.Name
		w.info["driver"] = w.config.Driver
		w.info["version"] = w.config.Version
		w.info["description"] = w.config.Description
		w.info["update"] = w.config.Update
		w.info["create"] = w.config.Create
	}

	return w.info
}

func (w *WorkerInfo) Body() []byte {
	return w.config.Body
}

func (w *WorkerInfo) ConfigType() reflect.Type {
	return w.configType
}
