package proxy

import (
	"github.com/eolinker/eosc/log"
	"io"
	"net"
)

func doProxy(from, to net.Conn) {
	defer func() {
		from.Close()
		to.Close()
	}()

	go func() {

		_, err := io.Copy(from, to)
		if err != nil {
			log.DebugF("copy to client:%s", err.Error())
			return
		}
	}()
	_, err := io.Copy(to, from)
	if err != nil {
		log.DebugF("copy to unix:%s", err.Error())
		return
	}
}
