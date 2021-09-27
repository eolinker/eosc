package traffic

import (
	"net"
)

type tTcpListener struct {
	net.Listener
}

func (t *tTcpListener) Close() error {
	panic("implement me")
}
