package eoscli

import (
	"github.com/urfave/cli/v2"
)

func Stop(stop cli.ActionFunc) *cli.Command {
	return &cli.Command{
		Name:        "stop",
		Usage:       "stop eosc server",
		Action:      stop,
		Subcommands: []*cli.Command{},
	}
}

//StopFunc 停止节点
func StopFunc(c *cli.Context) error {
	// 判断程序是否存在
	if !CheckPIDFILEAlreadyExists() {
		ClearPid()
		return nil
	}
	return stopProcess()
}
