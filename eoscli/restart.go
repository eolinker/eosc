package eoscli

import "github.com/urfave/cli/v2"

func Restart() *cli.Command {
	return &cli.Command{
		Name:  "restart",
		Usage: "restart goku server",

		Action: RestartFunc,
	}
}

func RestartFunc(c *cli.Context) error {
	return restartProcess()
}
