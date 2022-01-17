package eoscli

import (
	"fmt"

	"github.com/eolinker/eosc/env"
	"github.com/urfave/cli/v2"
)

func Restart() *cli.Command {
	return &cli.Command{
		Name:  "restart",
		Usage: fmt.Sprintf("restart %s server", env.AppName()),

		Action: RestartFunc,
	}
}

func RestartFunc(c *cli.Context) error {
	// 先check插件版本是否存在，若不存在，则先下载插件后才能执行restart

	return restartProcess()
}
