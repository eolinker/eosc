package process_admin

import (
	"encoding/json"
	"fmt"
	"github.com/eolinker/eosc"
	open_api "github.com/eolinker/eosc/open-api"
	"github.com/eolinker/eosc/setting"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var (
	settingApi = NewSettingApi()
)

type SettingApi struct {
	datas setting.ISettings
}

func RegisterSetting(router *httprouter.Router) {
	router.GET("/setting/:name", open_api.CreateHandleFunc(settingApi.Get))
	router.POST("/setting/:name", open_api.CreateHandleFunc(settingApi.Set))
	router.PUT("/setting/:name", open_api.CreateHandleFunc(settingApi.Set))
}
func (oe *SettingApi) request(req *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {

	switch req.Method {
	case http.MethodGet:
		return oe.Get(req, params)
	case http.MethodPut, http.MethodPost:
		return oe.Set(req, params)
	}

	return http.StatusMethodNotAllowed, nil, nil, http.StatusText(http.StatusMethodNotAllowed)
}
func (oe *SettingApi) Set(req *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {
	name := params.ByName("name")
	driver, has := oe.datas.GetDriver(name)
	if !has {
		return http.StatusNotFound, nil, nil, http.StatusText(http.StatusNotFound)
	}

	idata, err := GetData(req)
	if err != nil {
		return http.StatusServiceUnavailable, nil, nil, http.StatusText(http.StatusServiceUnavailable)
	}
	encode, err := idata.Encode()
	if err != nil {
		return http.StatusServiceUnavailable, nil, nil, http.StatusText(http.StatusServiceUnavailable)
	}

	obj, err := driver.Set(encode)
	if err != nil {
		return 0, nil, nil, nil
	}
	eventData, _ := json.Marshal(obj)
	return http.StatusOK, nil, []*open_api.EventResponse{{
		Event:     eosc.EventSet,
		Namespace: eosc.NamespaceWorker,
		Key:       fmt.Sprintf("%s@setting", name),
		Data:      eventData,
	}}, obj
}

func (oe *SettingApi) Get(req *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {
	name := params.ByName("name")
	driver, has := oe.datas.GetDriver(name)
	if !has {
		return http.StatusNotFound, nil, nil, http.StatusText(http.StatusNotFound)
	}

	return http.StatusOK, nil, nil, driver.Get()
}

func NewSettingApi() *SettingApi {

	return &SettingApi{
		datas: setting.GetSettings(),
	}
}
