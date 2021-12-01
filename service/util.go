package service

import (
	"github.com/eolinker/eosc/env"
)

func MasterServerAddr(pid int) string {

	return env.SocketAddr("master", pid)
}

func WorkerServerAddr(pid int) string {

	return env.SocketAddr("worker", pid)
}
