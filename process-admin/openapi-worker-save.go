package process_admin

import (
	"encoding/json"
	"fmt"
	"github.com/eolinker/eosc/process-admin/marshal"
	"time"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	open_api "github.com/eolinker/eosc/open-api"

	"net/http"

	"github.com/julienschmidt/httprouter"
)

type BaseArg struct {
	Id          string `json:"id,omitempty" yaml:"id"`
	Name        string `json:"name,omitempty" yaml:"name"`
	Driver      string `json:"driver,omitempty" yaml:"driver"`
	Description string `json:"description" yaml:"description"`
	Version     string `json:"version" yaml:"version"`
}

//func NewBaseArg() *BaseArg {
//	return &BaseArg{
//		Version: genVersion(),
//	}
//}

func genVersion() string {
	return time.Now().Format("20060102150405")
}

func (oe *WorkerApi) Add(r *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {
	profession := params.ByName("profession")
	isSkip := true
	isSkip, status, header, events, body = oe.compatibleSetting(profession, r, params)
	if isSkip {
		return
	}
	decoder, err := marshal.GetData(r)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	cb := new(BaseArg)
	errUnmarshal := decoder.UnMarshal(cb)
	if errUnmarshal != nil {
		return http.StatusInternalServerError, nil, nil, errUnmarshal
	}

	name := cb.Name
	if cb.Version == "" {
		cb.Version = genVersion()
	}
	transaction := oe.admin.Begin(r.Context())

	obj, err := transaction.Update(profession, name, cb.Driver, cb.Version, cb.Description, decoder)
	if err != nil {
		transaction.Rollback()

		return http.StatusInternalServerError, nil, nil, err
	}
	transaction.Commit()

	return http.StatusOK, nil, []*open_api.EventResponse{{
		Event:     eosc.EventSet,
		Namespace: eosc.NamespaceWorker,
		Key:       obj.Id(),
		Data:      obj.ConfigData(),
	}}, obj.Detail()
}
func (oe *WorkerApi) Patch(r *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {
	profession := params.ByName("profession")
	isSkip := true
	isSkip, status, header, events, body = oe.compatibleSetting(profession, r, params)
	if isSkip {
		return
	}
	name := params.ByName("name")
	if name == "" {
		return http.StatusInternalServerError, nil, nil, "require name"
	}
	decoder, err := marshal.GetData(r)
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

	workerInfo, err := oe.admin.GetEmployee(profession, name)
	if err != nil {
		return 0, nil, nil, "not exist"
	}
	current := make(map[string]interface{})
	json.Unmarshal(workerInfo.Body(), &current)

	for k, v := range options {
		if v != nil {
			log.Debug("patch set:", k, "=", v)
			current[k] = v
		} else {
			log.Debug("patch delete:", k)

			delete(current, k)
		}
	}
	description := workerInfo.Description()
	if v, has := options["description"]; has {
		description = v.(string)
	}
	data, _ := json.Marshal(current)
	log.Debug("patch betfor:", string(workerInfo.Body()))
	log.Debug("patch after:", string(data))
	version := genVersion()
	if v, has := options[version]; has {
		t, ok := v.(string)
		if ok {
			version = t
		}
	}
	decoder = marshal.JsonData(data)
	transaction := oe.admin.Begin(r.Context())
	obj, err := transaction.Update(profession, name, workerInfo.Driver(), version, description, decoder)
	if err != nil {
		transaction.Rollback()
		return http.StatusInternalServerError, nil, nil, err
	}
	transaction.Commit()
	return http.StatusOK, nil, []*open_api.EventResponse{{
		Event:     eosc.EventSet,
		Namespace: eosc.NamespaceWorker,
		Key:       obj.Id(),
		Data:      obj.ConfigData(),
	}}, obj.Detail()
}
func (oe *WorkerApi) Save(r *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {

	profession := params.ByName("profession")
	isSkip := true
	isSkip, status, header, events, body = oe.compatibleSetting(profession, r, params)
	if isSkip {
		return
	}
	name := params.ByName("name")
	if name == "" {
		return http.StatusInternalServerError, nil, nil, "require name"
	}
	decoder, err := marshal.GetData(r)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	cb := new(BaseArg)
	errUnmarshal := decoder.UnMarshal(cb)
	if errUnmarshal != nil {
		return http.StatusInternalServerError, nil, nil, errUnmarshal
	}
	if cb.Version == "" {
		cb.Version = genVersion()
	}
	transaction := oe.admin.Begin(r.Context())

	obj, err := transaction.Update(profession, name, cb.Driver, cb.Version, cb.Description, decoder)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	transaction.Commit()
	return http.StatusOK, nil, []*open_api.EventResponse{{
		Event:     eosc.EventSet,
		Namespace: eosc.NamespaceWorker,
		Key:       obj.Id(),
		Data:      obj.ConfigData(),
	}}, obj.Detail()
}

func (oe *WorkerApi) Delete(r *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {

	profession := params.ByName("profession")
	isSkip := true
	isSkip, status, header, events, body = oe.compatibleSetting(profession, r, params)
	if isSkip {
		return
	}
	name := params.ByName("name")
	id, ok := eosc.ToWorkerId(name, profession)
	if !ok {
		return http.StatusNotFound, nil, nil, fmt.Sprintf("invalid name:%s for %s", name, profession)
	}
	p, has := oe.admin.GetProfession(profession)
	if !has {
		return http.StatusNotFound, nil, nil, fmt.Sprintf("invalid profession:%s", profession)
	}
	if p.Mod == eosc.ProfessionConfig_Singleton {
		return http.StatusForbidden, nil, nil, fmt.Sprintf("not allow delete %s for %s", name, profession)
	}
	transaction := oe.admin.Begin(r.Context())
	wInfo, err := transaction.Delete(id)
	if err != nil {
		transaction.Rollback()
		return 404, nil, nil, err
	}
	transaction.Commit()
	return http.StatusOK, nil, []*open_api.EventResponse{{
		Event:     eosc.EventDel,
		Namespace: eosc.NamespaceWorker,
		Key:       id,
		Data:      nil,
	}}, wInfo.Detail()
}
