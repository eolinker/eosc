package process_worker

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolinker/eosc/log"
	"strings"

	"github.com/eolinker/eosc"
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
			w := new(eosc.WorkerConfig)
			err := json.Unmarshal(data, w)
			if err != nil {
				return err
			}

			return ws.workers.Set(w.Id, w.Profession, w.Name, w.Driver, w.Body, ws.variableManager.GetAll())
		}
	case eosc.NamespaceVariable:
		{
			var tmp map[string]string
			err := json.Unmarshal(data, &tmp)
			if err != nil {
				return err
			}
			_, _, err = ws.variableManager.SetByNamespace(key, tmp)
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

	for namespace, config := range eventData {
		for _, c := range config {
			switch namespace {
			case eosc.NamespaceProfession:
				{
					p := new(eosc.ProfessionConfig)
					err := json.Unmarshal(c, p)
					if err != nil {
						continue
					}
					pc = append(pc, p)
				}
			case eosc.NamespaceWorker:
				{
					w := new(eosc.WorkerConfig)
					err := json.Unmarshal(c, w)
					if err != nil {
						continue
					}
					wc = append(wc, w)
				}
			case eosc.NamespaceVariable:
				{
					var tmp map[string]string
					err := json.Unmarshal(c, &tmp)
					if err != nil {
						continue
					}
					target := make(map[string]map[string]string)
					for key, value := range tmp {
						name := "default"
						index := strings.Index(key, "@")
						if index > 0 && len(key) > index+1 {
							name = key[index+1:]
						}
						if _, ok := target[name]; !ok {
							target[name] = make(map[string]string)
						}
						target[name][key] = value
					}
					for key, value := range target {
						ws.variableManager.SetByNamespace(key, value)
					}
				}
			}
		}
	}

	ws.professionManager.Reset(pc)
	ws.onceInit.Do(func() {
		for _, h := range ws.initHandler {
			h()
		}
	})
	ws.workers.Reset(wc, ws.variableManager.GetAll())

	return nil
}
