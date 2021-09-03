package eoscli

import (
	"github.com/urfave/cli/v2"
)

func Leave(leave cli.ActionFunc) *cli.Command {
	return &cli.Command{
		Name:  "leave",
		Usage: "leave the cluster",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "id",
				Usage:    "id of node",
				Required: true,
			},
		},
		Action: leave,
	}
}
