package process_admin

import (
	"fmt"
	admin_o "github.com/eolinker/eosc/process-admin/admin-o"
	"github.com/eolinker/eosc/utils"
	"net/http"
	"time"

	"github.com/eolinker/eosc/log"
	open_api "github.com/eolinker/eosc/open-api"
	"github.com/eolinker/eosc/utils/zip"
	"github.com/ghodss/yaml"
	"github.com/julienschmidt/httprouter"
)

type ExportApi struct {
	version      map[string]string
	adminHandler admin_o.AdminController
}

func NewExportApi(extenders *ExtenderData, adminHandler admin_o.AdminController) *ExportApi {
	return &ExportApi{adminHandler: adminHandler, version: extenders.versions()}
}

func (oe *ExportApi) Register(router *httprouter.Router) {
	router.GET("/export", open_api.CreateHandleFunc(oe.export))

}
func (oe *ExportApi) export(r *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {

	workerData := oe.adminHandler.AllWorkers(r.Context())
	extenderList := oe.version
	professionList := oe.adminHandler.ListProfession(r.Context())

	id := time.Now().Format("2006-01-02 150405")

	exportData := getExportData(utils.GroupBy(workerData, admin_o.GetProfession))

	extenderData := make([]interface{}, 0, len(extenderList))
	for k, v := range extenderList {
		if v != "inner" {
			extenderData = append(extenderData, fmt.Sprintf("%s:%s", k, v))
		}
	}
	exportData["extenders"], _ = yamlEncode("extenders", extenderData)

	professionData := make([]interface{}, 0, len(professionList))
	for _, p := range professionList {
		professionData = append(professionData, p.ProfessionConfig)
	}
	exportData["professions"], _ = yamlEncode("professions", professionData)
	fileName := fmt.Sprintf("export_%s.zip", id)
	content, err := zip.CompressFile(exportData)
	if err != nil {
		return 500, nil, nil, err
	}
	header = make(http.Header)
	header.Add("Content-Type", "application/octet-stream")
	header.Add("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	return 200, header, nil, content

}
func yamlEncode[T any](k string, v []T) ([]byte, error) {
	newValue := map[string][]T{
		k: v,
	}
	d, err := yaml.Marshal(newValue)
	if err != nil {
		log.Errorf("marshal error	%s	%s", k, err.Error())
		return nil, err
	}
	return d, nil
}
func getExportData(value map[string][]*admin_o.WorkerInfo) map[string][]byte {
	data := make(map[string][]byte)
	for k, vs := range value {
		utils.ArrayType(vs, func(t *admin_o.WorkerInfo) any {
			return t.Detail()
		})
		data[fmt.Sprintf("profession-%s", k)], _ = yamlEncode(k, vs)
	}
	return data
}
