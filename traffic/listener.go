package traffic

import (
	"errors"
	"github.com/eolinker/eosc/log"
	"net"
	"os"
)

var (
	ErrorNotTcpListener = errors.New("not tcp port-reqiure")
)

type tListener struct {
	*net.TCPListener

	file      *os.File
	fileError error
}

func newTTcpListener(listener *net.TCPListener) *tListener {
	return &tListener{TCPListener: listener}
}
func (t *tListener) shutdown() {

	if t.file != nil {
		t.file.Close()
	}
	t.Close()

	log.Debug("tListener close done")
}
