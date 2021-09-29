package admin_open_api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/eolinker/eosc/process-master/admin"

	"github.com/eolinker/eosc/log"

	"github.com/eolinker/eosc"
	"github.com/julienschmidt/httprouter"
)

var _ iOpenAdmin = (*OpenAdmin)(nil)

func CreateHandler() admin.CreateHandler {
	return new(createHandler)
}

type OpenAdmin struct {
	prefix string
	admin  eosc.IAdmin
}

func (o *OpenAdmin) export(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	data := o.all()
	if len(data) < 1 {
		writeResultError(w, 500, errors.New("no data"))
		return
	}
	id := time.Now().Format("2006-01-02 150405")

	exportData := getExportData(data)
	fileName := fmt.Sprintf("export_%s.zip", id)
	content, err := CompressFile(exportData)
	if err != nil {
		writeResultError(w, 500, err)
		return
	}

	w.Header().Add("Content-Type", "application/octet-stream")
	w.Header().Add("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	w.Write(content)
}

func (o *OpenAdmin) all() map[string][]interface{} {
	professions := o.admin.ListProfessions()
	data := make(map[string][]interface{})
	for _, p := range professions {
		ws, err := o.admin.ListEmployees(p.Name)
		if err != nil {
			log.Errorf("read data error	%s	%s", p.Name, err.Error())
			continue
		}
		data[p.Name] = ws
	}
	return data
}

func (o *OpenAdmin) delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	profession := params.ByName("profession")
	name := params.ByName("name")

	err := o.admin.Delete(profession, name)
	if err != nil {
		writeResultError(w, 404, err)
		return
	}
	writeResult(w, []byte("{}"))
}

func (o *OpenAdmin) genUrl(url string) string {
	u := strings.TrimPrefix(url, "/")
	return fmt.Sprintf("%s/%s", o.prefix, u)
}

func (o *OpenAdmin) getFields(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprintf(w, "getFields:%v", params)
}

func (o *OpenAdmin) getRenders(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	//
	//profession := params.ByName("profession")
	//
	//renders, err := o.admin.Renders(profession)
	//if err != nil {
	//	writeResultError(w, 500, err)
	//	return
	//}
	//writeResult(w, renders)
}

func (o *OpenAdmin) getEmployeesByProfession(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	profession := params.ByName("profession")

	es, err := o.admin.ListEmployees(profession)
	if err != nil {
		writeResult(w, err)
		return
	}
	writeResult(w, es)
}

func (o *OpenAdmin) getEmployeeByName(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	profession := params.ByName("profession")
	name := params.ByName("name")
	eo, err := o.admin.GetEmployee(profession, name)
	if err != nil {
		writeResultError(w, 404, err)
		return
	}
	writeResult(w, eo)
}

func (o *OpenAdmin) getDriversByProfession(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	profession := params.ByName("profession")

	ds, err := o.admin.Drivers(profession)
	if err != nil {
		writeResultError(w, 500, err)
		return
	}
	writeResult(w, ds)
}

func (o *OpenAdmin) getDriverInfo(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	profession := params.ByName("profession")
	driver := params.ByName("driver")
	ds, err := o.admin.DriverInfo(profession, driver)
	if err != nil {
		writeResultError(w, 500, err)
		return
	}
	writeResult(w, ds)
}

func (o *OpenAdmin) getDriversItemByProfession(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	profession := params.ByName("profession")

	ds, err := o.admin.DriversItem(profession)
	if err != nil {
		writeResultError(w, 500, err)
		return
	}
	writeResult(w, ds)
}
func (o *OpenAdmin) getRender(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	//profession := params.ByName("profession")
	//driver := params.ByName("driver")
	////render, err := o.admin.Render(profession, driver)
	//if err != nil {
	//	writeResultError(w, 500, err)
	//	return
	//}
	//writeResult(w, render)
}

