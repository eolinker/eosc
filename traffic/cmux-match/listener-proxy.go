package cmuxMatch

import (
	"errors"
	"net"
	"sync"
)

var (
	ErrorListenerClosed = errors.New("closed")
)

type ListenerProxy struct {
	ch   chan net.Conn
	root net.Listener
	once sync.Once
	done chan struct{}
}

func NewListenerProxy(l net.Listener) *ListenerProxy {

	p := &ListenerProxy{root: l, ch: make(chan net.Conn, 10), done: make(chan struct{})}
	go p.doAccept()
	return p
}
func (l *ListenerProxy) Replace() *ListenerProxy {

	n := &ListenerProxy{root: l.root, ch: make(chan net.Conn, 10), done: make(chan struct{})}

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
	l.root.Close()
}
func (l *ListenerProxy) Addr() net.Addr {
	return l.root.Addr()
}
