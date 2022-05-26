package process_admin

import (
	"encoding/json"
	"github.com/eolinker/eosc"
)

type WorkerInfo struct {
	worker       eosc.IWorker
	config       *eosc.WorkerConfig
	appendLabels []string
	attr         map[string]interface{}
	info         map[string]interface{}
}

func NewWorkerInfo(worker eosc.IWorker, id string, profession string, name, driver, create, update string, config interface{}, appendLabels []string) *WorkerInfo {

	body, _ := json.Marshal(config)

	return &WorkerInfo{
		appendLabels: appendLabels,
		worker:       worker,
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

func (w *WorkerInfo) reset(driver string, config interface{}, worker eosc.IWorker) {
	w.config.Update = eosc.Now()
	w.config.Driver = driver
	w.config.Body, _ = json.Marshal(config)
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
		m["name"] = w.config.Profession
		m["driver"] = w.config.Driver
		m["update"] = w.config.Update
		m["create"] = w.config.Create

		w.attr = m
	}

	return w.attr
}
func (w *WorkerInfo) Info() interface{} {
	if w.info == nil {
		detail := w.toDetails()
		w.info = make(map[string]interface{})
		for _, label := range w.appendLabels {
			w.info[label] = detail[label]
		}
		w.info["id"] = w.config.Id
		w.info["profession"] = w.config.Profession
		w.info["name"] = w.config.Profession
		w.info["driver"] = w.config.Driver
		w.info["update"] = w.config.Update
		w.info["create"] = w.config.Create
	}

	return w.info
}
