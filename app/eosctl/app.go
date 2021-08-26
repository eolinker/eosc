package main

import (
	"github.com/urfave/cli/v2"
)

func newApp() *cli.App {
	return &cli.App{
		Name:  "eosctl",
		Usage: "a eosc controller process",
	}
}
