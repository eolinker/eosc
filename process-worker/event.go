package process_worker

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
)

func (ws *WorkerServer) setEvent(namespace string, key string, data []byte) (errResult error) {

	switch namespace {

	case eosc.NamespaceProfession:
		{
			p := new(eosc.ProfessionConfig)
			err := json.Unmarshal(data, p)
			if err != nil {
				log.Error("unmarshal profession data error:", err)
				return err
			}

			return ws.professionManager.Set(key, p)
		}
	case eosc.NamespaceWorker:
		{
			log.Debug("setEvent NamespaceWorker:")

			w := new(eosc.WorkerConfig)
			err := json.Unmarshal(data, w)
			if err != nil {
				return err
			}
			log.Debug("NamespaceWorker:", w.Profession, w.Name)

			if w.Profession == "setting" {
				err = ws.settings.SettingWorker(w.Name, w.Body)
				return err
			}

			return ws.workers.Set(w.Id, w.Profession, w.Name, w.Driver, w.Body)
		}
	case eosc.NamespaceVariable:
		{
			var tmp map[string]string
			err := json.Unmarshal(data, &tmp)
			if err != nil {
				return err
			}

			wids, oldValue, err := ws.variableManager.SetByNamespace(key, tmp)
			if err != nil {
				return err
			}
			defer func() {
				if errResult != nil {
					ws.variableManager.SetByNamespace(key, oldValue)
				}
			}()
			for _, id := range wids {
				profession, _, success := eosc.SplitWorkerId(id)
				if !success {
					continue
				}
				if profession == "setting" {
					ws.settings.Update(id)
				} else {
					ws.workers.Update(id)
				}
			}
			return err
		}
	case eosc.NamespaceCustomer:
		var value map[string]string
		err := json.Unmarshal(data, &value)
		if err != nil {
			return err
		}
		ws.customerVar.Set(key, value)
		return nil
	default:
		return nil
		//return errors.New(fmt.Sprintf("namespace %s is not existed.", namespace))
	}

}

func (ws *WorkerServer) delEvent(namespace string, key string) error {
	switch namespace {
	case eosc.NamespaceProfession:
		{
			return ws.professionManager.Delete(key)
		}
	case eosc.NamespaceWorker:
		{
			return ws.workers.Del(key)
		}
	case eosc.NamespaceCustomer:
		{
			ws.customerVar.Set(key, nil)
			return nil
		}

	case eosc.NamespaceVariable:
		{
			ws.variableManager.SetByNamespace(key, make(map[string]string))
			return nil
		}
	default:
		return errors.New(fmt.Sprintf("namespace %s is not existed.", namespace))
	}
}

func (ws *WorkerServer) resetEvent(data []byte) error {
	eventData := make(map[string]map[string][]byte)
	if len(data) > 0 {
		err := json.Unmarshal(data, &eventData)
		if err != nil {
			return err
		}
	}

	pc := make([]*eosc.ProfessionConfig, 0)
	wc := make([]*eosc.WorkerConfig, 0)
	settings := make([]*eosc.WorkerConfig, 0)
	// 第一步 处理 profession
	if config, has := eventData[eosc.NamespaceProfession]; has {
		for _, c := range config {
			p := new(eosc.ProfessionConfig)
			err := json.Unmarshal(c, p)
			if err != nil {
				continue
			}
			pc = append(pc, p)
		}
		ws.professionManager.Reset(pc)
	}
	// profession的初始化
	ws.onceInit.Do(func() {
		for _, h := range ws.initHandler {
			h()
		}
	})
	if config, has := eventData[eosc.NamespaceCustomer]; has {
		for key, c := range config {
			var value map[string]string
			err := json.Unmarshal(c, &value)
			if err != nil {
				return err
			}
			ws.customerVar.Set(key, value)
		}
	}
	// 处理环境变量
	if config, has := eventData[eosc.NamespaceVariable]; has {
		for key, c := range config {
			value := make(map[string]string)
			err := json.Unmarshal(c, &value)
			if err != nil {
				continue
			}
			ws.variableManager.SetByNamespace(key, value)
		}
	}
	// 处理worker
	if config, has := eventData[eosc.NamespaceWorker]; has {
		for _, c := range config {
			w := new(eosc.WorkerConfig)
			err := json.Unmarshal(c, w)
			if err != nil {
				continue
			}
			log.Debug("init read worker:", w.Profession, ":", w.Name)
			if w.Profession == "setting" {
				settings = append(settings, w)
			} else {
				wc = append(wc, w)
			}
		}
	}
	// 处理setting
	for _, w := range settings {

		err := ws.settings.SettingWorker(w.Name, w.Body)
		if err != nil {
			log.Warn("set setting :", err)
		}
	}
	// 处理其他worker
	ws.workers.Reset(wc)

	return nil
}
