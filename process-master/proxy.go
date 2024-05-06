package process_master

import (
	"github.com/eolinker/eosc/log"
	"net"
)

func doServer(ln net.Listener, handler func(conn net.Conn)) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Error(" server accept error:", err)
			return
		}
		go handler(conn)
	}
}
