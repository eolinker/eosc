package listener

import (
	"errors"
	"fmt"
	"github.com/eolinker/eosc"
	"net"
	"sync"
)


var (
	locker = sync.Mutex{}
	data   = eosc.NewUntyped()
)

func ListenTCP(port int, sign string) (net.Listener, error) {
	key := fmt.Sprintf("tcp:%d", port)
	locker.Lock()
	defer locker.Unlock()
	o, has := data.Get(key)
	if !has {

		tcpl, err := net.ListenTCP("tcp4", &net.TCPAddr{
			IP:   net.IPv4zero,
			Port: port,
			Zone: "",
		})
		if err != nil {
			return nil, err
		}
		l := newTcpListener(tcpl, sign, port)
		data.Set(key, l)
		return l, nil
	} else {
		l := o.(*tcpListener)
		if l.sing == sign {
			return l, nil
		}
		return nil, errors.New("port is ")
	}
}

type tcpListener struct {
	net.Listener
	sing string
	port int
}

func (l *tcpListener) Close() error {
	key := fmt.Sprintf("tcp:%d", l.port)

	data.Del(key)
	return l.Listener.Close()
}
func newTcpListener(listener net.Listener, sing string, port int) net.Listener {
	return &tcpListener{Listener: listener, sing: sing, port: port}
}
