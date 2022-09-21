package cmuxMatch

import (
	"context"
	"net"
	"sync"
)

type shutListener struct {
	lock sync.RWMutex
	ch   chan net.Conn

	closeTemp chan struct{}
	addr      net.Addr
	ctx       context.Context
	cancel    context.CancelFunc

	last net.Listener
}

func (m *shutListener) Addr() net.Addr {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.addr
}

func (m *shutListener) Accept() (net.Conn, error) {

	select {
	case <-m.closeTemp:
		return nil, ErrorListenerClosed
	case <-m.ctx.Done():
		return nil, ErrorListenerClosed
	case conn, ok := <-m.ch:
		if ok {
			return conn, nil
		}

		return nil, ErrorListenerClosed
	}
}
func (m *shutListener) doAccept(l net.Listener) {

	for {

		accept, err := l.Accept()

		if err != nil {
			if ne, ok := err.(net.Error); ok {
				if ne.Timeout() {
					continue
				}
			}

			break
		}
		m.ch <- accept
	}

}
func (m *shutListener) Close() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	return m.close()
}
func (m *shutListener) close() error {
	if m.cancel != nil {
		m.cancel()
		m.cancel = nil
	}
	return nil
}
func (m *shutListener) Shutdown() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	return m.close()
}
func (m *shutListener) reset(l net.Listener) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.addr = l.Addr()

	m.last = l

	go m.doAccept(l)
}
func newListener() *shutListener {
	ctx, cancelFunc := context.WithCancel(context.Background())

	return &shutListener{
		closeTemp: make(chan struct{}, 1),
		ctx:       ctx,
		cancel:    cancelFunc,
		ch:        make(chan net.Conn, 1),
	}
}
