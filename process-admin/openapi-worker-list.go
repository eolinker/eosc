package process_admin

import (
	"encoding/json"
	"github.com/eolinker/eosc/log"
	open_api "github.com/eolinker/eosc/open-api"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (oe *WorkerApi) getEmployeesByProfession(r *http.Request, params httprouter.Params) (status int, header http.Header, event *open_api.EventResponse, body interface{}) {
	profession := params.ByName("profession")

	es, err := oe.workers.ListEmployees(profession)
	if err != nil {
		return 500, nil, nil, err
	}

	out, _ := json.Marshal(es)
	log.Debug("getEmployeesByProfession:", string(out))
	return 200, nil, nil, out
}

func (oe *WorkerApi) getEmployeeByName(r *http.Request, params httprouter.Params) (status int, header http.Header, event *open_api.EventResponse, body interface{}) {
	profession := params.ByName("profession")
	name := params.ByName("name")
	eo, err := oe.workers.GetEmployee(profession, name)
	if err != nil {
		return 404, nil, nil, err
	}
	return 200, nil, nil, eo.Detail()
}
