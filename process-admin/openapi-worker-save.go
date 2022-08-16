package process_admin

import (
	"encoding/json"
	"fmt"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	open_api "github.com/eolinker/eosc/open-api"

	"github.com/julienschmidt/httprouter"
	"net/http"
)

type BaseArg struct {
	Id          string `json:"id,omitempty" yaml:"id"`
	Name        string `json:"name,omitempty" yaml:"name"`
	Driver      string `json:"driver,omitempty" yaml:"driver"`
	Description string `json:"description" yaml:"description"`
}

func (oe *WorkerApi) Add(r *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {
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

	obj, err := oe.workers.Update(profession, name, cb.Driver, cb.Description, decoder)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	eventData, _ := json.Marshal(obj.config)

	return http.StatusOK, nil, []*open_api.EventResponse{{
		Event:     eosc.EventSet,
		Namespace: eosc.NamespaceWorker,
		Key:       obj.config.Id,
		Data:      eventData,
	}}, obj.Detail()
}
func (oe *WorkerApi) Patch(r *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {
	profession := params.ByName("profession")
	name := params.ByName("name")
	if name == "" {
		return http.StatusInternalServerError, nil, nil, "require name"
	}
	decoder, err := GetData(r)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}

	options := make(map[string]interface{})
	err = decoder.UnMarshal(&options)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	if len(options) == 0 {
		return http.StatusInternalServerError, nil, nil, "nothing to patch"

	}
	workerInfo, err := oe.workers.GetEmployee(profession, name)
	if err != nil {
		return 0, nil, nil, nil
	}
	current := make(map[string]interface{})
	json.Unmarshal(workerInfo.config.Body, &current)

	for k, v := range options {
		if v != nil {
			log.Debug("patch set:", k, "=", v)
			current[k] = v
		} else {
			log.Debug("patch delete:", k)

			delete(current, k)
		}
	}
	description := workerInfo.config.Description
	if v, has := options["description"]; has {
		description = v.(string)
	}
	data, _ := json.Marshal(current)
	log.Debug("patch betfor:", string(workerInfo.config.Body))
	log.Debug("patch after:", string(data))
	decoder = JsonData(data)
	obj, err := oe.workers.Update(profession, name, workerInfo.config.Driver, description, decoder)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}

	eventData, _ := json.Marshal(obj.config)
	return http.StatusOK, nil, []*open_api.EventResponse{{
		Event:     eosc.EventSet,
		Namespace: eosc.NamespaceWorker,
		Key:       obj.config.Id,
		Data:      eventData,
	}}, obj.Detail()
}
func (oe *WorkerApi) Save(r *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {

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
	obj, err := oe.workers.Update(profession, name, cb.Driver, cb.Description, decoder)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}

	eventData, _ := json.Marshal(obj.config)
	return http.StatusOK, nil, []*open_api.EventResponse{{
		Event:     eosc.EventSet,
		Namespace: eosc.NamespaceWorker,
		Key:       obj.config.Id,
		Data:      eventData,
	}}, obj.Detail()
}

func (oe *WorkerApi) Delete(r *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {

	profession := params.ByName("profession")
	name := params.ByName("name")
	id, ok := eosc.ToWorkerId(name, profession)
	if !ok {
		return http.StatusNotFound, nil, nil, fmt.Sprintf("invalid name:%s for %s", name, profession)
	}
	p, has := oe.workers.professions.Get(profession)
	if !has {
		return http.StatusNotFound, nil, nil, fmt.Sprintf("invalid profession:%s", profession)
	}
	if p.Mod == eosc.ProfessionConfig_Singleton {
		return http.StatusForbidden, nil, nil, fmt.Sprintf("not allow delete %s for %s", name, profession)

	}
	wInfo, err := oe.workers.Delete(id)
	if err != nil {
		return 404, nil, nil, err
	}
	return http.StatusOK, nil, []*open_api.EventResponse{{
		Event:     eosc.EventDel,
		Namespace: eosc.NamespaceWorker,
		Key:       id,
		Data:      nil,
	}}, wInfo.Detail()
}
