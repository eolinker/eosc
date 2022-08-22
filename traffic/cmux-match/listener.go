package cmuxMatch

import (
	"context"
	"fmt"
	"net"
	"sync"
)

type shutListener struct {
	lock   sync.RWMutex
	ch     chan net.Conn
	addr   net.Addr
	ctx    context.Context
	cancel context.CancelFunc

	last net.Listener
}

func (m *shutListener) Addr() net.Addr {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.addr
}

func (m *shutListener) Accept() (net.Conn, error) {

	select {
	case <-m.ctx.Done():
		fmt.Println("Accept done")
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
		fmt.Println("shut listener: start")
		accept, err := l.Accept()
		fmt.Println("shut listener: end")

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
	return nil
	//m.lock.Lock()
	//defer m.lock.Unlock()

	//return m.close()
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
		ctx:    ctx,
		cancel: cancelFunc,
		ch:     make(chan net.Conn, 1),
	}
}
