package mixl

import (
	"net"
	"sync"
	"sync/atomic"

	"github.com/eolinker/eosc/log"
)

type MixListener struct {
	lock       sync.RWMutex
	listeners  map[string]net.Listener
	acceptChan chan net.Conn
	addr       net.Addr
	closeCh    chan struct{}
	sitClose   int32
	wg         sync.WaitGroup
}

func NewMixListener(port int, listeners ...net.Listener) *MixListener {
	lm := make(map[string]net.Listener)

	for _, listener := range listeners {
		name := listener.Addr().String()
		lm[name] = listener
	}
	ml := &MixListener{
		lock:       sync.RWMutex{},
		listeners:  lm,
		acceptChan: make(chan net.Conn, 10),
		addr: &net.TCPAddr{
			Port: port,
		},
		closeCh: make(chan struct{}),
		//sitClose: 0,
		wg: sync.WaitGroup{},
	}
	atomic.StoreInt32(&ml.sitClose, 0)

	for _, l := range listeners {
		ml.wg.Add(1)
		go ml.accept(l)
	}
	return ml
}

//	func (m *MixListener) Add(ls ...*net.TCPListener) {
//		isClosed := atomic.LoadInt32(&m.sitClose)
//		if isClosed == 0 {
//			m.lock.Lock()
//
//			for _, l := range ls {
//				name := l.Addr().String()
//				if _, has := m.listeners[name]; !has {
//					m.listeners[name] = l
//					m.wg.Add(1)
//					go m.accept(l)
//				}
//
//			}
//
//			m.lock.Unlock()
//		}
//	}
//
//	func (m *MixListener) Listeners() []*net.TCPListener {
//		m.lock.RLock()
//		defer m.lock.RUnlock()
//		ls := make([]*net.TCPListener, 0, len(m.listeners))
//		for _, l := range m.listeners {
//			ls = append(ls, l)
//		}
//		return ls
//
// }
func (m *MixListener) accept(l net.Listener) {
	defer func() {
		m.lock.Lock()
		defer m.lock.Unlock()
		name := l.Addr().String()
		delete(m.listeners, name)
		m.wg.Done()
	}()
	for {
		conn, err := l.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok {
				if ne.Timeout() {
					continue
				}
			}
			log.Info("listener accept:", l.Addr().Network(), ":", err)
			return
		}
		isClose := atomic.LoadInt32(&m.sitClose)
		if isClose == 0 {
			m.acceptChan <- conn
		} else {
			conn.Close()
			l.Close()
			return
		}
	}
}
func (m *MixListener) Accept() (net.Conn, error) {
	select {
	case conn, ok := <-m.acceptChan:
		{
			if ok {
				return conn, nil
			}
		}
	case <-m.closeCh:
	}

	return nil, net.ErrClosed

}

func (m *MixListener) Close() error {
	isClosed := atomic.SwapInt32(&m.sitClose, 1)
	if isClosed == 0 {

		m.lock.Lock()
		listeners := m.listeners
		m.lock.Unlock()
		for _, l := range listeners {
			l.Close()
		}
		m.wg.Wait()
	skip:
		for {
			select {
			case conn, ok := <-m.acceptChan:
				if ok {
					conn.Close()
				}
			default:
				close(m.acceptChan)
				break skip
			}
		}
		close(m.closeCh)
	}

	return nil
}

func (m *MixListener) Addr() net.Addr {

	return m.addr
}
