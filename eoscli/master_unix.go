// +build linux freebsd darwin

package eoscli

import (
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/eolinker/eosc/process"
)

func StartMaster(args []string, extra []*os.File) (*exec.Cmd, error) {

	cmd, err := process.Cmd("master", args)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	//cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	cmd.ExtraFiles = extra
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}
	e := cmd.Start()
	if e != nil {
		log.Println(e)
		return nil, e
	}

	return cmd, nil
}
