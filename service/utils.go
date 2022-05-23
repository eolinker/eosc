package service

import (
	"fmt"

	"github.com/eolinker/eosc/env"
)

func ServerAddr(pid int, name string) string {

	return env.SocketAddr(name, pid)
}

func ServerUnixAddr(pid int, name string) string {
	return env.SocketAddr(fmt.Sprintf("unix-%s", name), pid)
}
