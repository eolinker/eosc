package traffic_http_fast

import (
	"sync"

	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/traffic"
)

var _ IHttpTraffic = (*HttpTraffic)(nil)

type IHttpTraffic interface {
	Set(port int, srv *HttpService)
	Get(port int) (IService, bool)
	All() map[int]IService
	ShutDown(port int)
	Close()
}

type HttpTraffic struct {
	locker sync.Mutex
	tf     traffic.ITraffic
	srvs   map[int]*HttpService
}

func (h *HttpTraffic) ShutDown(port int) {

	h.locker.Lock()
	defer h.locker.Unlock()

	if s, has := h.srvs[port]; has {
		log.Debug("http traffic shutdown,port is ", port)
		s.ShutDown()

		delete(h.srvs, port)
	}
	return
}

func (h *HttpTraffic) Close() {
	h.locker.Lock()
	defer h.locker.Unlock()

	for _, s := range h.srvs {
		s.ShutDown()
	}

	return
}

func (h *HttpTraffic) Set(port int, srv *HttpService) {
	h.locker.Lock()
	defer h.locker.Unlock()
	h.srvs[port] = srv
}

func (h *HttpTraffic) Get(port int) (IService, bool) {

	h.locker.Lock()
	defer h.locker.Unlock()

	srv, has := h.srvs[port]
	if has {
		return srv, true
	}
	return nil, false
	//log.Debug("http traffic get:", port)
	//listener, err := h.tf.ListenTcp("", port)
	//
	//if err != nil {
	//	srv = NewHttpService(nil)
	//} else {
	//	srv = NewHttpService(listener)
	//}
	//h.srvs[port] = srv
	//return srv
}

func (h *HttpTraffic) All() map[int]IService {
	h.locker.Lock()
	defer h.locker.Unlock()
	srv := make(map[int]IService)
	for k, v := range h.srvs {
		srv[k] = v
	}
	return srv
}

func NewHttpTraffic(tf traffic.ITraffic) *HttpTraffic {
	return &HttpTraffic{
		locker: sync.Mutex{},
		tf:     tf,
		srvs:   make(map[int]*HttpService),
	}
}
