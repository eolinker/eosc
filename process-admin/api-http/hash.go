package api_http

import (
	"context"
	"fmt"
	open_api "github.com/eolinker/eosc/open-api"
	"github.com/eolinker/eosc/process-admin/admin"
	"github.com/eolinker/eosc/process-admin/marshal"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
)

type HashApi struct {
	admin admin.AdminController
}

func NewHashApi(admin admin.AdminController) *HashApi {
	return &HashApi{admin: admin}
}

func (h *HashApi) Register(router *httprouter.Router) {
	router.GET("/hash/:key", open_api.CreateHandleFunc(h.Get))
	router.POST("/hash/:key", open_api.CreateHandleFunc(h.PutKey))
	router.PUT("/hash/:key", open_api.CreateHandleFunc(h.PutKey))

	router.POST("/hash/:key/:field", open_api.CreateHandleFunc(h.PostField))
	router.PUT("/hash/:key/:field", open_api.CreateHandleFunc(h.PostField))

	router.DELETE("/hash/:key", open_api.CreateHandleFunc(h.DeleteKey))
	router.DELETE("/hash/:key/:field", open_api.CreateHandleFunc(h.DeleteField))
}

func (h *HashApi) Get(req *http.Request, params httprouter.Params) (status int, header http.Header, body interface{}) {
	key := params.ByName("key")
	if key == "" {
		return http.StatusBadRequest, nil, "key is required"
	}
	all, b := h.admin.GetHashAll(req.Context(), key)
	if !b {
		return http.StatusNotFound, nil, fmt.Sprintf("not found key: %s", key)
	}
	return http.StatusOK, nil, all
}

func (h *HashApi) PostField(req *http.Request, params httprouter.Params) (status int, header http.Header, body interface{}) {
	key := params.ByName("key")
	field := params.ByName("field")
	if key == "" || field == "" {
		return http.StatusBadRequest, nil, "key and field is required"
	}

	value, err := io.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	err = h.admin.Transaction(req.Context(), func(ctx context.Context, api admin.AdminApiWrite) error {
		return api.SetHashValue(ctx, key, field, string(value))
	})
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	all, _ := h.admin.GetHashAll(req.Context(), key)
	return http.StatusOK, nil, all
}

func (h *HashApi) PutKey(req *http.Request, params httprouter.Params) (status int, header http.Header, body interface{}) {
	key := params.ByName("key")
	if key == "" {
		return http.StatusBadRequest, nil, "key is required"
	}

	decoder, err := marshal.GetData(req)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	cb := make(map[string]string)
	errUnmarshal := decoder.UnMarshal(&cb)
	if errUnmarshal != nil {
		return http.StatusInternalServerError, nil, errUnmarshal
	}
	err = h.admin.Transaction(req.Context(), func(ctx context.Context, api admin.AdminApiWrite) error {
		return api.SetHash(ctx, key, cb)
	})
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, nil, cb
}

func (h *HashApi) DeleteKey(req *http.Request, params httprouter.Params) (status int, header http.Header, body interface{}) {
	key := params.ByName("key")
	if key == "" {
		return http.StatusBadRequest, nil, "key is required"
	}
	oldValue, has := h.admin.GetHashAll(req.Context(), key)
	if !has {
		return http.StatusNotFound, nil, fmt.Sprintf("not found key: %s", key)
	}
	err := h.admin.Transaction(req.Context(), func(ctx context.Context, api admin.AdminApiWrite) error {
		return api.DeleteHashAll(ctx, key)
	})
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, nil, oldValue
}

func (h *HashApi) DeleteField(req *http.Request, params httprouter.Params) (status int, header http.Header, body interface{}) {
	key := params.ByName("key")
	if key == "" {
		return http.StatusBadRequest, nil, "key is required"
	}
	field := params.ByName("field")
	if field == "" {
		return http.StatusBadRequest, nil, "field is required"
	}
	_, has := h.admin.GetHash(req.Context(), key, field)
	if !has {
		return http.StatusNotFound, nil, fmt.Sprintf("not found %s.%s", key, field)
	}
	err := h.admin.Transaction(req.Context(), func(ctx context.Context, api admin.AdminApiWrite) error {
		return api.DeleteHash(ctx, key, field)
	})
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	newValue, has := h.admin.GetHashAll(req.Context(), key)
	if !has {
		return http.StatusOK, nil, map[string]string{}
	}
	return http.StatusOK, nil, newValue
}
