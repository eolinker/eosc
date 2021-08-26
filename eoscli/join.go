package eoscli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func Join() *cli.Command {
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

func join(c *cli.Context) error {
	ip := c.String("broadcast_ip")
	port := c.Int("broadcast_port")
	fmt.Printf("broadcast ip is %s, broadcast port is %d\n", ip, port)
	fmt.Printf("eosc is joining, please wait...\n")

	// TODO 集群加入操作
	return nil
}
