//go:build !windows
// +build !windows

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

	"github.com/eolinker/eosc/env"
	process_admin "github.com/eolinker/eosc/process-admin"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/eoscli"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/process"
	process_helper "github.com/eolinker/eosc/process-helper"
	process_master "github.com/eolinker/eosc/process-master"
	process_worker "github.com/eolinker/eosc/process-worker"
)

func init() {
	process.Register(eosc.ProcessWorker, process_worker.Process)
	process.Register(eosc.ProcessAdmin, process_admin.Process)
	process.Register(eosc.ProcessMaster, process_master.Process)
	process.Register(eosc.ProcessHelper, process_helper.Process)
}

func main() {

	if process.Run() {
		log.Close()
		return
	}
	if env.IsDebug() {
		if process.RunDebug(eosc.ProcessMaster) {
			log.Info("debug done")
		} else {
			log.Error("debug not run")
		}
		log.Close()
		return
	}
	//utils.InitStdTransport()

	app := eoscli.NewApp()
	app.Default()

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
	}
	log.Close()
}
