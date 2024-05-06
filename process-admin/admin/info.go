/*
 * Copyright (c) 2024. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package admin

import (
	"encoding/json"
	"reflect"

	"github.com/eolinker/eosc"
)

type WorkerInfo struct {
	worker     eosc.IWorker
	config     *eosc.WorkerConfig
	attr       map[string]interface{}
	info       map[string]interface{}
	configType reflect.Type
}

func GetProfession(w *WorkerInfo) string {
	return w.config.Profession

}
func NewWorkerInfo(worker eosc.IWorker, cf *eosc.WorkerConfig, configType reflect.Type) *WorkerInfo {
	if cf.Create == "" {
		cf.Create = eosc.Now()
	}
	if cf.Update == "" {
		cf.Update = eosc.Now()
	}
	return &WorkerInfo{
		worker:     worker,
		config:     cf,
		configType: configType,
		attr:       nil,
	}
}

func (w *WorkerInfo) reset(cf *eosc.WorkerConfig, worker eosc.IWorker, configType reflect.Type) {

	cf.Create = w.config.Create
	if cf.Create == "" {
		cf.Create = eosc.Now()
	}

	if cf.Update == "" {
		cf.Update = eosc.Now()
	}

	w.config = cf

	w.configType = configType
	w.worker = worker
	w.info = nil
	w.attr = nil
}

func (w *WorkerInfo) Detail() interface{} {
	return w.toDetails()
}
func (w *WorkerInfo) ConfigData() []byte {
	configData, _ := json.Marshal(w.config)
	return configData
}
func (w *WorkerInfo) toDetails() map[string]interface{} {
	if w.attr == nil {
		m := make(map[string]interface{})
		_ = json.Unmarshal(w.config.Body, &m)
		m["id"] = w.config.Id
		m["profession"] = w.config.Profession
		m["name"] = w.config.Name
		m["driver"] = w.config.Driver
		m["version"] = w.config.Version
		m["description"] = w.config.Description
		m["update"] = w.config.Update
		m["create"] = w.config.Create
		m["matches"] = w.config.Matches
		w.attr = m
	}

	return w.attr
}
func (w *WorkerInfo) Info(appendLabels ...string) interface{} {
	if w.info == nil {
		detail := w.toDetails()
		info := make(map[string]interface{})
		for _, label := range appendLabels {
			info[label] = detail[label]
		}
		info["id"] = w.config.Id
		info["profession"] = w.config.Profession
		info["name"] = w.config.Name
		info["driver"] = w.config.Driver
		info["version"] = w.config.Version
		info["description"] = w.config.Description
		info["update"] = w.config.Update
		info["create"] = w.config.Create
		w.info = info
	}

	return w.info
}
func (w *WorkerInfo) Matches() map[string]string {
	return w.config.Matches
}
func (w *WorkerInfo) Body() []byte {
	return w.config.Body
}
func (w *WorkerInfo) Id() string {
	return w.config.Id
}

func (w *WorkerInfo) ConfigType() reflect.Type {
	return w.configType
}
func (w *WorkerInfo) Description() string {
	return w.config.Description
}
func (w *WorkerInfo) Driver() string {
	return w.config.Driver
}
func (w *WorkerInfo) Version() string {
	return w.config.Version
}
func (w *WorkerInfo) GetWorker() eosc.IWorker {
	return w.worker
}
