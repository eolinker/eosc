package eoscli

import (
	"github.com/urfave/cli/v2"
)

func Restart(restart cli.ActionFunc) *cli.Command {
	return &cli.Command{
		Name:  "restart",
		Usage: "restart goku server",

		Action: restart,
	}
}
