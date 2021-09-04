package eoscli

import "github.com/urfave/cli/v2"

var CmdInfo = "info"

func Info(info cli.ActionFunc) *cli.Command {
	return &cli.Command{
		Name:   CmdInfo,
		Usage:  "display information of the node",
		Action: info,
	}
}
