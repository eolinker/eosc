package process_admin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	open_api "github.com/eolinker/eosc/open-api"
	"github.com/eolinker/eosc/professions"
	"github.com/eolinker/eosc/utils/schema"
	"github.com/julienschmidt/httprouter"
)

type ProfessionApi struct {
	data       professions.IProfessions
	workerData *WorkerDatas
}

func NewProfessionApi(data professions.IProfessions, ws *WorkerDatas) *ProfessionApi {
	return &ProfessionApi{data: data, workerData: ws}
}

func (pi *ProfessionApi) Register(router *httprouter.Router) {
	router.Handle(http.MethodGet, "/profession", open_api.CreateHandleFunc(pi.All))
	router.Handle(http.MethodGet, "/profession/:profession", open_api.CreateHandleFunc(pi.Detail))
	router.Handle(http.MethodGet, "/profession/:profession/drivers", open_api.CreateHandleFunc(pi.Drivers))
	router.Handle(http.MethodGet, "/profession/:profession/driver", open_api.CreateHandleFunc(pi.DriverInfo))
	router.Handle(http.MethodPut, "/profession/:profession/drivers", open_api.CreateHandleFunc(pi.ResetDrivers))
	router.Handle(http.MethodPost, "/profession/:profession/drivers", open_api.CreateHandleFunc(pi.ResetDrivers))

	router.Handle(http.MethodPut, "/profession/:profession/driver", open_api.CreateHandleFunc(pi.SetDriver))
	router.Handle(http.MethodPost, "/profession/:profession/driver", open_api.CreateHandleFunc(pi.AddDriver))
	router.Handle(http.MethodDelete, "/profession/:profession/driver", open_api.CreateHandleFunc(pi.Delete))
	router.GET("/profession/:profession/skill", open_api.CreateHandleFunc(pi.Skill))

}

type ProfessionInfo struct {
	Name   string   `json:"name,omitempty"`
	Label  string   `json:"label,omitempty"`
	Desc   string   `json:"desc,omitempty"`
	Driver []string `json:"driver,omitempty"`
}

