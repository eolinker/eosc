package process_admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	open_api "github.com/eolinker/eosc/open-api"
	"github.com/eolinker/eosc/setting"
	"github.com/julienschmidt/httprouter"
)

type SettingApi struct {
	workers  *Workers
	settings setting.ISettings
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
	driver, has := oe.settings.GetDriver(name)
	if !has {
		return http.StatusNotFound, nil, nil, http.StatusText(http.StatusNotFound)
	}
	if driver.Mode() == eosc.SettingModeReadonly {
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
	configType := driver.ConfigType()
	if driver.Mode() == eosc.SettingModeSingleton {
		err := oe.settings.SettingWorker(name, inputData, oe.variable)
		if err != nil {
			return http.StatusServiceUnavailable, nil, nil, err.Error()
		}
		wc := &eosc.WorkerConfig{
			Id:          fmt.Sprintf("%s@setting", name),
			Profession:  Setting,
			Name:        name,
			Driver:      name,
			Create:      eosc.Now(),
			Update:      eosc.Now(),
			Body:        inputData,
			Description: "",
		}
		eventData, _ := json.Marshal(wc)
		return http.StatusOK, nil, []*open_api.EventResponse{{
			Event:     eosc.EventSet,
			Namespace: eosc.NamespaceWorker,
			Key:       wc.Id,
			Data:      eventData,
		}}, setting.FormatConfig(inputData, configType)
	} else {
		events, body, err = oe.batchSet(inputData, driver, configType)
		if err != nil {
			status = http.StatusServiceUnavailable
			body = err.Error()
			log.Debug("batch set:", name, ":", string(inputData))
			log.Info("batch set:", name, ":", err)
			return
		}
		status = http.StatusOK
		return
	}
}

func (oe *SettingApi) batchSet(inputData []byte, driver eosc.ISetting, configType reflect.Type) ([]*open_api.EventResponse, interface{}, error) {
	type BatchWorkerInfo struct {
		id         string
		profession string
		name       string
		driver     string
		desc       string
		configBody IData
	}
	inputList := splitConfig(inputData)
	cfgs := make(map[string]BatchWorkerInfo, len(inputList))
	allWorkers := toSet(driver.AllWorkers())
	events := make([]*open_api.EventResponse, 0, len(allWorkers))
	responseBody := make([]interface{}, 0, len(inputList))
	for _, inp := range inputList {
		configData, _ := inp.Encode()
		cfg, _, err2 := oe.variable.Unmarshal(configData, configType)
		if err2 != nil {

			return nil, nil, err2
		}
		profession, workerName, driverName, desc, errCk := driver.Check(cfg)
		if errCk != nil {

			return nil, nil, errCk
		}
		id, _ := eosc.ToWorkerId(workerName, profession)
		if allWorkers[id] {
			delete(allWorkers, id)
		}
		cfgs[id] = BatchWorkerInfo{
			id:         id,
			profession: profession,
			name:       workerName,
			driver:     driverName,
			desc:       desc,
			configBody: inp,
		}
	}
	idtoDelete := make([]string, 0, len(allWorkers))
	for id := range allWorkers {
		idtoDelete = append(idtoDelete, id)
	}

	cannotDelete := oe.workers.DeleteTest(idtoDelete...)
	if len(cannotDelete) > 0 {
		return nil, nil, fmt.Errorf("should not delete:%s", strings.Join(cannotDelete, ","))
	}
	version := genVersion()
	for id, cfg := range cfgs {
		info, errSet := oe.workers.Update(cfg.profession, cfg.name, cfg.driver, version, cfg.desc, cfg.configBody)
		if errSet != nil {
			log.Warnf("bath set skip %s by error:%v", id, ":", errSet)
			continue
		}
		configData, _ := json.Marshal(info.config)
		responseBody = append(responseBody, info.Detail())
		events = append(events, &open_api.EventResponse{
			Event:     eosc.EventSet,
			Namespace: eosc.NamespaceWorker,
			Key:       id,
			Data:      configData,
		})
	}

	for _, id := range idtoDelete {
		oe.workers.Delete(id)
		events = append(events, &open_api.EventResponse{
			Event:     eosc.EventDel,
			Namespace: eosc.NamespaceWorker,
			Key:       id,
			Data:      nil,
		})
	}
	return events, responseBody, nil
}

func (oe *SettingApi) Get(req *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {
	name := params.ByName("name")
	_, has := oe.settings.GetDriver(name)
	if !has {
		return http.StatusNotFound, nil, nil, http.StatusText(http.StatusNotFound)
	}

	return http.StatusOK, nil, nil, oe.settings.GetConfig(name)
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
			datas.SettingWorker(name, config.Body, variable)
		}
	}
	return &SettingApi{
		workers:  workers,
		variable: variable,
		settings: datas,
	}
}

func toSet(ids []string) map[string]bool {
	s := make(map[string]bool)
	for _, id := range ids {
		s[id] = true
	}
	return s
}
