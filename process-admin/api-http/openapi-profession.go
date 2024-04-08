package api_http

import (
	"context"
	"encoding/json"
	"fmt"
	admin "github.com/eolinker/eosc/process-admin/admin"
	"github.com/eolinker/eosc/process-admin/data"
	"github.com/eolinker/eosc/process-admin/model"
	"io"
	"net/http"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	open_api "github.com/eolinker/eosc/open-api"
	"github.com/eolinker/eosc/utils/schema"
	"github.com/julienschmidt/httprouter"
)

type ProfessionApi struct {
	data admin.AdminController
}

func NewProfessionApi(ws admin.AdminController) *ProfessionApi {
	return &ProfessionApi{data: ws}
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
	router.Handle(http.MethodDelete, "/profession/:profession/driver", open_api.CreateHandleFunc(pi.DeleteDriver))
	router.GET("/profession/:profession/skill", open_api.CreateHandleFunc(pi.Skill))

}

func (pi *ProfessionApi) Skill(req *http.Request, params httprouter.Params) (status int, header http.Header, body interface{}) {
	name := params.ByName("profession")
	skill := req.URL.Query().Get("skill")
	if skill == "" {
		return http.StatusBadRequest, nil, "skill invalid"
	}
	pn, has := pi.data.GetProfession(req.Context(), name)
	if !has {
		return http.StatusNotFound, nil, data.ErrorNotExist
	}
	dependencies := pn.Dependencies
	dps := make(map[string]bool)
	for _, dependency := range dependencies {
		dps[dependency] = true
	}
	all := pi.data.AllWorkers(req.Context())

	ws := make([]interface{}, 0, len(all))

	for _, w := range all {
		wk := w.GetWorker()
		if wk != nil {
			if wk.CheckSkill(skill) {
				ws = append(ws, w.Info(pn.AppendLabels...))
			}
		}
		log.Debug("worker is: ", wk)
	}
	return http.StatusOK, nil, ws
}

func (pi *ProfessionApi) All(r *http.Request, params httprouter.Params) (status int, header http.Header, body interface{}) {
	list := pi.data.ListProfession(r.Context())
	res := make([]*model.ProfessionInfo, 0, len(list))
	for _, p := range list {
		drivers := make([]string, 0, len(p.Drivers))
		for _, d := range p.Drivers {
			drivers = append(drivers, d.Name)
		}
		res = append(res, &model.ProfessionInfo{
			Name:   p.Name,
			Label:  p.Label,
			Desc:   p.Desc,
			Driver: drivers,
		})
	}
	return 200, nil, res
}

func (pi *ProfessionApi) Detail(r *http.Request, params httprouter.Params) (status int, header http.Header, body interface{}) {
	name := params.ByName("profession")
	profession, has := pi.data.GetProfession(r.Context(), name)
	if !has {
		return http.StatusNotFound, nil, data.ErrorNotExist
	}
	return 200, nil, profession.ProfessionConfig
}

