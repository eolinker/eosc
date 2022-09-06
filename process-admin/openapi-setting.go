package process_admin

import (
	"encoding/json"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	open_api "github.com/eolinker/eosc/open-api"
	"github.com/eolinker/eosc/setting"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type SettingApi struct {
	workers  *Workers
	datas    setting.ISettings
	variable eosc.IVariable
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
	if driver.ReadOnly() {
		return http.StatusMethodNotAllowed, nil, nil, http.StatusText(http.StatusMethodNotAllowed)
	}

	idata, err := GetData(req)
	if err != nil {
		return http.StatusServiceUnavailable, nil, nil, http.StatusText(http.StatusServiceUnavailable)
	}
	inputData, err := idata.Encode()
	if err != nil {
		return http.StatusServiceUnavailable, nil, nil, http.StatusText(http.StatusServiceUnavailable)
	}
	output, toUpdate, toDelete, err := oe.datas.Set(name, inputData, oe.variable)
	if err != nil {
		return 0, nil, nil, nil
	}

	eventResponse := make([]*open_api.EventResponse, 0, len(toUpdate)+len(toDelete))
	for _, m := range toUpdate {
		eventData, _ := json.Marshal(m)
		eventResponse = append(eventResponse, &open_api.EventResponse{
			Event:     eosc.EventSet,
			Namespace: eosc.NamespaceWorker,
			Key:       m.Id,
			Data:      eventData,
		})
		oe.workers.set(m.Id, m.Profession, m.Name, m.Driver, m.Description, m.Body)
	}
	for _, delId := range toDelete {

		eventResponse = append(eventResponse, &open_api.EventResponse{
			Event:     eosc.EventDel,
			Namespace: eosc.NamespaceWorker,
			Key:       delId,
			Data:      nil,
		})
		oe.workers.Delete(delId)
	}

	return http.StatusOK, nil, eventResponse, output
}

func (oe *SettingApi) Get(req *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {
	name := params.ByName("name")
	_, has := oe.datas.GetDriver(name)
	if !has {
		return http.StatusNotFound, nil, nil, http.StatusText(http.StatusNotFound)
	}

	return http.StatusOK, nil, nil, oe.datas.GetConfig(name)
}

func NewSettingApi(init map[string][]byte, workers *Workers, variable eosc.IVariable) *SettingApi {
	datas := setting.GetSettings()
	for id, conf := range init {

		_, name, _ := eosc.SplitWorkerId(id)
		_, has := datas.GetDriver(name)
		log.Debug("init setting id: ", id, " conf: ", string(conf), " ", has)
		if has {
			config := new(eosc.WorkerConfig)
			err := json.Unmarshal(conf, config)
			if err != nil {
				log.Warn("init setting Unmarshal WorkerConfig:", err)
				continue
			}
			log.Debug("init setting id body: ", id, " conf: ", string(config.Body), " ", has)
			datas.Set(name, config.Body, variable)
		}
	}
	return &SettingApi{
		workers:  workers,
		variable: variable,
		datas:    datas,
	}
}
