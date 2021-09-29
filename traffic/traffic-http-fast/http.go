package traffic_http_fast

import (
	"crypto/tls"
	"errors"
	"net"
	"strings"
	"sync"

	"github.com/eolinker/eosc/log"

	"github.com/valyala/fasthttp"
)

var _ IService = (*HttpService)(nil)
var (
	errorCertificateNotExit = errors.New("not exist cert")
)

type IService interface {
	SetHttps(handler fasthttp.RequestHandler, certs map[string]*tls.Certificate)
	SetHttp(handler fasthttp.RequestHandler)
	ShutDown()
}

const (
	serviceNotInit = 0
	serviceHttp    = 1
	serviceHttps   = 2
)

type HttpService struct {
	locker sync.Mutex
	certs  *Certs
	status int
	last   net.Listener
	inner  net.Listener
	srv    *fasthttp.Server
}

func (h *HttpService) SetHttps(handler fasthttp.RequestHandler, certs map[string]*tls.Certificate) {
	log.Debug("http Service SetHttps")

	h.locker.Lock()
	defer h.locker.Unlock()

	h.certs = newCerts(certs)
	h.srv.Handler = handler

	if h.inner != nil && h.status != serviceHttps {
		// http to https
		h.status = serviceHttps

		if h.last != nil {
			h.last.Close()
		}

		h.last = tls.NewListener(newNotClose(h.inner), &tls.Config{GetCertificate: h.GetCertificate})
		go h.srv.Serve(h.last)
	}
}

func (h *HttpService) SetHttp(handler fasthttp.RequestHandler) {
	log.Debug("http Service SetHttp")
	h.locker.Lock()
	defer h.locker.Unlock()
	h.srv.Handler = handler
	if h.inner != nil && h.status != serviceHttp {
		h.status = serviceHttp
		if h.last != nil {
			h.last.Close()
		}
		h.certs = nil

		h.last = newNotClose(h.inner)
		go h.srv.Serve(h.last)
	}
}

//GetCertificate 获取证书配置
func (h *HttpService) GetCertificate(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
	if h.certs == nil {
		return nil, errorCertificateNotExit
	}
	certificate, has := h.certs.Get(strings.ToLower(info.ServerName))
	if !has {
		return nil, errorCertificateNotExit
	}

	return certificate, nil
}

func (h *HttpService) ShutDown() {
	h.locker.Lock()
	defer h.locker.Unlock()
	h.srv.Shutdown()
	if h.inner != nil {
		h.last.Close()
		h.last = nil
		h.inner.Close()
	}
}

func NewHttpService(listener net.Listener) *HttpService {
	return &HttpService{
		inner: listener,
		srv:   &fasthttp.Server{},
	}
}
