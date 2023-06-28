package router

import (
	"fmt"
	"net/http"
	"strings"
)

type routerData struct {
	data map[string]*routerConfig
}

func (rs *routerData) Parse() (*http.ServeMux, error) {
	tmpSet := make(map[string]struct{}, len(rs.data))
	mux := new(http.ServeMux)

	for _, v := range rs.data {
		//添加前缀
		path := RouterPrefix + strings.Trim(v.Path, "/")
		//查重
		if _, exist := tmpSet[path]; exist {
			return nil, fmt.Errorf("parse Router fail. path=%s", v.Path)
		}
		tmpSet[path] = struct{}{}

		if !strings.HasSuffix(v.Path, "/") {
			mux.Handle(path, v.Router)
		}

		mux.Handle(fmt.Sprint(path, "/"), v.Router)
	}
	return mux, nil
}

func (rs *routerData) set(r *routerConfig) *routerData {
	rs.data[r.Id] = r
	return rs
}

func (rs *routerData) Set(id string, path string, router http.Handler) *routerData {
	r := &routerConfig{
		Id:     id,
		Path:   path,
		Router: router,
	}

	delta := 1
	if _, exist := rs.data[id]; exist {
		delta = 0
	}
	return rs.clone(delta).set(r)
}

func (rs *routerData) Delete(id string) *routerData {
	return rs.clone(0).delete(id)
}

func (rs *routerData) delete(id string) *routerData {
	delete(rs.data, id)
	return rs
}

func (rs *routerData) clone(delta int) *routerData {
	if delta < 0 {
		delta = 0
	}
	if rs == nil || len(rs.data) == 0 {
		return &routerData{data: make(map[string]*routerConfig, 1)}
	}

	data := make(map[string]*routerConfig, len(rs.data)+delta)
	for k, v := range rs.data {
		data[k] = v
	}
	return &routerData{data: data}
}
