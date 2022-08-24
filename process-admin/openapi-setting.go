package process_admin

import (
	"encoding/json"
	"fmt"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	open_api "github.com/eolinker/eosc/open-api"
	"github.com/eolinker/eosc/setting"
	"github.com/eolinker/eosc/variable"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type SettingApi struct {
	datas    setting.ISettings
	variable variable.IVariable
}

func (oe *SettingApi) RegisterSetting(router *httprouter.Router) {
	router.GET("/setting/:name", open_api.CreateHandleFunc(oe.Get))
	router.POST("/setting/:name", open_api.CreateHandleFunc(oe.Set))
	router.PUT("/setting/:name", open_api.CreateHandleFunc(oe.Set))
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
	inputData, err := idata.Encode()
	if err != nil {
		return http.StatusServiceUnavailable, nil, nil, http.StatusText(http.StatusServiceUnavailable)
	}

	err = driver.Set(inputData)
	if err != nil {
		return 0, nil, nil, nil
	}
	id, _ := eosc.ToWorkerId(name, Setting)
	eventData, _ := json.Marshal(eosc.WorkerConfig{
		Id:          id,
		Profession:  Setting,
		Name:        name,
		Driver:      name,
		Create:      eosc.Now(),
		Update:      eosc.Now(),
		Body:        inputData,
		Description: id,
	})

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

func NewSettingApi(init map[string][]byte, variable variable.IVariable) *SettingApi {
	datas := setting.GetSettings()

	for id, conf := range init {
		name, _, _ := eosc.SplitWorkerId(id)
		driver, has := datas.GetDriver(name)
		if has {
			config := new(eosc.WorkerConfig)
			err := json.Unmarshal(conf, config)
			if err != nil {
				log.Warn("init setting Unmarshal WorkerConfig:", err)
				continue
			}
			cfg, vs, err := variable.Unmarshal(config.Body, driver.ConfigType())
			if err != nil {
				log.Warn("init setting", name, " Unmarshal variable:", err)
				continue
			}

			err = driver.Set(cfg)
			if err != nil {
				log.Warn("init setting set:", err)

				continue
			}
			variable.SetVariablesById(id, vs)
		}
	}
	return &SettingApi{
		variable: variable,
		datas:    datas,
	}
}
