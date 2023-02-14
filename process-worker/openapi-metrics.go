package process_worker

import (
	routerworker "github.com/eolinker/eosc/process-worker/router-worker"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type MetricsApi struct {
}

func (ma *MetricsApi) Register(router *httprouter.Router) {
	router.GET(routerworker.RouterPrefix+"metrics/:name", ma.GetMetrics)

}

func NewMetricsApi() *MetricsApi {
	return &MetricsApi{}
}

func (ma *MetricsApi) GetMetrics(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	str := "rawPath=" + r.URL.Path
	w.Write([]byte(str))
}
