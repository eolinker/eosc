package eoscli

import (
	"os"
	"syscall"

	eosc_args "github.com/eolinker/eosc/eosc-args"
	"github.com/eolinker/eosc/log"
)

func restartProcess() error {
	pid, err := readPid()

	if err != nil {
		return err
	}

	log.Debugf("app %s pid:%d is restart,please wait...\n", eosc_args.AppName(), pid)

	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Signal(syscall.SIGUSR1)
}

func stopProcess() error {
	log.Debugf("app %s is stopping,please wait...\n", eosc_args.AppName())
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
