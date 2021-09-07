package eoscli

import (
	"github.com/urfave/cli/v2"
)

func Restart(x cli.ActionFunc) *cli.Command {
	return &cli.Command{
		Name:  "restart",
		Usage: "restart goku server",

		Action: x,
	}
}

func RestartFunc(c *cli.Context) error {
	return restartProcess()
}