func (pi *ProfessionApi) Skill(req *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {
	name := params.ByName("profession")
	skill := req.URL.Query().Get("skill")
	if skill == "" {
		return http.StatusBadRequest, nil, nil, "skill invalid"
	}
	pn, has := pi.data.Get(name)
	if !has {
		return http.StatusNotFound, nil, nil, ErrorNotExist
	}
	dependencies := pn.Dependencies
	dps := make(map[string]bool)
	for _, dependency := range dependencies {
		dps[dependency] = true
	}
	ws := make([]interface{}, 0, pi.workerData.Count())
	all := pi.workerData.All()

	for _, w := range all {
		if w.worker != nil {
			if w.worker.CheckSkill(skill) {
				ws = append(ws, w.Info(pn.AppendLabels...))
			}
		}
		log.Debug("worker is: ", w.worker)
	}
	return http.StatusOK, nil, nil, ws
}

func (pi *ProfessionApi) All(r *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {
	list := pi.data.List()
	res := make([]*ProfessionInfo, 0, len(list))
	for _, p := range list {
		drivers := make([]string, 0, len(p.Drivers))
		for _, d := range p.Drivers {
			drivers = append(drivers, d.Name)
		}
		res = append(res, &ProfessionInfo{
			Name:   p.Name,
			Label:  p.Label,
			Desc:   p.Desc,
			Driver: drivers,
		})
	}
	return 200, nil, nil, res
}

func (pi *ProfessionApi) Detail(r *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {
	name := params.ByName("profession")
	profession, has := pi.data.Get(name)
	if !has {
		return http.StatusNotFound, nil, nil, ErrorNotExist
	}
	return 200, nil, nil, profession.ProfessionConfig
}

func (pi *ProfessionApi) Drivers(r *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {
	name := params.ByName("profession")
	profession, has := pi.data.Get(name)
	if !has {
		return http.StatusNotFound, nil, nil, ErrorNotExist
	}
	ds := profession.GetDrivers()
	type DriverInfo struct {
		Id    string `json:"id,omitempty"`
		Name  string `json:"name,omitempty"`
		Desc  string `json:"desc,omitempty"`
		Label string `json:"label,omitempty"`
	}
	dsi := make([]*DriverInfo, 0, len(ds))
	for _, d := range ds {
		dsi = append(dsi, &DriverInfo{
			Id:    d.Id,
			Name:  d.Name,
			Desc:  d.Desc,
			Label: d.Label,
		})
	}
	return http.StatusOK, nil, nil, dsi
}

func (pi *ProfessionApi) DriverInfo(r *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {
	professionName := params.ByName("profession")

	driverName := r.URL.Query().Get("name")
	if driverName == "" {
		return http.StatusInternalServerError, nil, nil, "invalid driver name"

	}
	profession, has := pi.data.Get(professionName)
	if !has {
		return http.StatusNotFound, nil, nil, ErrorNotExist
	}
	driverInfo, has := profession.DriverConfig(driverName)
	if !has {
		return http.StatusInternalServerError, nil, nil, fmt.Errorf("%s in %s:%w", driverName, professionName, ErrorNotExist)
	}
	type DriverDetail struct {
		*eosc.DriverConfig
		Render *schema.Schema `json:"render"`
	}
	df, ok := profession.GetDriver(driverName)
	if !ok {
		return http.StatusInsufficientStorage, nil, nil, ErrorExtenderNotWork
	}
	render, err := schema.Generate(df.ConfigType(), nil)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}

	return http.StatusOK, nil, nil, &DriverDetail{
		DriverConfig: driverInfo,
		Render:       render,
	}
}

func (pi *ProfessionApi) SetDriver(r *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {
	name := params.ByName("profession")
	driverName := r.URL.Query().Get("name")
	profession, has := pi.data.Get(name)
	if !has {
		return http.StatusNotFound, nil, nil, ErrorNotExist
	}
	driverConfig, has := profession.DriverConfig(driverName)
	if !has {
		return http.StatusInternalServerError, nil, nil, ErrorNotExist
	}

	bodyData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	err = json.Unmarshal(bodyData, driverConfig)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	pConfig := profession.ProfessionConfig

	err = pi.data.Set(name, pConfig)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	data, _ := json.Marshal(pConfig)
	return http.StatusOK, nil, []*open_api.EventResponse{{
		Event:     eosc.EventSet,
		Namespace: eosc.NamespaceProfession,
		Key:       name,
		Data:      data,
	}}, data

}
func (pi *ProfessionApi) AddDriver(r *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {
	name := params.ByName("profession")
	profession, has := pi.data.Get(name)
	if !has {
		return http.StatusNotFound, nil, nil, ErrorNotExist
	}
	bodyData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	driverConfig := new(eosc.DriverConfig)
	err = json.Unmarshal(bodyData, driverConfig)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}

	_, has = profession.DriverConfig(driverConfig.Name)
	if has {
		return http.StatusInsufficientStorage, nil, nil, "driver name duplicate:" + driverConfig.Name
	}
	pConfig := profession.ProfessionConfig
	pConfig.Drivers = append(profession.Drivers, driverConfig)

	err = pi.data.Set(name, pConfig)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	data, _ := json.Marshal(pConfig)
	return http.StatusOK, nil, []*open_api.EventResponse{{
		Event:     eosc.EventSet,
		Namespace: eosc.NamespaceProfession,
		Key:       name,
		Data:      data,
	}}, data

}
func (pi *ProfessionApi) ResetDrivers(r *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {
	name := params.ByName("profession")
	profession, has := pi.data.Get(name)
	if !has {
		return http.StatusNotFound, nil, nil, ErrorNotExist
	}

	bodyData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	driverConfigs := make([]*eosc.DriverConfig, 0, len(profession.Drivers))
	err = json.Unmarshal(bodyData, &driverConfigs)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	pConfig := profession.ProfessionConfig
	pConfig.Drivers = driverConfigs
	err = pi.data.Set(name, pConfig)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	data, _ := json.Marshal(pConfig)
	return http.StatusOK, nil, []*open_api.EventResponse{{
		Event:     eosc.EventSet,
		Namespace: eosc.NamespaceProfession,
		Key:       name,
		Data:      data,
	}}, data

}
func (pi *ProfessionApi) Delete(r *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {
	name := params.ByName("profession")
	driverName := r.URL.Query().Get("name")
	profession, has := pi.data.Get(name)
	if !has {
		return http.StatusNotFound, nil, nil, ErrorNotExist
	}

	_, has = profession.DriverConfig(driverName)
	if !has {
		return http.StatusNotFound, nil, nil, fmt.Sprintf("driver [%s] in %s not exits", driverName, name)
	}
	pConfig := profession.ProfessionConfig
	index := -1
	for i, d := range pConfig.Drivers {
		if d.Name == driverName {
			index = i
			break
		}
	}
	if index > -1 {
		pConfig.Drivers = append(pConfig.Drivers[:index], pConfig.Drivers[index+1:]...)
	}
	err := pi.data.Set(name, pConfig)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	data, _ := json.Marshal(pConfig)
	return http.StatusOK, nil, []*open_api.EventResponse{{
		Event:     eosc.EventSet,
		Namespace: eosc.NamespaceProfession,
		Key:       name,
		Data:      data,
	}}, data
}
