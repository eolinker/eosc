package cmuxMatch

import (
	"net"
	"testing"
)

func TestNewMatch(t *testing.T) {
	listener, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		return
	}
	errChan := make(chan error)
	match := NewMatch(listener)
	h1l := match.Match(Http1)

	go runTestHTTPServer(errChan, h1l)

	rpcl := match.Match(Any)

	go runTestRPCServer(errChan, rpcl)

	runTestHTTP1Client(t, "127.0.0.1:9090")
	<-errChan

	// output
	//
}
