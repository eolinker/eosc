package traffic

import (
	"errors"
	"net"
	"os"
	"sync"

	"github.com/eolinker/eosc"

	"github.com/eolinker/eosc/log"
)

var (
	ErrorNotTcpListener = errors.New("not tcp listener")
)

type iRemove interface {
	remove(name string)
}

type tListener struct {
	listener  *net.TCPListener
	once      sync.Once
	parent    iRemove
	file      *os.File
	fileError error
	addr      net.Addr
	name      string
}

func (t *tListener) Accept() (net.Conn, error) {
	if t.listener != nil {
		return t.listener.Accept()
	}
	return nil, ErrorInvalidListener
}
func (t *tListener) Addr() net.Addr {
	return t.addr
}

func newTTcpListener(listener *net.TCPListener) *tListener {
	addr := listener.Addr()
	log.Debug("new tcp listener...", eosc.TypeNameOf(listener), " ", addrToName(addr))
	return &tListener{listener: listener, addr: addr, name: addrToName(addr)}
}
func (t *tListener) Close() error {
	log.Debug("tListener close try")
	t.once.Do(func() {

		log.Info("shutdown listener:", t.name)
		if t.parent != nil {
			t.parent.remove(t.name)
		}

		if t.file != nil {
			t.file.Close()
		}
		if t.listener != nil {
			err := t.listener.Close()
			if err != nil {
				log.Warn("close listener:", err)
			}
			t.listener = nil
		}
	})
	log.Debug("tListener close done")
	return nil
}

func (t *tListener) File() (*os.File, error) {
	if t.file == nil && t.fileError == nil {
		//if tcp, ok := t.listener.(*net.TCPListener); ok {

		t.file, t.fileError = t.listener.File()
		log.Debug("get tcp file...", t.name)
		err := t.listener.Close()
		if err != nil {
			log.Error("tcp listener close error: ", err)
		}
		log.Debug("listener is closed...", t.name)
		t.listener = nil
		//} else {
		//	t.fileError = ErrorNotTcpListener
		//}
	}
	return t.file, t.fileError
}
