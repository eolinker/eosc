package admin_open_api

import (
	"errors"
	"fmt"
	"github.com/eolinker/eosc"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"strings"
)
var _ iOpenAdmin = (*OpenAdmin)(nil)

type OpenAdmin struct {
	prefix string
	admin  eosc.IAdmin
}

func (o *OpenAdmin) delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	profession:=params.ByName("profession")
	name:=params.ByName("name")

	info,err:=o.admin.Delete(profession,name)
	if err!= nil{
		writeResultError(w,500,err)
		return
	}
	writeResult(w,info)
}

func (o *OpenAdmin) genUrl(url string) string{
	u := strings.TrimPrefix(url,"/")
	return fmt.Sprintf("%s/%s",o.prefix,u)
}

func (o *OpenAdmin) getFields(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprintf(w,"getFields:%v",params)
}

func (o *OpenAdmin) getRenders(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	profession:=params.ByName("profession")

	renders,err:= o.admin.Renders(profession)
	if  err!= nil{
		writeResultError(w,500,err)
		return
	}
	writeResult(w,renders)
}

func (o *OpenAdmin) getEmployeesByProfession(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	profession:=params.ByName("profession")

	es,err:= o.admin.ListEmployees(profession)
	if err!= nil{
		writeResult(w,err)
		return
	}
	writeResult(w,es)
}

func (o *OpenAdmin) getEmployeeByName(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	profession:=params.ByName("profession")
	name:= params.ByName("name")
	eo,err:= o.admin.GetEmployee(profession,name)
	if err!= nil{
		writeResult(w,err)
		return
	}
	writeResult(w,eo)
}

func (o *OpenAdmin) getDriversByProfession(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	profession:=params.ByName("profession")

	ds,err:= o.admin.Drivers(profession)
	if err!= nil{
		writeResultError(w,500,err)
		return
	}
	writeResult(w,ds)
}

func (o *OpenAdmin) getDriverInfo(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	profession:=params.ByName("profession")
	driver:=params.ByName("driver")
	ds,err:= o.admin.DriverInfo(profession,driver)
	if err!= nil{
		writeResultError(w,500,err)
		return
	}
	writeResult(w,ds)
}

func (o *OpenAdmin) getDriversItemByProfession(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	profession:=params.ByName("profession")

	ds,err:= o.admin.DriversItem(profession)
	if err!= nil{
		writeResultError(w,500,err)
		return
	}
	writeResult(w,ds)
}
func (o *OpenAdmin) getRender(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	profession:=params.ByName("profession")
	driver:=params.ByName("driver")
	render,err:= o.admin.Render(profession,driver)
	if  err!= nil{
		writeResultError(w,500,err)
		return
	}
	writeResult(w,render)
}

func (o *OpenAdmin) Save(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	profession:= params.ByName("profession")
	name:= params.ByName("name")
	data,err:=ioutil.ReadAll(r.Body)

	if err!= nil{

		writeResultError(w,500,err)
		return
	}

	idata:=JsonData(data)
	cb:=new(baseConfig)
	errUnmarshal:= idata.UnMarshal(cb)
	if errUnmarshal!= nil{
		writeResultError(w,500,err)
		return
	}
	if name == ""{
		name = cb.Name
	}
	if name == ""{
		writeResultError(w,500,errors.New("require name"))
		return
	}

	winfo,err:=o.admin.Update(profession,name,cb.Driver,JsonData(data))
	if err!= nil{
		writeResultError(w,500,err)

		return
	}

	writeResult(w,winfo)
}

func (o *OpenAdmin) getProfessions(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	professions := o.admin.ListProfessions()
	writeResult(w,professions)
}

type iOpenAdmin interface {
	// GET /
	getProfessions(w http.ResponseWriter,r *http.Request,params httprouter.Params)
	// GET /:profession
	getEmployeesByProfession(w http.ResponseWriter,r *http.Request,params httprouter.Params)
	// GET /:profession/:name
	getEmployeeByName(w http.ResponseWriter,r *http.Request,params httprouter.Params)

	//GET /:profession/_render
	getRenders(w http.ResponseWriter,r *http.Request,params httprouter.Params)
	//GET /:profession/_render/:driver
	getRender(w http.ResponseWriter,r *http.Request,params httprouter.Params)

	//POST /:profession/
	//POST /:profession/:name
	//POST /:profession/:id
	//POST /:profession/:name
	//POST /:profession/:id
	Save(w http.ResponseWriter,r *http.Request,params httprouter.Params)

	//GET /:profession/:name/:fieldName
	getFields(w http.ResponseWriter,r *http.Request,params httprouter.Params)

	//GET /:profession/_driver
	getDriversByProfession(w http.ResponseWriter,r *http.Request,params httprouter.Params)
	//GET /:profession/_driver/:driver
	getDriverInfo(w http.ResponseWriter,r *http.Request,params httprouter.Params)
	//GET /:profession/_driver/item
	getDriversItemByProfession(w http.ResponseWriter,r *http.Request,params httprouter.Params)
	//DELETE /:profession/:name
	//DELETE /:profession/:id
	delete(w http.ResponseWriter,r *http.Request,params httprouter.Params)

	genUrl(url string)string
}

func NewOpenAdmin(prefix string, admin eosc.IAdmin) *OpenAdmin {
	p := strings.TrimSpace(prefix)
	if len(p) == 0{
		p = "/"
	}else if p[0]!= '/'{
		p = "/"+p
	}
	p =  strings.TrimSuffix(p,"/")

	return &OpenAdmin{
		prefix: p,
		admin:  admin,
	}
}

func (o *OpenAdmin) GenHandler() (http.Handler, error) {
	var admin iOpenAdmin = o
	router := httprouter.New()
	router.GET(admin.genUrl("/"),admin.getProfessions)
 	router.GET(admin.genUrl("/:profession"),admin.getEmployeesByProfession)
	router.GET(admin.genUrl("/:profession/:action"), func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		action:=params.ByName("action")
		switch action {
		case "_render":
			admin.getRenders(w,r,params)
			return
		case "_driver":
			admin.getDriversByProfession(w,r,params)
		default:
			rename(params,"action","name")
			admin.getEmployeeByName(w,r,params)
		}
	})
	router.GET(admin.genUrl("/:profession/:action/:key"), func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		action:=params.ByName("action")
		switch action {
		case "_render":
			rename(params,"key","driver")

			admin.getRender(w,r, params)

			return
		case "_driver":
			key:= params.ByName("key")
			if strings.ToLower(key) == "item"{
				admin.getDriversItemByProfession(w,r,params)
			}else{
				rename(params,"key","driver")
				admin.getDriverInfo(w,r,params)
			}
			return
		default:

			rename(params,"action","name")
			rename(params,"key","fieldName")
			admin.getFields(w,r, params)
		}

	})

	router.POST(admin.genUrl("/:profession"),admin.Save)
	router.POST(admin.genUrl("/:profession/:name"),admin.Save)
	router.PUT(admin.genUrl("/:profession/:name"),admin.Save)
	router.DELETE(admin.genUrl("/:profession/:name"),admin.delete)
	return router,nil
}

func rename(ps httprouter.Params, source, target string)  {
	for i := range ps {
		if ps[i].Key == source {
			  ps[i].Key = target
		}
	}
}