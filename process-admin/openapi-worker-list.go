package process_admin

import (
	"encoding/json"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	open_api "github.com/eolinker/eosc/open-api"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

func (oe *WorkerApi) getEmployeesByProfession(r *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {
	profession := params.ByName("profession")
	isSkip := true
	isSkip, status, header, events, body = oe.compatibleSetting(profession, r, params)
	if isSkip {
		return
	}
	es, err := oe.admin.ListWorker(r.Context(), profession)
	if err != nil {
		return 500, nil, nil, err
	}

	out, _ := json.Marshal(es)
	log.Debug("getEmployeesByProfession:", string(out))
	return 200, nil, nil, out
}

func (oe *WorkerApi) getEmployeeByName(r *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {
	profession := params.ByName("profession")

	isSkip := true
	isSkip, status, header, events, body = oe.compatibleSetting(profession, r, params)
	if isSkip {
		return
	}

	name := params.ByName("name")
	id, ok := eosc.ToWorkerId(name, profession)
	if !ok {
		return http.StatusServiceUnavailable, nil, nil, "worker id is invalid"
	}
	eo, err := oe.admin.GetWorker(r.Context(), id)
	if err != nil {
		return 404, nil, nil, err
	}
	return 200, nil, nil, eo.Detail()
}

func (oe *WorkerApi) compatibleSetting(profession string, r *http.Request, params httprouter.Params) (isSkip bool, status int, header http.Header, events []*open_api.EventResponse, body interface{}) {
	if strings.ToLower(profession) != admin.Setting {
		isSkip = false
		return
	}
	isSkip = true
	status, header, events, body = oe.settingRequest(r, params)
	return
}