func (o *OpenAdmin) Save(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	profession := params.ByName("profession")
	name := params.ByName("name")

	idata, err := GetData(r)
	if err != nil {
		writeResultError(w, 500, err)
		return
	}
	data, err := idata.Encode()
	if err != nil {
		writeResultError(w, 500, err)
		return
	}
	cb := new(baseConfig)
	errUnmarshal := idata.UnMarshal(cb)
	if errUnmarshal != nil {
		writeResultError(w, 500, errUnmarshal)
		return
	}
	if name == "" {
		name = cb.Name
	}

	if name == "" {
		writeResultError(w, 500, errors.New("require name"))
		return
	}

	err = o.admin.Update(profession, name, cb.Driver, data)
	if err != nil {
		writeResultError(w, 500, err)

		return
	}
	employee, err := o.admin.GetEmployee(profession, name)
	if err != nil {

		info := make(map[string]interface{})
		idata.UnMarshal(&info)

		info["profession"] = profession
		info["id"] = eosc.ToWorkerId(name, profession)
		info["create"] = eosc.Now()
		info["update"] = eosc.Now()
		employee = info
	}
	writeResult(w, employee)
}

func (o *OpenAdmin) getProfessions(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	professions := o.admin.ListProfessions()
	writeResult(w, professions)
}

type iOpenAdmin interface {
	// GET /
	getProfessions(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	// GET /:profession
	getEmployeesByProfession(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	// GET /:profession/:name
	getEmployeeByName(w http.ResponseWriter, r *http.Request, params httprouter.Params)

	//GET /:profession/_render
	//getRenders(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	////GET /:profession/_render/:driver
	//getRender(w http.ResponseWriter, r *http.Request, params httprouter.Params)

	//POST /:profession/
	//POST /:profession/:name
	//POST /:profession/:id
	//POST /:profession/:name
	//POST /:profession/:id
	Save(w http.ResponseWriter, r *http.Request, params httprouter.Params)

	//GET /:profession/:name/:fieldName
	getFields(w http.ResponseWriter, r *http.Request, params httprouter.Params)

	//GET /:profession/_driver
	getDriversByProfession(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	//GET /:profession/_driver/:driver
	getDriverInfo(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	//GET /:profession/_driver/item
	getDriversItemByProfession(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	//DELETE /:profession/:name
	//DELETE /:profession/:id
	delete(w http.ResponseWriter, r *http.Request, params httprouter.Params)

	genUrl(url string) string

	//GET /_export
	export(w http.ResponseWriter, r *http.Request, params httprouter.Params)
}

func NewOpenAdmin(admin eosc.IAdmin) *OpenAdmin {

	return &OpenAdmin{
		admin: admin,
	}
}

func (o *OpenAdmin) GenHandler() http.Handler {
	var openAdmin iOpenAdmin = o
	router := httprouter.New()
	router.GET(openAdmin.genUrl("/"), openAdmin.getProfessions)
	router.GET(openAdmin.genUrl("/:profession"), func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		switch params.ByName("profession") {
		case "_export":
			openAdmin.export(w, r, params)
		default:
			openAdmin.getEmployeesByProfession(w, r, params)
		}
	})
	router.GET(openAdmin.genUrl("/:profession/:action"), func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		action := params.ByName("action")
		switch action {
		case "_render":
			//openAdmin.getRenders(w, r, params)
			return
		case "_driver":
			openAdmin.getDriversByProfession(w, r, params)
		default:
			rename(params, "action", "name")
			openAdmin.getEmployeeByName(w, r, params)
		}
	})
	router.GET(openAdmin.genUrl("/:profession/:action/:key"), func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		action := params.ByName("action")
		switch action {
		case "_render":
			rename(params, "key", "driver")

			//openAdmin.getRender(w, r, params)

			return
		case "_driver":
			key := params.ByName("key")
			if strings.ToLower(key) == "item" {
				openAdmin.getDriversItemByProfession(w, r, params)
			} else {
				rename(params, "key", "driver")
				openAdmin.getDriverInfo(w, r, params)
			}
			return
		default:

			rename(params, "action", "name")
			rename(params, "key", "fieldName")
			openAdmin.getFields(w, r, params)
		}

	})

	router.POST(openAdmin.genUrl("/:profession"), openAdmin.Save)
	router.POST(openAdmin.genUrl("/:profession/:name"), openAdmin.Save)
	router.PUT(openAdmin.genUrl("/:profession/:name"), openAdmin.Save)
	router.DELETE(openAdmin.genUrl("/:profession/:name"), openAdmin.delete)
	return router
}

func rename(ps httprouter.Params, source, target string) {
	for i := range ps {
		if ps[i].Key == source {
			ps[i].Key = target
		}
	}
}
