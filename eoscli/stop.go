package eoscli

import (
	"github.com/urfave/cli/v2"
)

var CmdStop = "stop"

func Stop(stop cli.ActionFunc) *cli.Command {
	return &cli.Command{
		Name:   "stop",
		Usage:  "stop eosc server",
		Action: stop,
	}
}
