package eoscli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func Start(start cli.ActionFunc) *cli.Command {
	return &cli.Command{
		Name:  "start",
		Usage: "start eosc server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "admin-ip",
				Aliases: []string{"ip"},
				Usage:   "ip for the admin process",
				Value:   "0.0.0.0",
			},
			&cli.IntFlag{
				Name:    "admin-port",
				Aliases: []string{"port", "p"},
				Usage:   "port for the admin process",
				Value:   9400,
			},
			&cli.BoolFlag{
				Name:  "join",
				Usage: "join the cluster if true",
			},
			&cli.StringFlag{
				Name:  "broadcast-ip",
				Usage: "ip for the node broadcast, required when join is true",
			},
			&cli.IntFlag{
				Name:  "broadcast-port",
				Usage: "port for the node broadcast, required when join is true",
				Value: 9401,
			},
			&cli.IntFlag{
				Name:    "cluster-addr",
				Aliases: []string{"addr"},
				Usage:   "port for the node broadcast",
			},
		},
		Action: start,
	}
}

func start(c *cli.Context) error {
	fmt.Println("eosc start")
	return nil
}
