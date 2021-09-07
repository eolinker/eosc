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

func CreateLevel(pre, affter cli.ActionFunc) *cli.Command {
	return Leave(func(c *cli.Context) error {
		err := pre(c)
		if err != nil {
			return err
		}
		// todo cli level

		return affter(c)

	})
}
