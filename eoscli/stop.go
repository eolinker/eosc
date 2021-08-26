package eoscli

import (
	"os"

	"github.com/urfave/cli/v2"
)

func Stop() *cli.Command {
	return &cli.Command{
		Name:        "stop",
		Usage:       "stop eosc server",
		Action:      stop,
		Subcommands: []*cli.Command{},
	}
}

func stop(c *cli.Context) error {
	os.Exit(0)
	return nil
}
