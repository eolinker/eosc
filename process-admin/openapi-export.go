package process_admin

import (
	"fmt"
	"github.com/eolinker/eosc/log"
	open_api "github.com/eolinker/eosc/open-api"
	"github.com/eolinker/eosc/professions"
	"github.com/eolinker/eosc/utils/zip"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/yaml.v3"
	"net/http"
	"time"
)

type ExportApi struct {
	extenders  *ExtenderData
	workers    *Workers
	profession professions.IProfessions
}

func NewExportApi(extenders *ExtenderData, profession professions.IProfessions, workers *Workers) *ExportApi {
	return &ExportApi{extenders: extenders, workers: workers, profession: profession}
}

func (oe *ExportApi) Register(router *httprouter.Router) {
	router.GET("/export", open_api.CreateHandleFunc(oe.export))

}
func (oe *ExportApi) export(r *http.Request, params httprouter.Params) (status int, header http.Header, events []*open_api.EventResponse, body interface{}) {
	workerData := oe.allWorker()
	extenderList := oe.extenders.versions()
	professionList := oe.profession.List()

	id := time.Now().Format("2006-01-02 150405")

	exportData := getExportData(workerData)

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
func yamlEncode(k string, v []interface{}) ([]byte, error) {
	newValue := map[string][]interface{}{
		k: v,
	}
	d, err := yaml.Marshal(newValue)
	if err != nil {
		log.Errorf("marshal error	%s	%s", k, err.Error())
		return nil, err
	}
	return d, nil
}
func getExportData(value map[string][]interface{}) map[string][]byte {
	data := make(map[string][]byte)
	for k, v := range value {
		data[fmt.Sprintf("profession-%s", k)], _ = yamlEncode(k, v)
	}
	return data
}

func (oe *ExportApi) allWorker() map[string][]interface{} {
	ps := oe.workers.Export()
	data := make(map[string][]interface{})
	for key, pl := range ps {
		list := make([]interface{}, 0, len(pl))
		for _, p := range pl {
			list = append(list, p.Detail())
		}
		data[key] = list
	}
	return data
}
