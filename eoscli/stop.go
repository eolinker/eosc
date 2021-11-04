package eoscli

import (
	"fmt"

	"github.com/eolinker/eosc/env"
	"github.com/urfave/cli/v2"
)

var CmdStop = "stop"

func Stop() *cli.Command {
	return &cli.Command{
		Name:   "stop",
		Usage:  fmt.Sprintf("stop %s server", env.AppName()),
		Action: StopFunc,
	}
}

//StopFunc 停止节点
func StopFunc(c *cli.Context) error {
	pidDir := env.PidFileDir()
	// 判断程序是否存在
	if !CheckPIDFILEAlreadyExists(pidDir) {
		ClearPid(pidDir)
		return nil
	}
	return stopProcess()
}
