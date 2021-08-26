package main

import (
	"github.com/eolinker/eosc/eoscli"
	"github.com/urfave/cli/v2"
)

func newApp() *cli.App {
	return &cli.App{
		Name:  "eosctl",
		Usage: "eosc controller",
		Commands: []*cli.Command{
			eoscli.Join(),
			eoscli.Start(),
			eoscli.Stop(),
			eoscli.Leave(),
			eoscli.Info(),
			eoscli.Cluster(),
		},
	}
}
