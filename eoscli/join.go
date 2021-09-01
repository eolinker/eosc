package eoscli

import (
	"github.com/urfave/cli/v2"
)

func Join(join cli.ActionFunc) *cli.Command {
	return &cli.Command{
		Name:  "join",
		Usage: "join the cluster",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "broadcast-ip",
				Aliases:  []string{"ip"},
				Usage:    "ip for the node broadcast",
				Required: true,
			},
			&cli.IntFlag{
				Name:    "broadcast-port",
				Aliases: []string{"p", "port"},
				Usage:   "port for the node broadcast",
				Value:   9401,
			},
			&cli.StringSliceFlag{
				Name:    "cluster-addr",
				Aliases: []string{"addr"},
				Usage:   "<scheme>://<ip>:<port> of any node that already in cluster",
			},
		},
		Action: join,
	}
}
