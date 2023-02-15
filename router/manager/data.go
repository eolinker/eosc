package manager

import (
	"net/http"
	"strings"
)

var _ IRouterData = (*RouterData)(nil)

type IRouterData interface {
	Set(id string, path string, router http.Handler) IRouterData
	Delete(id string) IRouterData
	Parse() (*http.ServeMux, error)
}

type RouterData struct {
	data map[string]*Router
}

func (rs *RouterData) Parse() (*http.ServeMux, error) {
	//TODO 做一个panic recovery以防mux.Handle panic

	//TODO 去重
	mux := new(http.ServeMux)
	for _, v := range rs.data {
		//添加前缀
		path := RouterPrefix + strings.Trim(v.Path, "/")
		mux.Handle(path, v.Router)
	}
	return mux, nil
}

func (rs *RouterData) set(r *Router) *RouterData {
	rs.data[r.Id] = r
	return rs
}

func (rs *RouterData) Set(id string, path string, router http.Handler) IRouterData {
	r := &Router{
		Id:     id,
		Path:   path,
		Router: router,
	}

	if _, exist := rs.data[id]; exist {
		return rs.clone(0).set(r)
	}
	return rs.clone(1).set(r)
}

func (rs *RouterData) Delete(id string) IRouterData {
	return rs.clone(0).delete(id)
}

func (rs *RouterData) delete(id string) IRouterData {
	delete(rs.data, id)
	return rs
}

func (rs *RouterData) clone(delta int) *RouterData {
	if delta < 0 {
		delta = 0
	}
	if rs == nil || len(rs.data) == 0 {
		return &RouterData{data: make(map[string]*Router, 1)}
	}

	data := make(map[string]*Router, len(rs.data)+delta)
	for k, v := range rs.data {
		data[k] = v
	}
	return &RouterData{data: data}
}
