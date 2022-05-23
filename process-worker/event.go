package process_worker

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"

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

			return ws.workers.Set(w.Id, w.Profession, w.Name, w.Driver, w.Body)
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

	professionTyp := reflect.TypeOf(&eosc.ProfessionConfig{})
	workerTyp := reflect.TypeOf(&eosc.WorkerConfig{})
	for namespace, config := range eventData {
		for _, c := range config {
			switch namespace {
			case eosc.NamespaceProfession:
				{

					p, err := toConfig(c, professionTyp)
					if err != nil {
						continue
					}
					pc = append(pc, p.(*eosc.ProfessionConfig))
				}
			case eosc.NamespaceWorker:
				{
					w, err := toConfig(c, workerTyp)
					if err != nil {
						continue
					}
					wc = append(wc, w.(*eosc.WorkerConfig))
				}
			}
		}
	}

	ws.professionManager.Reset(pc)
	ws.workers.Reset(wc)
	bean.Injection(&ws.workers)
	return nil
}

func toConfig(data []byte, typ reflect.Type) (interface{}, error) {
	cfg := newConfig(typ)
	err := json.Unmarshal(data, cfg)
	return cfg, err
}

func newConfig(t reflect.Type) interface{} {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return reflect.New(t).Interface()
}
