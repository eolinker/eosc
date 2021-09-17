package admin

import (
	"net/http"

	"github.com/eolinker/eosc/admin"
)

type WorkerInfo map[string]interface{}

type Admin struct {
	professions admin.IProfessions
	workers     admin.IWorkers
	handler     http.Handler
}

func (a *Admin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if a.handler == nil {
		http.NotFound(w, r)
		return
	}
	a.handler.ServeHTTP(w, r)
}

func NewAdmin(professions admin.IProfessions, workers admin.IWorkers, prefix string) *Admin {
	a := &Admin{professions: professions, workers: workers}
	a.handler = load(a, prefix)
	return a
}
