package eoscli

import (
	"fmt"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/process"
	"github.com/urfave/cli/v2"
)

func Master() *cli.Command {
	return &cli.Command{
		Name:  "debug",
		Usage: "run as master",
		Subcommands: []*cli.Command{
			{
				Name:  "master",
				Usage: "debug as master",
				Action: func(context *cli.Context) error {
					log.Info("run master")
					if process.RunDebug(eosc.ProcessMaster) {

						log.Info("debug master done")
						return nil
					} else {
						return fmt.Errorf("debug master done")
					}
				},
			},
			{
				Name:  "admin",
				Usage: "debug as admin",
				Action: func(context *cli.Context) error {
					log.Info("run admin")
					if process.RunDebug(eosc.ProcessAdmin) {
						log.Info("debug admin done")
						return nil
					} else {
						return fmt.Errorf("debug admin done")
					}
				},
			},
			{
				Name:  "worker",
				Usage: "debug as worker",
				Action: func(context *cli.Context) error {
					log.Info("run worker")
					if process.RunDebug(eosc.ProcessWorker) {

						log.Info("debug worker done")
						return nil
					} else {
						return fmt.Errorf("debug worker done")
					}
				},
			},
		},
	}
}
