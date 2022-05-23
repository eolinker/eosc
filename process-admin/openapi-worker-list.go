package process_admin

import (
	open_api "github.com/eolinker/eosc/open-api"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (oe *WorkerApi) getEmployeesByProfession(r *http.Request, params httprouter.Params) (status int, header http.Header, event *open_api.EventResponse, body interface{}) {
	profession := params.ByName("profession")
	if profession == "_export" {
		return oe.export(r, params)
	}
	es, err := oe.workers.ListEmployees(profession)
	if err != nil {
		return 500, nil, nil, err
	}
	rs := make([]interface{}, 0, len(es))
	for _, e := range es {
		rs = append(rs, e.toAttr())
	}
	return 200, nil, nil, es
}

func (oe *WorkerApi) getEmployeeByName(r *http.Request, params httprouter.Params) (status int, header http.Header, event *open_api.EventResponse, body interface{}) {
	profession := params.ByName("profession")
	name := params.ByName("name")
	eo, err := oe.workers.GetEmployee(profession, name)
	if err != nil {
		return 404, nil, nil, err
	}
	return 200, nil, nil, eo.toAttr()
}
