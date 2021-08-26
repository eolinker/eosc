package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func Join() *cli.Command {
	return &cli.Command{
		Name: "join",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "broadcast_ip",
				Usage:    "",
				Required: true,
				Value:    "127.0.0.1",
			},
			&cli.IntFlag{
				Name:  "broadcast_port",
				Usage: "",
				Value: 9401,
			},
		},
		Action: join,
	}
}

func join(c *cli.Context) error {
	fmt.Println("the")
	// TODO 集群加入操作
	return nil
}
