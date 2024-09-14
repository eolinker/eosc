package api_http

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/process-admin/admin"
	"github.com/eolinker/eosc/process-admin/marshal"

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

func (oe *WorkerApi) Add(r *http.Request, params httprouter.Params) (status int, header http.Header, body interface{}) {
	profession := params.ByName("profession")
	isSkip := true
	isSkip, status, header, body = oe.compatibleSetting(profession, r, params)
	if isSkip {
		return
	}
	decoder, err := marshal.GetData(r)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	cb := new(eosc.WorkerConfig)
	errUnmarshal := decoder.UnMarshal(cb)

	if errUnmarshal != nil {
		return http.StatusInternalServerError, nil, errUnmarshal
	}
	cb.Body, _ = decoder.Encode()
	if cb.Version == "" {
		cb.Version = admin.GenVersion()
	}
	cb.Profession = profession

	var out *admin.WorkerInfo
	err = oe.admin.Transaction(r.Context(), func(ctx context.Context, api admin.AdminApiWrite) error {
		worker, err := api.SetWorker(ctx, cb)
		if err != nil {
			return err
		}
		out = worker
		return nil
	})
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, nil, out.Detail()
}
func (oe *WorkerApi) Patch(r *http.Request, params httprouter.Params) (status int, header http.Header, body interface{}) {
	profession := params.ByName("profession")
	isSkip := true
	isSkip, status, header, body = oe.compatibleSetting(profession, r, params)
	if isSkip {
		return
	}
	name := params.ByName("name")
	if name == "" {
		return http.StatusInternalServerError, nil, "require name"
	}
	id, ok := eosc.ToWorkerId(name, profession)
	if !ok {
		return http.StatusInternalServerError, nil, fmt.Errorf("invalid name:%s", name)
	}
	decoder, err := marshal.GetData(r)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	options := make(map[string]interface{})
	err = decoder.UnMarshal(&options)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	if len(options) == 0 {
		return http.StatusInternalServerError, nil, "nothing to patch"

	}
	workerInfo, has := oe.admin.GetWorker(r.Context(), id)
	if !has {
		return 0, nil, "not exist"
	}
	current := make(map[string]interface{})
	_ = json.Unmarshal(workerInfo.Body(), &current)

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
	version := admin.GenVersion()
	if v, has := options["version"]; has {
		t, ok := v.(string)
		if ok {
			version = t
		}
	}
	matchLabels := workerInfo.Matches()
	if v, has := options["matchLabels"]; has {
		if vl, ok := v.(map[string]string); ok {
			matchLabels = vl
		}
	}
	data, _ := json.Marshal(current)
	cf := &eosc.WorkerConfig{
		Id:          id,
		Name:        name,
		Profession:  profession,
		Driver:      workerInfo.Driver(),
		Description: description,
		Version:     version,
		Body:        data,
		Matches:     matchLabels,
		Update:      eosc.Now(),
	}

	decoder = marshal.JsonData(data)
	err = oe.admin.Transaction(r.Context(), func(ctx context.Context, api admin.AdminApiWrite) error {
		w, err := api.SetWorker(ctx, cf)
		if err != nil {
			workerInfo = w
		}
		return err
	})

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, nil, workerInfo.Detail()
}
func (oe *WorkerApi) Save(r *http.Request, params httprouter.Params) (status int, header http.Header, body interface{}) {

	profession := params.ByName("profession")
	isSkip := true
	isSkip, status, header, body = oe.compatibleSetting(profession, r, params)
	if isSkip {
		return
	}
	name := params.ByName("name")
	if name == "" {
		return http.StatusInternalServerError, nil, "require name"
	}
	decoder, err := marshal.GetData(r)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	cb := new(eosc.WorkerConfig)
	errUnmarshal := decoder.UnMarshal(cb)
	if errUnmarshal != nil {
		return http.StatusInternalServerError, nil, errUnmarshal
	}

	if cb.Version == "" {
		cb.Version = admin.GenVersion()
	}
	cb.Profession = profession
	cb.Id, _ = eosc.ToWorkerId(name, profession)
	cb.Name = name
	cb.Update = eosc.Now()
	cb.Create = eosc.Now()
	cb.Body, _ = decoder.Encode()
	var out *admin.WorkerInfo
	err = oe.admin.Transaction(r.Context(), func(ctx context.Context, api admin.AdminApiWrite) error {
		w, err := api.SetWorker(ctx, cb)
		if err != nil {
			return err
		}
		out = w
		return nil
	})
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, nil, out.Detail()
}

func (oe *WorkerApi) Delete(r *http.Request, params httprouter.Params) (status int, header http.Header, body interface{}) {

	profession := params.ByName("profession")
	isSkip := true
	isSkip, status, header, body = oe.compatibleSetting(profession, r, params)
	if isSkip {
		return
	}
	name := params.ByName("name")
	id, ok := eosc.ToWorkerId(name, profession)
	if !ok {
		return http.StatusNotFound, nil, fmt.Sprintf("invalid name:%s for %s", name, profession)
	}
	p, has := oe.admin.GetProfession(r.Context(), profession)
	if !has {
		return http.StatusNotFound, nil, fmt.Sprintf("invalid profession:%s", profession)
	}
	if p.Mod == eosc.ProfessionConfig_Singleton {
		return http.StatusForbidden, nil, fmt.Sprintf("not allow delete %s for %s", name, profession)
	}

	var out *admin.WorkerInfo
	err := oe.admin.Transaction(r.Context(), func(ctx context.Context, api admin.AdminApiWrite) error {
		worker, err := api.DeleteWorker(ctx, id)
		if err != nil {
			return err
		}
		out = worker
		return nil
	})
	if err != nil {
		return http.StatusNotFound, nil, err
	}

	return http.StatusOK, nil, out.Detail()
}
