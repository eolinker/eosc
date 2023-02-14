package process_worker

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type MetricsApi struct {
}

func (ma *MetricsApi) Register(router *httprouter.Router) {
	router.GET("/apinto/metrics/:name", ma.GetMetrics)

}

func NewMetricsApi() *MetricsApi {
	return &MetricsApi{}
}

func (ma *MetricsApi) GetMetrics(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	w.Write([]byte(param.ByName("name")))
}
