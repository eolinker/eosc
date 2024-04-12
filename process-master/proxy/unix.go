package proxy

import (
	"fmt"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/service"
	"net"
	"os/exec"
	"time"
)

type UnixProxy struct {
	addr    string
	name    string
	timeout time.Duration
}

func NewUnixProxy(name string) *UnixProxy {
	return &UnixProxy{name: name}
}
func (uc *UnixProxy) Update(process *exec.Cmd) {
	log.DebugF("unix client update: %s %s", uc.name, process)
	if process == nil {
		uc.addr = ""
		return
	}
	uc.addr = service.ServerAddr(process.Process.Pid, uc.name)
}
func (uc *UnixProxy) dialToProcess() (net.Conn, error) {
	if uc.addr == "" {
		return nil, fmt.Errorf("%s rocess not init", uc.name)

	}
	return net.DialTimeout("unix", uc.addr, uc.timeout)

}
func (uc *UnixProxy) ProxyToUnix(conn net.Conn) {
	targetConn, err := uc.dialToProcess()
	if err != nil {
		conn.Close()
		log.DebugF("dial to process:%s", err.Error())
		return
	}
	doProxy(conn, targetConn)

}
