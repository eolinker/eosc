package eoscli

import (
	"os"
	"syscall"

	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/process"
)

func restartProcess() error {
	pid, err := readPid()

	if err != nil {
		return err
	}

	log.Debugf("app %s pid:%d is restart,please wait...\n", process.AppName(), pid)

	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Signal(syscall.SIGUSR1)
}

func stopProcess() error {
	log.Debugf("app %s is stopping,please wait...\n", process.AppName())
	pid, err := readPid()
	if err != nil {
		return err
	}
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Signal(os.Interrupt)
}
