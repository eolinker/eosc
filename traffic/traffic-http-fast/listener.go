package traffic_http_fast

import (
	"net"
	"syscall"

	"github.com/eolinker/eosc/log"
)

type listenerNotClose struct {
	inner net.Listener

	addr net.Addr
}

func (l *listenerNotClose) Accept() (net.Conn, error) {
	log.Debug("accept:start")
	if l.inner == nil {
		log.Debug("accept:nil")
		return nil, syscall.EINVAL
	}
	accept, err := l.inner.Accept()
	if err != nil {
		log.Debug("accept: error")

		return nil, err
	}
	log.Debug("accept: done")

	return accept, nil
}

func (l *listenerNotClose) Addr() net.Addr {

	return l.addr
}

func (l *listenerNotClose) Close() error {
	l.inner = nil
	return nil
}

func newNotClose(inner net.Listener) *listenerNotClose {
	log.Debug("new not close listener:", inner.Addr())
	return &listenerNotClose{inner: inner, addr: inner.Addr()}
}
