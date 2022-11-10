package traffic

import (
	"net"
	"testing"
	"time"
)

func Test_listener(t *testing.T) {
	tcp, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   nil,
		Port: 9988,
	})
	if err != nil {
		t.Error(err)
		return
	}
	tcp.SetDeadline(time.Now().Add(time.Second * 10))
	acceptTCP, err := tcp.AcceptTCP()
	if err != nil {
		t.Log(err)
		return
	}
	acceptTCP.Close()
}
