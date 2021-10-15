// +build linux freebsd darwin

package eoscli

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/eolinker/eosc"

	env "github.com/eolinker/eosc/env"

	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/process"
)

func StartMaster(args []string, extra []*os.File) (*exec.Cmd, error) {

	cmd, err := process.Cmd(eosc.ProcessMaster, args)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if env.IsDebug() {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	//cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	cmd.ExtraFiles = extra
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}
	e := cmd.Start()
	if e != nil {
		log.Error(e)
		return nil, e
	}

	return cmd, nil
}
