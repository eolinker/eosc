package process_worker

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/variable"
)

func (ws *WorkerServer) setEvent(namespace string, key string, data []byte) error {

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
				_, err = ws.settings.Set(w.Name, w.Body, ws.variableManager)
				return err
			}

			return ws.workers.Set(w.Id, w.Profession, w.Name, w.Driver, w.Body, ws.variableManager)
		}
	case eosc.NamespaceVariable:
		{
			var tmp map[string]string
			err := json.Unmarshal(data, &tmp)
			if err != nil {
				return err
			}

			wids, clone, err := ws.variableManager.Check(key, tmp)
			if err != nil {
				return err

			}
			ws.variableManager.SetByNamespace(key, tmp)
			for _, id := range wids {
				profession, name, success := eosc.SplitWorkerId(id)
				if !success {
					continue
				}
				if profession == "setting" {
					ws.settings.Update(name, clone)
				} else {
					ws.workers.Update(id, clone)
				}
			}
			return err
		}
	default:
		return errors.New(fmt.Sprintf("namespace %s is not existed.", namespace))
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
	case eosc.NamespaceVariable:
		{
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
	for namespace, config := range eventData {

		switch namespace {
		case eosc.NamespaceProfession:
			{
				for _, c := range config {
					p := new(eosc.ProfessionConfig)
					err := json.Unmarshal(c, p)
					if err != nil {
						continue
					}
					pc = append(pc, p)
				}
			}
		case eosc.NamespaceWorker:
			{
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
		case eosc.NamespaceVariable:
			{
				ws.variableManager = variable.NewVariables(config)
			}
		}
	}

	ws.professionManager.Reset(pc)
	ws.onceInit.Do(func() {
		for _, h := range ws.initHandler {
			h()
		}
	})
	for _, w := range settings {

		_, err := ws.settings.Set(w.Name, w.Body, ws.variableManager)
		log.Warn("set setting :", err)
	}
	ws.workers.Reset(wc, ws.variableManager)

	return nil
}
