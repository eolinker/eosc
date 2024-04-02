package process_admin

import (
	"encoding/json"
	"fmt"
	"github.com/eolinker/eosc/log"
	open_api "github.com/eolinker/eosc/open-api"
	"github.com/julienschmidt/httprouter"
	"io"
	"mime"
	"net/http"
)

type ExtenderOpenApi struct {
	extenders *ExtenderData
}

func NewExtenderOpenApi(extenders *ExtenderData) *ExtenderOpenApi {
	return &ExtenderOpenApi{extenders: extenders}
}
func (oe *ExtenderOpenApi) Register(router *httprouter.Router) {

	router.Handle(http.MethodGet, "/extender", open_api.CreateHandleFunc(oe.List))
	router.Handle(http.MethodGet, "/extender/:id", open_api.CreateHandleFunc(oe.Info))
	router.Handle(http.MethodGet, "/extender/:id/:name", open_api.CreateHandleFunc(oe.Render))
	router.Handle(http.MethodPut, "/extender", open_api.CreateHandleFunc(oe.SET))
	router.Handle(http.MethodPost, "/extender", open_api.CreateHandleFunc(oe.SET))
	router.Handle(http.MethodDelete, "/extender/:id", open_api.CreateHandleFunc(oe.Delete))

}

func (oe *ExtenderOpenApi) Delete(r *http.Request, params httprouter.Params) (status int, header http.Header, body interface{}) {
	id := params.ByName("id")
	group, project := readProject(id)
	version := r.URL.Query().Get("v")

	projectInfo, err := oe.extenders.Delete(group, project, version)
	if err != nil {
		return 0, nil, err.Error()
	}

	return 200, nil, projectInfo.toInfo()

}
func (oe *ExtenderOpenApi) SET(r *http.Request, params httprouter.Params) (status int, header http.Header, body interface{}) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	log.Debug("Content-Type:", r.Header.Get("Content-Type"))
	log.Debug("mediaType:", mediaType, err)

	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Sprintf("Content-Type must application/json")
	}
	if mediaType != "application/json" {
		return http.StatusInternalServerError, nil, fmt.Sprintf("Content-Type must application/json")
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Debug(err)
		return http.StatusInternalServerError, nil, err.Error()
	}
	r.Body.Close()
	log.Debug("body:", string(data))
	type ParamSet struct {
		Group   string `json:"group"`
		Project string `json:"project"`
		Version string `json:"version"`
	}
	p := new(ParamSet)
	err = json.Unmarshal(data, p)
	if err != nil {
		log.Debug(err)
		return http.StatusInternalServerError, nil, err.Error()
	}
	log.Debug(p)
	projectInfo, ok, err := oe.extenders.SetVersion(p.Group, p.Project, p.Version)
	if err != nil {
		log.Debug(err)
		return http.StatusInternalServerError, nil, err.Error()
	}
	if ok {
		return 200, nil, projectInfo.toInfo()
	} else {
		return 200, nil, projectInfo.toInfo()
	}

}

func (oe *ExtenderOpenApi) List(r *http.Request, params httprouter.Params) (status int, header http.Header, body interface{}) {

	return 200, nil, oe.extenders.List()
}
func (oe *ExtenderOpenApi) Info(r *http.Request, params httprouter.Params) (status int, header http.Header, body interface{}) {
	id := params.ByName("id")

	info, ok := oe.extenders.GetInfo(readProject(id))
	if !ok {
		return 404, nil, fmt.Sprintf("extender{%s} not install", id)
	}
	return 200, nil, info
}
func (oe *ExtenderOpenApi) Render(r *http.Request, params httprouter.Params) (status int, header http.Header, body interface{}) {
	id := params.ByName("id")
	name := params.ByName("name")
	group, project := readProject(id)
	info, ok := oe.extenders.GetRender(group, project, name)
	if !ok {
		return 404, nil, fmt.Sprintf("extender{%s} not install", id)
	}
	if info.Render == nil {
		return http.StatusServiceUnavailable, nil, fmt.Sprintf("extender{%s:%s:%s} not work", group, project, name)
	}
	return 200, nil, info
}
