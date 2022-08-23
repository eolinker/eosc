package process_admin

import (
	"encoding/json"
	"fmt"
	"github.com/eolinker/eosc"
	open_api "github.com/eolinker/eosc/open-api"
	"github.com/eolinker/eosc/variable"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type VariableApi struct {
	extenderData *ExtenderData
	workers      *Workers
	variableData variable.IVariable
}

func NewVariableApi(extenderData *ExtenderData, workers *Workers, variableData variable.IVariable) *VariableApi {
	return &VariableApi{extenderData: extenderData, workers: workers, variableData: variableData}
}

func (oe *VariableApi) Register(router *httprouter.Router) {
	router.GET("/variable", open_api.CreateHandleFunc(oe.getAll))
	router.GET("/variable/:namespace", open_api.CreateHandleFunc(oe.getByNamespace))
	router.GET("/variable/:namespace/:key", open_api.CreateHandleFunc(oe.getByKey))
	router.POST("/variable/:namespace", open_api.CreateHandleFunc(oe.setByNamespace))
	router.PUT("/variable/:namespace", open_api.CreateHandleFunc(oe.setByNamespace))
	
}

func (oe *VariableApi) getAll(r *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {
	namespaces := oe.variableData.Namespaces()
	all := make(map[string]interface{})
	for _, namespace := range namespaces {
		data, has := oe.variableData.GetByNamespace(namespace)
		if !has {
			continue
		}
		all[namespace] = variable.TrimNamespace(data)
	}
	return http.StatusOK, nil, nil, all
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
	return http.StatusOK, nil, nil, variable.TrimNamespace(data)
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
	decoder, err := GetData(r)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	cb := make(map[string]string)
	errUnmarshal := decoder.UnMarshal(&cb)
	if errUnmarshal != nil {
		return http.StatusInternalServerError, nil, nil, errUnmarshal
	}
	
	variables := variable.FillNamespace(namespace, cb)
	
	affectIds, err := oe.variableData.Check(namespace, variables)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, fmt.Sprintf("namespace{%s} not found", namespace)
	}
	
	parse := variable.NewParse(variables)
	es := make([]*open_api.EventResponse, 0, len(affectIds)+1)
	data, _ := json.Marshal(variables)
	es = append(es, &open_api.EventResponse{
		Event:     "set",
		Namespace: "variable",
		Key:       namespace,
		Data:      data,
	})
	
	for _, id := range affectIds {
		profession, name, success := eosc.SplitWorkerId(id)
		if !success {
			continue
		}
		info, err := oe.workers.GetEmployee(profession, name)
		if err != nil {
			return http.StatusInternalServerError, nil, nil, fmt.Sprintf("worker(%s) not found, error is %s", id, err)
		}
		_, _, err = parse.Unmarshal(info.Body(), info.configType)
		if err != nil {
			return http.StatusInternalServerError, nil, nil, fmt.Sprintf("unmarshal error:%s,body is '%s'", err, string(info.Body()))
		}
		eventData, _ := json.Marshal(info.config)
		es = append(es, &open_api.EventResponse{
			Event:     "set",
			Namespace: eosc.NamespaceWorker,
			Key:       info.config.Id,
			Data:      eventData,
		})
	}
	oe.variableData.SetByNamespace(namespace, variables)
	return http.StatusOK, nil, es, map[string]interface{}{
		"namespace": namespace,
		"variables": cb,
	}
}
