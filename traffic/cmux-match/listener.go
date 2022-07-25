package cmuxMatch

import (
	"fmt"
	"net"
	"sync"
)

type shutListener struct {
	lock sync.RWMutex
	net.Listener
}

func (m *shutListener) Accept() (net.Conn, error) {

	m.lock.RLock()
	defer m.lock.RUnlock()
	if m.Listener == nil {
		return nil, net.ErrClosed
	}

	accept, err := m.Listener.Accept()
	if err != nil {
		return nil, err
	}
	fmt.Println("accept:", accept.RemoteAddr().String())
	return accept, nil
}

func (m *shutListener) Close() error {
	return nil
}
func (m *shutListener) Shutdown() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	err := m.Listener.Close()
	m.Listener = nil
	return err
}
func (m *shutListener) reset(l net.Listener) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.Listener = l
}
func newListener() *shutListener {
	return &shutListener{}
}
