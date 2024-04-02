package process_admin

import (
	open_api "github.com/eolinker/eosc/open-api"
	admin_o "github.com/eolinker/eosc/process-admin/admin-o"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type WorkerApi struct {
	admin admin_o.AdminController

	settingRequest func(req *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{})
}

func (oe *WorkerApi) Register(router *httprouter.Router) {
	router.GET("/api/:profession", open_api.CreateHandleFunc(oe.getEmployeesByProfession))
	router.GET("/api/:profession/:name", open_api.CreateHandleFunc(oe.getEmployeeByName))
	router.POST("/api/:profession", open_api.CreateHandleFunc(oe.Add))
	router.PUT("/api/:profession/:name", open_api.CreateHandleFunc(oe.Save))
	router.POST("/api/:profession/:name", open_api.CreateHandleFunc(oe.Save))
	router.DELETE("/api/:profession/:name", open_api.CreateHandleFunc(oe.Delete))
	router.PATCH("/api/:profession/:name", open_api.CreateHandleFunc(oe.Patch))

}

func NewWorkerApi(workers admin_o.AdminController, settingRequest func(req *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{})) *WorkerApi {

	return &WorkerApi{admin: workers, settingRequest: settingRequest}
}
