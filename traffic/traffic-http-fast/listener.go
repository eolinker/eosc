package traffic_http_fast

import (
	"net"
	"syscall"
)

type listenerNotClose struct {
	inner net.Listener

	addr net.Addr
}

func (l *listenerNotClose) Accept() (net.Conn, error) {
	if l.inner == nil {
		return nil, syscall.EINVAL
	}
	return l.inner.Accept()
}

func (l *listenerNotClose) Addr() net.Addr {

	return l.addr
}

func (l *listenerNotClose) Close() error {
	l.inner = nil
	return nil
}

func newNotClose(inner net.Listener) *listenerNotClose {
	return &listenerNotClose{inner: inner, addr: inner.Addr()}
}
