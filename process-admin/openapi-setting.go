package process_admin

import (
	"context"
	"github.com/eolinker/eosc"
	open_api "github.com/eolinker/eosc/open-api"
	admin "github.com/eolinker/eosc/process-admin/admin"
	"github.com/eolinker/eosc/process-admin/marshal"
	"github.com/eolinker/eosc/setting"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type SettingApi struct {
	//workers  workers.IAdmin
	settings setting.ISettings
	variable eosc.IVariable

	data admin.AdminController
}

func (oe *SettingApi) RegisterSetting(router *httprouter.Router) {
	router.GET("/setting/:name", open_api.CreateHandleFunc(oe.Get))
	router.POST("/setting/:name", open_api.CreateHandleFunc(oe.Set))
	router.PUT("/setting/:name", open_api.CreateHandleFunc(oe.Set))
}
func (oe *SettingApi) request(req *http.Request, params httprouter.Params) (status int, header http.Header, body interface{}) {

	switch req.Method {
	case http.MethodGet:
		return oe.Get(req, params)
	case http.MethodPut, http.MethodPost:
		return oe.Set(req, params)
	}

	return http.StatusMethodNotAllowed, nil, http.StatusText(http.StatusMethodNotAllowed)
}
func (oe *SettingApi) Set(req *http.Request, params httprouter.Params) (status int, header http.Header, body interface{}) {
	name := params.ByName("name")
	idata, err := marshal.GetData(req)
	if err != nil {
		return http.StatusServiceUnavailable, nil, http.StatusText(http.StatusServiceUnavailable)
	}

	err = oe.data.Transaction(req.Context(), func(ctx context.Context, api admin.AdminApiWrite) error {
		return api.SetSetting(ctx, name, idata)
	})
	if err != nil {
		return http.StatusServiceUnavailable, nil, err.Error()
	}
	getSetting, has := oe.data.GetSetting(req.Context(), name)
	if has {
		return http.StatusOK, nil, getSetting
	}
	return http.StatusOK, nil, nil

}

func (oe *SettingApi) Get(req *http.Request, params httprouter.Params) (status int, header http.Header, body interface{}) {
	name := params.ByName("name")
	config, has := oe.data.GetSetting(req.Context(), name)
	if !has {
		return http.StatusNotFound, nil, http.StatusText(http.StatusNotFound)
	}

	return http.StatusOK, nil, config
}

func NewSettingApi(workers admin.AdminController) *SettingApi {

	return &SettingApi{
		data: workers,
	}
}