func (pi *ProfessionApi) Drivers(r *http.Request, params httprouter.Params) (status int, header http.Header, body interface{}) {
	name := params.ByName("profession")
	profession, has := pi.data.GetProfession(r.Context(), name)
	if !has {
		return http.StatusNotFound, nil, data.ErrorNotExist
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
	return http.StatusOK, nil, dsi
}

func (pi *ProfessionApi) DriverInfo(r *http.Request, params httprouter.Params) (status int, header http.Header, body interface{}) {
	professionName := params.ByName("profession")

	driverName := r.URL.Query().Get("name")
	if driverName == "" {
		return http.StatusInternalServerError, nil, "invalid driver name"

	}
	profession, has := pi.data.GetProfession(r.Context(), professionName)
	if !has {
		return http.StatusNotFound, nil, data.ErrorNotExist
	}
	driverInfo, has := profession.DriverConfig(driverName)
	if !has {
		return http.StatusInternalServerError, nil, fmt.Errorf("%s in %s:%w", driverName, professionName, data.ErrorNotExist)
	}
	type DriverDetail struct {
		*eosc.DriverConfig
		Render *schema.Schema `json:"render"`
	}
	df, ok := profession.GetDriver(driverName)
	if !ok {
		return http.StatusInsufficientStorage, nil, data.ErrorExtenderNotWork
	}
	render, err := schema.Generate(df.ConfigType(), nil)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, nil, &DriverDetail{
		DriverConfig: driverInfo,
		Render:       render,
	}
}

func (pi *ProfessionApi) SetDriver(r *http.Request, params httprouter.Params) (status int, header http.Header, body interface{}) {
	name := params.ByName("profession")
	driverName := r.URL.Query().Get("name")
	profession, has := pi.data.GetProfession(r.Context(), name)
	if !has {
		return http.StatusNotFound, nil, data.ErrorNotExist
	}
	driverConfig, has := profession.DriverConfig(driverName)
	if !has {
		return http.StatusInternalServerError, nil, data.ErrorNotExist
	}

	bodyData, err := io.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	err = json.Unmarshal(bodyData, driverConfig)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	pConfig := profession.ProfessionConfig

	err = pi.data.Transaction(r.Context(), func(ctx context.Context, api admin.AdminApiWrite) error {
		return api.SetProfession(name, pConfig)
	})
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	data, _ := json.Marshal(pConfig)

	return http.StatusOK, nil, data

}
func (pi *ProfessionApi) AddDriver(r *http.Request, params httprouter.Params) (status int, header http.Header, body interface{}) {
	name := params.ByName("profession")
	profession, has := pi.data.GetProfession(r.Context(), name)
	if !has {
		return http.StatusNotFound, nil, data.ErrorNotExist
	}
	bodyData, err := io.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	driverConfig := new(eosc.DriverConfig)
	err = json.Unmarshal(bodyData, driverConfig)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	_, has = profession.DriverConfig(driverConfig.Name)
	if has {
		return http.StatusInsufficientStorage, nil, "driver name duplicate:" + driverConfig.Name
	}
	pConfig := profession.ProfessionConfig
	pConfig.Drivers = append(profession.Drivers, driverConfig)

	err = pi.data.Transaction(r.Context(), func(ctx context.Context, api admin.AdminApiWrite) error {
		return api.SetProfession(name, pConfig)
	})
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	data, _ := json.Marshal(pConfig)
	return http.StatusOK, nil, data

}
func (pi *ProfessionApi) ResetDrivers(r *http.Request, params httprouter.Params) (status int, header http.Header, body interface{}) {
	name := params.ByName("profession")
	profession, has := pi.data.GetProfession(r.Context(), name)
	if !has {
		return http.StatusNotFound, nil, data.ErrorNotExist
	}

	bodyData, err := io.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	driverConfigs := make([]*eosc.DriverConfig, 0, len(profession.Drivers))
	err = json.Unmarshal(bodyData, &driverConfigs)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	pConfig := profession.ProfessionConfig
	pConfig.Drivers = driverConfigs
	err = pi.data.Transaction(r.Context(), func(ctx context.Context, api admin.AdminApiWrite) error {
		return api.SetProfession(name, pConfig)
	})
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	data, _ := json.Marshal(pConfig)

	return http.StatusOK, nil, data

}
func (pi *ProfessionApi) DeleteDriver(r *http.Request, params httprouter.Params) (status int, header http.Header, body interface{}) {
	name := params.ByName("profession")
	driverName := r.URL.Query().Get("name")
	profession, has := pi.data.GetProfession(r.Context(), name)
	if !has {
		return http.StatusNotFound, nil, data.ErrorNotExist
	}

	_, has = profession.DriverConfig(driverName)
	if !has {
		return http.StatusNotFound, nil, fmt.Sprintf("driver [%s] in %s not exits", driverName, name)
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
	err := pi.data.Transaction(r.Context(), func(ctx context.Context, api admin.AdminApiWrite) error {
		return api.SetProfession(name, pConfig)
	})

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	data, _ := json.Marshal(pConfig)
	return http.StatusOK, nil, data

}
