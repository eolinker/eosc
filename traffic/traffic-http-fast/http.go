package traffic_http_fast

import (
	"errors"
	"net"
	"sync"

	"github.com/eolinker/eosc/log"

	"github.com/valyala/fasthttp"
)

var _ IService = (*HttpService)(nil)
var (
	errorCertificateNotExit = errors.New("not exist cert")
)

type IService interface {
	Set(handler fasthttp.RequestHandler)
	ShutDown()
}

type HttpService struct {
	locker sync.Mutex
	status int
	inner  net.Listener
	srv    *fasthttp.Server
}

func (h *HttpService) Set(handler fasthttp.RequestHandler) {
	h.locker.Lock()
	defer h.locker.Unlock()
	if handler == nil {
		h.srv.Handler = NotFound
	}
	h.srv.Handler = handler

}

func (h *HttpService) ShutDown() {
	h.locker.Lock()
	defer h.locker.Unlock()

	h.srv.Shutdown()
	h.srv = nil
	log.Debug("http service shutdown done")
}

func NewHttpService(listener net.Listener) *HttpService {
	s := &HttpService{
		srv: &fasthttp.Server{Handler: NotFound, DisablePreParseMultipartForm: true},
	}
	go s.srv.Serve(listener)
	log.Debug("new http service:", listener.Addr())
	return s
}

func NotFound(ctx *fasthttp.RequestCtx) {
	ctx.NotFound()
	return
}
