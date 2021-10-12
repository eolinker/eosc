package traffic

import (
	"errors"
	"net"
	"os"
	"sync"

	"github.com/eolinker/eosc/log"
)

var (
	ErrorNotTcpListener = errors.New("not tcp listener")
)

type iRemove interface {
	remove(name string)
}

type tListener struct {
	net.Listener
	once      sync.Once
	parent    iRemove
	file      *os.File
	fileError error
	addr      net.Addr
	name      string
}

func (t *tListener) Addr() net.Addr {
	return t.addr
}

func newTTcpListener(listener net.Listener, parent iRemove) *tListener {
	addr := listener.Addr()

	return &tListener{Listener: listener, parent: parent, addr: addr, name: addrToName(addr)}
}
func (t *tListener) Close() error {
	log.Debug("tListener close try")
	t.once.Do(func() {

		log.Info("shutdown listener:", t.name)
		t.parent.remove(t.name)

		if t.file != nil {
			t.file.Close()
		}
		if t.Listener != nil {
			err := t.Listener.Close()
			if err != nil {
				log.Warn("close listener:", err)
			}
		}
	})
	log.Debug("tListener close done")
	return nil
}

func (t *tListener) File() (*os.File, error) {
	if t.file == nil && t.fileError == nil {
		if tcp, ok := t.Listener.(*net.TCPListener); ok {

			t.file, t.fileError = tcp.File()
			t.Listener.Close()
			t.Listener = nil
		} else {
			t.fileError = ErrorNotTcpListener
		}
	}
	return t.file, t.fileError
}
