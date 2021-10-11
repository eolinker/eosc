package traffic

import (
	"net"

	"github.com/eolinker/eosc/log"
)

type tListener struct {
	net.Listener
	parent ITraffic
}

func newTTcpListener(listener net.Listener, parent ITraffic) *tListener {
	return &tListener{Listener: listener, parent: parent}
}

func (t *tListener) Close() error {
	name := toName(t.Listener)
	log.Info("shutdown listener:", name)
	t.parent.remove(name)
	return nil
}
