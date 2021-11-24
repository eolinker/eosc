package admin_open_api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/eolinker/eosc/process-master/extenders"

	"github.com/julienschmidt/httprouter"
)

type ExtenderAdmin struct {
	pre  string
	data extenders.ITypedPlugin
}

func (e *ExtenderAdmin) GenHandler() http.Handler {
	router := httprouter.New()
	router.GET(e.pre+"/extender/:id", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		e.ExtenderInfo(w, r, params)
		return
	})
	router.GET(e.pre+"/extenders", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		e.Extenders(w, r, params)
		return
	})
	return router
}

func NewExtenderAdmin(pre string, data extenders.ITypedPlugin) *ExtenderAdmin {
	return &ExtenderAdmin{pre: strings.TrimSuffix(pre, "/"), data: data}
}

func (e *ExtenderAdmin) ExtenderInfo(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	plugin, has := e.data.GetPluginByID(id)
	if !has {
		writeResultError(w, 200, errors.New("the plugin does not exist"))
		return
	}
	writeResult(w, plugin)
}

func (e *ExtenderAdmin) Extenders(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	plugins := e.data.GetPlugins()

	writeResult(w, plugins)
}
