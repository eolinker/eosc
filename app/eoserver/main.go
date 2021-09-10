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
	"os"

	"github.com/eolinker/eosc/eoscli"
	"github.com/eolinker/eosc/helper"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/master"
	"github.com/eolinker/eosc/process"
	"github.com/eolinker/eosc/worker"
)

func init() {

	process.Register("worker", worker.Process)
	process.Register("master", master.Process)
	process.Register("helper", helper.Process)
}

func main() {

	if process.Run() {
		log.Close()
		return
	}
	app := eoscli.NewApp()
	app.AppendCommand(
		eoscli.Start(eoscli.StartFunc),
		eoscli.Join(eoscli.JoinFunc),
		eoscli.Stop(eoscli.StopFunc),
		eoscli.Info(eoscli.InfoFunc),
		eoscli.Leave(eoscli.LeaveFunc),
		eoscli.Cluster(eoscli.ClustersFunc),
		eoscli.Restart(eoscli.RestartFunc),
		eoscli.Env(eoscli.EnvFunc),
	)
	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
	}
	log.Close()
}
