package listener

import (
	"net"

	"github.com/eolinker/eosc/traffic"
)

type Listener struct {
	traffic traffic.ITraffic
}

func (l *Listener) ListenTcp(ip string, port int) (net.Listener, error) {
	panic("implement me")
}

func (l *Listener) Close() {
	panic("implement me")
}
