package open_api

import (
	"net/http"
)

type OpenApiWorkerProxy struct {
	workerHandler http.Handler
}

func NewOpenApiWorkerProxy(workerHandler http.Handler) *OpenApiWorkerProxy {
	p := &OpenApiWorkerProxy{
		workerHandler: workerHandler,
	}
	return p
}

func (p *OpenApiWorkerProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.workerHandler.ServeHTTP(w, r)
}
