package process_admin

import (
	"fmt"
	"github.com/eolinker/eosc/process-admin/marshal"
	"github.com/eolinker/eosc/process-admin/workers"
	"net/http"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	open_api "github.com/eolinker/eosc/open-api"
	"github.com/eolinker/eosc/variable"
	"github.com/julienschmidt/httprouter"
)

type VariableApi struct {
	extenderData *ExtenderData
	workers      workers.IWorkers
	variableData eosc.IVariable
	setting      eosc.ISettings
}

func NewVariableApi(extenderData *ExtenderData, workers workers.IWorkers, variableData eosc.IVariable, setting eosc.ISettings) *VariableApi {
	return &VariableApi{extenderData: extenderData, workers: workers, variableData: variableData, setting: setting}
}

func (oe *VariableApi) Register(router *httprouter.Router) {
	router.GET("/variable", open_api.CreateHandleFunc(oe.getAll))
	router.GET("/variable/:namespace", open_api.CreateHandleFunc(oe.getByNamespace))
	router.GET("/variable/:namespace/:key", open_api.CreateHandleFunc(oe.getByKey))
	router.POST("/variable/:namespace", open_api.CreateHandleFunc(oe.setByNamespace))
	router.PUT("/variable/:namespace", open_api.CreateHandleFunc(oe.setByNamespace))

}

func (oe *VariableApi) getAll(r *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {

	return http.StatusOK, nil, nil, oe.variableData
}

func (oe *VariableApi) getByNamespace(r *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {
	namespace := params.ByName("namespace")
	if namespace == "" {
		namespace = "default"
	}
	data, has := oe.variableData.GetByNamespace(namespace)
	if !has {
		return http.StatusNotFound, nil, nil, fmt.Sprintf("namespace{%s} not found", namespace)
	}
	return http.StatusOK, nil, nil, data
}

func (oe *VariableApi) getByKey(r *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {
	namespace := params.ByName("namespace")
	key := params.ByName("key")
	if namespace == "" {
		namespace = "default"
	}
	data, has := oe.variableData.GetByNamespace(namespace)
	if !has {
		return http.StatusNotFound, nil, nil, fmt.Sprintf("namespace{%s} not found", namespace)
	}
	value, ok := data[fmt.Sprintf("%s@%s", key, namespace)]
	if !ok {
		return http.StatusNotFound, nil, nil, fmt.Sprintf("key{%s} not found", key)
	}
	return http.StatusOK, nil, nil, value
}

func (oe *VariableApi) setByNamespace(r *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {
	namespace := params.ByName("namespace")
	if namespace == "" {
		namespace = "default"
	}
	decoder, err := marshal.GetData(r)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	cb := make(map[string]string)
	errUnmarshal := decoder.UnMarshal(&cb)
	if errUnmarshal != nil {
		return http.StatusInternalServerError, nil, nil, errUnmarshal
	}
	log.Debug("check variable...")
	affectIds, clone, err := oe.variableData.Check(namespace, cb)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, fmt.Sprintf("namespace{%s} not found", namespace)
	}
	log.Debug("parse variable...")
	parse := variable.NewParse(clone)
	workerToUpdate := make([]CacheItem, 0, len(affectIds))
	for _, id := range affectIds {
		profession, name, success := eosc.SplitWorkerId(id)
		if !success {
			continue
		}
		if profession != Setting {
			info, err := oe.workers.GetEmployee(profession, name)
			if err != nil {
				return http.StatusInternalServerError, nil, nil, fmt.Sprintf("worker(%s) not found, error is %s", id, err)
			}
			_, _, err = parse.Unmarshal(info.Body(), info.ConfigType())
			if err != nil {
				return http.StatusInternalServerError, nil, nil, fmt.Sprintf("unmarshal error:%s,body is '%s'", err, string(info.Body()))
			}
			workerToUpdate = append(workerToUpdate, CacheItem{
				id:         id,
				profession: profession,
			})
		} else {
			err := oe.setting.CheckVariable(name, clone)
			if err != nil {
				return http.StatusInternalServerError, nil, nil, fmt.Sprintf("setting %s unmarshal error:%s", name, err)
			}
			workerToUpdate = append(workerToUpdate, CacheItem{
				id:         name,
				profession: Setting,
			})
		}

	}
	log.Debug("update variable...")
	transaction := oe.workers.Begin(r.Context())
	for _, w := range workerToUpdate {
		if w.profession != Setting {
			err := transaction.Rebuild(w.id)
			if err != nil {
				transaction.Rollback()
				return http.StatusInternalServerError, nil, nil, err
			}
		} else {

			if err := oe.setting.Update(w.id, oe.variableData); err != nil {
				transaction.Rollback()
				return http.StatusInternalServerError, nil, nil, err
			}
		}
	}

	log.Debug("set variable...")

	if err := oe.variableData.SetByNamespace(namespace, cb); err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	log.Debug("set variable over...")
	data, _ := decoder.Encode()
	return http.StatusOK, nil, []*open_api.EventResponse{{
			Event:     "set",
			Namespace: "variable",
			Key:       namespace,
			Data:      data,
		},
		}, map[string]interface{}{
			"namespace": namespace,
			"variables": cb,
		}
}
