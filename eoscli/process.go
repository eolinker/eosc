package eoscli

import (
	"os"
	"syscall"

	"github.com/eolinker/eosc/env"
	"github.com/eolinker/eosc/log"
)

func restartProcess() error {
	pidDir := env.PidFileDir()
	pid, err := readPid(pidDir)

	if err != nil {
		return err
	}

	log.DebugF("app %s pid:%d is restart,please wait...\n", env.AppName(), pid)

	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Signal(syscall.SIGUSR1)
}

func stopProcess() error {
	pidDir := env.PidFileDir()
	log.DebugF("app %s is stopping,please wait...\n", env.AppName())
	pid, err := readPid(pidDir)
	if err != nil {
		return err
	}
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Signal(os.Interrupt)
}
