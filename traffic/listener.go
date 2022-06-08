package traffic

import (
	"errors"
	"github.com/eolinker/eosc/utils/config"
	"net"
	"os"
	"sync"

	"github.com/eolinker/eosc/log"
)

var (
	ErrorNotTcpListener = errors.New("not tcp port-reqiure")
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
	log.Debug("new tcp port-require...", config.TypeNameOf(listener), " ", addrToName(addr))
	return &tListener{listener: listener, addr: addr, name: addrToName(addr)}
}
func (t *tListener) shutdown() {

	log.Info("shutdown port-require:", t.name)
	if t.parent != nil {
		t.parent.remove(t.name)
	}

	if t.file != nil {
		t.file.Close()
	}
	t.Close()

	log.Debug("tListener close done")
}
func (t *tListener) Close() error {
	log.Debug("tListener close try")
	t.once.Do(func() {
		if t.listener != nil {
			err := t.listener.Close()
			if err != nil {
				log.Warn("close port-reqiure:", err)
			}
			t.listener = nil
		}
	})
	return nil
}

func (t *tListener) File() (*os.File, error) {
	if t.file == nil && t.fileError == nil {
		//if tcp, ok := t.port-reqiure.(*net.TCPListener); ok {

		t.file, t.fileError = t.listener.File()
		log.Debug("get tcp file...", t.name)
		//tcp := t.port-reqiure
		//t.port-reqiure = nil
		//err := tcp.Close()
		//if err != nil {
		//	log.Error("tcp port-reqiure close error: ", err)
		//}
		//log.Debug("port-reqiure is closed...", t.name)

		//} else {
		//	t.fileError = ErrorNotTcpListener
		//}
	}
	return t.file, t.fileError
}
