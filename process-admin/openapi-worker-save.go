package process_admin

import (
	"encoding/json"
	"github.com/eolinker/eosc"
	open_api "github.com/eolinker/eosc/open-api"

	"github.com/julienschmidt/httprouter"
	"net/http"
)

type BaseArg struct {
	Id string `json:"id,omitempty" yaml:"id"`
	//Profession string `json:"profession,omitempty" yaml:"profession"`
	Name   string `json:"name,omitempty" yaml:"name"`
	Driver string `json:"driver,omitempty" yaml:"driver"`
}

func (oe *WorkerApi) Add(r *http.Request, params httprouter.Params) (status int, header http.Header, event *open_api.EventResponse, body interface{}) {
	profession := params.ByName("profession")
	decoder, err := GetData(r)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	cb := new(BaseArg)
	errUnmarshal := decoder.UnMarshal(cb)
	if errUnmarshal != nil {
		return http.StatusInternalServerError, nil, nil, errUnmarshal
	}

	name := cb.Name

	obj, err := oe.workers.Update(profession, name, cb.Driver, decoder)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	eventData, _ := json.Marshal(obj)

	return http.StatusOK, nil, &open_api.EventResponse{
		Event:     eosc.EventSet,
		Namespace: eosc.NamespaceWorker,
		Key:       obj.config.Id,
		Data:      eventData,
	}, obj.toAttr()
}
func (oe *WorkerApi) Save(r *http.Request, params httprouter.Params) (status int, header http.Header, event *open_api.EventResponse, body interface{}) {

	profession := params.ByName("profession")
	name := params.ByName("name")
	if name == "" {
		return http.StatusInternalServerError, nil, nil, "require name"
	}
	decoder, err := GetData(r)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	cb := new(BaseArg)
	errUnmarshal := decoder.UnMarshal(cb)
	if errUnmarshal != nil {
		return http.StatusInternalServerError, nil, nil, errUnmarshal
	}

	obj, err := oe.workers.Update(profession, name, cb.Driver, decoder)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	eventData, _ := json.Marshal(obj)

	return http.StatusOK, nil, &open_api.EventResponse{
		Event:     eosc.EventSet,
		Namespace: eosc.NamespaceWorker,
		Key:       obj.config.Id,
		Data:      eventData,
	}, nil
}

func (oe *WorkerApi) delete(r *http.Request, params httprouter.Params) (status int, header http.Header, event *open_api.EventResponse, body interface{}) {

	profession := params.ByName("profession")
	name := params.ByName("name")
	id, ok := eosc.ToWorkerId(name, profession)
	if !ok {
		return http.StatusNotFound, nil, nil, "invalid id"
	}
	wInfo, err := oe.workers.Delete(profession, name)
	if err != nil {
		return 404, nil, nil, err
	}
	return http.StatusOK, nil, &open_api.EventResponse{
		Event:     eosc.EventDel,
		Namespace: eosc.NamespaceWorker,
		Key:       id,
		Data:      nil,
	}, wInfo.toAttr()
}