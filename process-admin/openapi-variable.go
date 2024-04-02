package process_admin

import (
	"context"
	"fmt"
	admin_o "github.com/eolinker/eosc/process-admin/admin-o"
	"github.com/eolinker/eosc/process-admin/marshal"
	"net/http"

	open_api "github.com/eolinker/eosc/open-api"
	"github.com/julienschmidt/httprouter"
)

type VariableApi struct {
	adminHandler admin_o.AdminController
}

func NewVariableApi(adminHandler admin_o.AdminController) *VariableApi {
	return &VariableApi{
		adminHandler: adminHandler,
	}
}

func (oe *VariableApi) Register(router *httprouter.Router) {
	router.GET("/variable", open_api.CreateHandleFunc(oe.getAll))
	router.GET("/variable/:namespace", open_api.CreateHandleFunc(oe.getByNamespace))
	router.GET("/variable/:namespace/:key", open_api.CreateHandleFunc(oe.getByKey))
	router.POST("/variable/:namespace", open_api.CreateHandleFunc(oe.setByNamespace))
	router.PUT("/variable/:namespace", open_api.CreateHandleFunc(oe.setByNamespace))

}

func (oe *VariableApi) getAll(r *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {

	status = http.StatusOK
	events, _ = oe.adminHandler.Transaction(r.Context(), func(ctx context.Context, api admin_o.AdminApiWrite) error {
		body = api.AllVariables(ctx)
		return nil
	})

	return
}

func (oe *VariableApi) getByNamespace(r *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {
	namespace := params.ByName("namespace")
	if namespace == "" {
		namespace = "default"
	}
	data, has := oe.adminHandler.GetVariables(r.Context(), namespace)
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
	value, has := oe.adminHandler.GetVariable(r.Context(), namespace, key)
	if !has {
		return http.StatusNotFound, nil, nil, fmt.Sprintf("namespace{%s} not found", namespace)
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
	event, err := oe.adminHandler.Transaction(r.Context(), func(ctx context.Context, api admin_o.AdminApiWrite) error {
		return api.SetVariable(ctx, namespace, cb)
	})
	if err != nil {
		return http.StatusInternalServerError, nil, nil, errUnmarshal
	}
	return http.StatusOK, nil, event, map[string]interface{}{
		"namespace": namespace,
		"variables": cb,
	}

}
