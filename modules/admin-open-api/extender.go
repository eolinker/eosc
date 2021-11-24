package admin_open_api

import (
	"errors"
	"net/http"

	"github.com/eolinker/eosc/process-master/extenders"

	"github.com/julienschmidt/httprouter"
)

type ExtenderAdmin struct {
	data extenders.ITypedPlugin
}

func NewExtenderAdmin(data extenders.ITypedPlugin) *ExtenderAdmin {
	return &ExtenderAdmin{data: data}
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
