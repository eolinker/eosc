package eoscli

import (
	"os"
	"syscall"

	"github.com/eolinker/eosc/env"
	"github.com/eolinker/eosc/log"
)

func restartProcess() error {
	pidDir := env.PidFilePath()
	pid, err := readPid(pidDir)

	if err != nil {
		return err
	}

	log.Debugf("app %s pid:%d is restart,please wait...\n", env.AppName(), pid)

	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Signal(syscall.SIGUSR1)
}

func stopProcess() error {
	pidDir := env.PidFilePath()
	log.Debugf("app %s is stopping,please wait...\n", env.AppName())
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
