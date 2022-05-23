package process_admin

import (
	"fmt"
	"github.com/eolinker/eosc/log"
	open_api "github.com/eolinker/eosc/open-api"
	"github.com/eolinker/eosc/utils/zip"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/yaml.v3"
	"net/http"
	"time"
)

func (oe *WorkerApi) export(r *http.Request, params httprouter.Params) (status int, header http.Header, event *open_api.EventResponse, body interface{}) {
	data := oe.all()
	if len(data) < 1 {
		return 500, nil, nil, "no data"
	}
	id := time.Now().Format("2006-01-02 150405")

	exportData := getExportData(data)
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

func getExportData(value map[string][]interface{}) map[string][]byte {
	data := make(map[string][]byte)
	for k, v := range value {
		newValue := map[string][]interface{}{
			k: v,
		}
		d, err := yaml.Marshal(newValue)
		if err != nil {
			log.Errorf("marshal error	%s	%s", k, err.Error())
			continue
		}
		data[k] = d
	}
	return data
}

func (oe *WorkerApi) all() map[string][]interface{} {
	ps := oe.workers.Export()
	data := make(map[string][]interface{})
	for key, pl := range ps {
		list := make([]interface{}, 0, len(pl))
		for _, p := range pl {
			list = append(list, p.toAttr())
		}
		data[key] = list
	}
	return data
}
