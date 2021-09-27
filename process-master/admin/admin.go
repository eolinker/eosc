package admin

import (
	"net/http"

	"github.com/eolinker/eosc"
)

type WorkerInfo map[string]interface{}

type Admin struct {
	professions eosc.IProfessionsData
	workers     eosc.IWorkersData
	handler     http.Handler
}

func (a *Admin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if a.handler == nil {
		http.NotFound(w, r)
		return
	}
	a.handler.ServeHTTP(w, r)
}

func NewAdmin(professions eosc.IProfessionsData, workers eosc.IWorkersData, prefix string) *Admin {
	a := &Admin{professions: professions, workers: workers}
	a.handler = load(a, prefix)
	return a
}