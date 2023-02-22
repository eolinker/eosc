package router

import (
	"github.com/eolinker/eosc/log"
	"net/http"
	"sync"
)

var _ IRouter = (*router)(nil)

var globalRouter *router

type IRouter interface {
	http.Handler
	Set(id string, path string, handler http.Handler) error
	Delete(id string)
}

type router struct {
	lock sync.RWMutex

	serverMux   *http.ServeMux
	routersData *routerData
}

// NewRouterManager 创建路由管理器
func newRouter() *router {
	return &router{
		serverMux:   &http.ServeMux{},
		routersData: new(routerData),
	}
}

func GetHandler() *router {
	return globalRouter
}

func AddPath(id string, path string, handler http.Handler) error {
	return globalRouter.Set(id, path, handler)
}

func DeletePath(id string) {
	globalRouter.Delete(id)
}

func (m *router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	m.serverMux.ServeHTTP(writer, request)
}

func (m *router) Set(id string, path string, handler http.Handler) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	routersData := m.routersData.Set(id, path, handler)
	serverMux, err := routersData.Parse()
	if err != nil {
		log.Error("parse router data error: ", err)
		return err
	}
	m.serverMux = serverMux
	m.routersData = routersData
	return nil
}

func (m *router) Delete(id string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	routersData := m.routersData.Delete(id)
	serverMux, err := routersData.Parse()
	if err != nil {
		log.Errorf("delete router:%s %s", id, err.Error())
		return
	}
	m.serverMux = serverMux
	m.routersData = routersData
	return
}
