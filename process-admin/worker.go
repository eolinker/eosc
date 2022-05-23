package process_admin

import (
	"encoding/json"
	"github.com/eolinker/eosc"
)

type WorkerInfo struct {
	worker eosc.IWorker
	config *eosc.WorkerConfig
	attr   interface{}
}

func NewWorkerInfo(worker eosc.IWorker, id, profession, name, driver, create, update string, config interface{}) *WorkerInfo {

	body, _ := json.Marshal(config)

	return &WorkerInfo{
		worker: worker,
		config: &eosc.WorkerConfig{
			Id:         id,
			Profession: profession,
			Name:       name,
			Driver:     driver,
			Create:     create,
			Update:     update,
			Body:       body,
		},
		attr: nil,
	}
}

func (w *WorkerInfo) reset(driver string, config interface{}) {
	w.config.Update = eosc.Now()
	w.config.Driver = driver
	w.config.Body, _ = json.Marshal(config)
	w.attr = nil
}

func (w *WorkerInfo) toAttr() interface{} {
	if w.attr != nil {
		return w.attr
	}
	m := make(map[string]interface{})
	json.Unmarshal(w.config.Body, &m)
	m["id"] = w.config.Id
	m["profession"] = w.config.Profession
	m["name"] = w.config.Profession
	m["driver"] = w.config.Driver
	m["update"] = w.config.Update
	m["create"] = w.config.Create

	w.attr = m
	return m
}
