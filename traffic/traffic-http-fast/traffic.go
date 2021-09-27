package traffic_http_fast

import (
	"sync"

	"github.com/eolinker/eosc/traffic"
)

var _ IHttpTraffic = (*HttpTraffic)(nil)

type IHttpTraffic interface {
	Get(port int) IService
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

func (h *HttpTraffic) Get(port int) IService {

	h.locker.Lock()
	defer h.locker.Unlock()

	if s, has := h.srvs[port]; has {
		return s
	}
	listener, err := h.tf.ListenTcp("", port)
	if err != nil {
		return nil
	}
	srv := NewHttpService(listener)
	h.srvs[port] = srv
	return srv
}

func NewHttpTraffic(tf traffic.ITraffic) *HttpTraffic {
	return &HttpTraffic{
		locker: sync.Mutex{},
		tf:     tf,
		srvs:   make(map[int]*HttpService),
	}
}
