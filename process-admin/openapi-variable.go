package process_admin

import (
	"encoding/json"
	"fmt"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	open_api "github.com/eolinker/eosc/open-api"
	"github.com/eolinker/eosc/variable"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

type VariableApi struct {
	extenderData *ExtenderData
	workerData   *WorkerDatas
	variableData variable.IVariable
}

func (oe *VariableApi) Register(router *httprouter.Router) {
	router.GET("/variable", open_api.CreateHandleFunc(oe.getAll))
	router.GET("/variable/:namespace", open_api.CreateHandleFunc(oe.getByNamespace))
	router.GET("/variable/:namespace/:key", open_api.CreateHandleFunc(oe.getByKey))
	router.POST("/variable/:namespace", open_api.CreateHandleFunc(oe.setByNamespace))
	router.PUT("/variable/:namespace", open_api.CreateHandleFunc(oe.setByNamespace))

}

func NewVariableApi() *VariableApi {
	return &VariableApi{variableData: variable.NewManager()}
}

func (oe *VariableApi) getAll(r *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {
	namespaces := oe.variableData.Namespaces()
	all := make(map[string]interface{})
	for _, namespace := range namespaces {
		data, has := oe.variableData.GetByNamespace(namespace)
		if !has {
			continue
		}
		all[namespace] = trimNamespace(data)
	}
	return http.StatusOK, nil, nil, all
}

func trimNamespace(origin map[string]string) map[string]string {
	target := make(map[string]string)
	for key, value := range origin {
		index := strings.Index(key, "@")
		if index > 0 {
			key = key[:index]
		}
		target[key] = value
	}
	return target
}

func fillNamespace(namespace string, origin map[string]string) map[string]string {
	target := make(map[string]string)
	for key, value := range origin {
		target[fmt.Sprintf("%s@%s", key, namespace)] = value
	}
	return target
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
	return http.StatusOK, nil, nil, trimNamespace(data)
}

func (oe *VariableApi) getByKey(r *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {
	namespace := params.ByName("namespace")
	key := params.ByName("key")
	if namespace == "" {
		namespace = "default"
	}
	data, has := oe.variableData.GetByNamespace(namespace)
	if !has {
		return 404, nil, nil, fmt.Sprintf("namespace{%s} not found", namespace)
	}
	value, ok := data[fmt.Sprintf("%s@%s", namespace, key)]
	if !ok {
		return http.StatusNotFound, nil, nil, fmt.Sprintf("namespace{%s} not found", namespace)
	}
	return http.StatusOK, nil, nil, map[string]string{
		key: value,
	}
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
	errUnmarshal := decoder.UnMarshal(cb)
	if errUnmarshal != nil {
		return http.StatusInternalServerError, nil, nil, errUnmarshal
	}
	variables := fillNamespace(namespace, cb)
	_, affectIds, err := oe.variableData.SetByNamespace(namespace, variables)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, fmt.Sprintf("namespace{%s} not found", namespace)
	}

	es := make([]*open_api.EventResponse, 0, len(affectIds)+1)
	data, _ := json.Marshal(variables)
	es = append(es, &open_api.EventResponse{
		Event:     "set",
		Namespace: "variable",
		Key:       namespace,
		Data:      data,
	})
	for _, id := range affectIds {
		info, has := oe.workerData.GetInfo(id)
		if !has {
			log.DebugF("worker(%s) not found", id)
			continue
		}
		es = append(es, &open_api.EventResponse{
			Event:     "set",
			Namespace: eosc.NamespaceWorker,
			Key:       info.config.Id,
			Data:      info.Body(),
		})
	}

	return http.StatusOK, nil, es, nil
}
