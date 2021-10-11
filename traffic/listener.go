package traffic

import (
	"net"
	"sync"

	"github.com/eolinker/eosc/log"
)

type iRemove interface {
	remove(name string)
}

type tListener struct {
	once sync.Once
	net.Listener
	parent iRemove
}

func newTTcpListener(listener net.Listener, parent iRemove) *tListener {

	return &tListener{Listener: listener, parent: parent}
}

func (t *tListener) Close() error {
	log.Debug("tListener close try")
	t.once.Do(func() {

		name := toName(t.Listener)
		log.Info("shutdown listener:", name)
		t.parent.remove(name)
		t.Listener.Close()

	})
	log.Debug("tListener close done")
	return nil
}
