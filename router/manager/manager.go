package manager

import (
	"github.com/eolinker/eosc/log"
	"net/http"
	"sync"
)

var _ IRouterManager = (*RouterManager)(nil)

var globalRouterManager *RouterManager

type IRouterManager interface {
	http.Handler
	Set(id string, path string, handler http.Handler) error
	Delete(id string)
}

type RouterManager struct {
	lock sync.RWMutex

	serverMux   *http.ServeMux
	routersData IRouterData
}

// NewRouterManager 创建路由管理器
func NewRouterManager() *RouterManager {
	return &RouterManager{routersData: new(RouterData)}
}

func GetGlobalRouterManager() IRouterManager {
	return globalRouterManager
}

func GlobalRouterManagerAdd(id string, path string, handler http.Handler) error {
	return globalRouterManager.Set(id, path, handler)
}

func GlobalRouterManagerDel(id string) {
	globalRouterManager.Delete(id)
}

func (m *RouterManager) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	m.serverMux.ServeHTTP(writer, request)
}

func (m *RouterManager) Set(id string, path string, handler http.Handler) error {
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

func (m *RouterManager) Delete(id string) {
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
