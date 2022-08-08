package process_admin

import (
	open_api "github.com/eolinker/eosc/open-api"
	"github.com/julienschmidt/httprouter"
)

type VariableApi struct {
}

func (oe *VariableApi) Register(router *httprouter.Router) {
	router.GET("/variable", open_api.CreateHandleFunc(nil))
	router.GET("/variable/:namespace", open_api.CreateHandleFunc(nil))
	router.POST("/variable/:namespace", open_api.CreateHandleFunc(nil))
	router.PUT("/variable/:namespace", open_api.CreateHandleFunc(nil))
	router.POST("/variable/:namespace/:key", open_api.CreateHandleFunc(nil))
}

func NewVariableApi() *VariableApi {

	return &VariableApi{}
}
