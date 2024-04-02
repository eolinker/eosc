package process_admin

import (
	"encoding/json"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/process-admin/admin"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

func (oe *WorkerApi) getEmployeesByProfession(r *http.Request, params httprouter.Params) (status int, header http.Header, body interface{}) {
	profession := params.ByName("profession")
	isSkip := true
	isSkip, status, header, body = oe.compatibleSetting(profession, r, params)
	if isSkip {
		return
	}
	_, has := oe.admin.GetProfession(r.Context(), profession)
	if !has {
		return http.StatusNotFound, nil, "profession not found"
	}

	es, err := oe.admin.ListWorker(r.Context(), profession)
	if err != nil {
		return 500, nil, err
	}

	out, _ := json.Marshal(es)
	log.Debug("getEmployeesByProfession:", string(out))
	return 200, nil, out
}

func (oe *WorkerApi) getEmployeeByName(r *http.Request, params httprouter.Params) (status int, header http.Header, body interface{}) {
	profession := params.ByName("profession")

	isSkip := true
	isSkip, status, header, body = oe.compatibleSetting(profession, r, params)
	if isSkip {
		return
	}

	name := params.ByName("name")
	id, ok := eosc.ToWorkerId(name, profession)
	if !ok {
		return http.StatusServiceUnavailable, nil, "worker id is invalid"
	}
	eo, err := oe.admin.GetWorker(r.Context(), id)
	if err != nil {
		return 404, nil, err
	}
	return 200, nil, eo.Detail()
}

func (oe *WorkerApi) compatibleSetting(profession string, r *http.Request, params httprouter.Params) (isSkip bool, status int, header http.Header, body interface{}) {
	if strings.ToLower(profession) != admin.Setting {
		isSkip = false
		return
	}
	isSkip = true
	status, header, body = oe.settingRequest(r, params)
	return
}
