package unix_proxy

import (
	"bufio"
	"net"
)

// The Conn type represents a WebSocket connection.
type Conn struct {
	net.Conn
	readBuf *bufio.Reader
}

func (c *Conn) Read(p []byte) (int, error) {
	return c.readBuf.Read(p)
}
