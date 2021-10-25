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

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/eoscli"
	"github.com/eolinker/eosc/helper"
	"github.com/eolinker/eosc/log"
	admin_open_api "github.com/eolinker/eosc/modules/admin-open-api"
	"github.com/eolinker/eosc/process"
	process_master "github.com/eolinker/eosc/process-master"
	"github.com/eolinker/eosc/process-master/admin"
	process_worker "github.com/eolinker/eosc/process-worker"
)

func init() {
	admin.Register("/api/", admin_open_api.CreateHandler())
	process.Register(eosc.ProcessWorker, func() {
		process_worker.Process(nil)
	})
	process.Register(eosc.ProcessMaster, process_master.Process)
	process.Register(eosc.ProcessHelper, helper.Process)
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
	app := eoscli.NewApp()
	app.Default()
	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
	}
	log.Close()
}
