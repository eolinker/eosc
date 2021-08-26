package cli

import "github.com/urfave/cli/v2"

func Start() *cli.Command {
	return &cli.Command{
		Name: "start",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "admin_ip",
				Usage: "",
				Value: "0.0.0.0",
			},
			&cli.IntFlag{
				Name:  "admin_port",
				Usage: "",
				Value: 9400,
			},
		},
		Action: join,
	}
}
