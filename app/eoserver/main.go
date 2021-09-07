//+build !windows

/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package main

import (
	"github.com/eolinker/eosc/eoscli"
	"github.com/eolinker/eosc/helper"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/master"
	"github.com/eolinker/eosc/process"
	"github.com/eolinker/eosc/worker"
	"github.com/urfave/cli/v2"
	"os"
)

func init() {

	process.Register("worker", worker.Process)
	process.Register("master", master.Process)
	process.Register("helper", helper.Process)
}

func main() {
	defer os.Exit(0)
	if process.Run() {
		return
	}
	log.InitDebug(true)
	app := eoscli.NewApp()
	app.AppendCommand(
		eoscli.Start(start),
		eoscli.Join(join),
		eoscli.Stop(stop),
		eoscli.Info(info),
		eoscli.Leave(leave),
		eoscli.Cluster(clusters),
		eoscli.Restart(func(c *cli.Context) error {
			return process.Restart()
		}),
	)
	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
	}


}
