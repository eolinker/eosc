package cmuxMatch

import (
	"net"
	"sync"

	"github.com/soheilhy/cmux"
)

var (
	ErrorListenerClosed = cmux.ErrListenerClosed
)

type ListenerProxy struct {
	ch       chan net.Conn
	root     net.Listener
	once     sync.Once
	done     chan struct{}
	shutDown chan struct{}
}

func NewListenerProxy(l net.Listener, shutDown chan struct{}) *ListenerProxy {

	p := &ListenerProxy{root: l, ch: make(chan net.Conn, 10), done: make(chan struct{}), shutDown: shutDown}
	go p.doAccept()
	return p
}
func (l *ListenerProxy) Replace() *ListenerProxy {

	n := &ListenerProxy{root: l.root, ch: make(chan net.Conn, 10), done: make(chan struct{}), shutDown: l.shutDown}

	go func() {
		l.Close()
		for c := range l.ch {
			n.ch <- c
		}

		n.doAccept()
	}()
	return n
}
func (l *ListenerProxy) doAccept() {
	defer close(l.ch)
	for {
		c, err := l.root.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok {
				if ne.Timeout() {
					continue
				}
			}
			l.root = nil
			l.ShutDown()
			return
		}
		l.ch <- c
		select {
		case <-l.done:

			return
		default:
			continue
		}
	}
}
func (l *ListenerProxy) Accept() (net.Conn, error) {
	select {
	case conn, ok := <-l.ch:
		if ok {
			return conn, nil
		}
		return nil, ErrorListenerClosed
	case <-l.done:
		return nil, ErrorListenerClosed
	}

}

func (l *ListenerProxy) Close() error {
	l.once.Do(func() {
		close(l.done)
	})
	return nil
}
func (l *ListenerProxy) ShutDown() {
	l.once.Do(func() {
		if l.root != nil {
			l.root.Close()
			l.root = nil
		}

		close(l.shutDown)
		close(l.done)
	})

}

func (l *ListenerProxy) Addr() net.Addr {
	return l.root.Addr()
}
