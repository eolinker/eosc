package eoscli

import "github.com/urfave/cli/v2"

func Info(info cli.ActionFunc) *cli.Command {
	return &cli.Command{
		Name:   "info",
		Usage:  "display information of the node",
		Action: info,
	}
}
