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

	"github.com/eolinker/eosc/helper"
	"github.com/eolinker/eosc/master"
	"github.com/eolinker/eosc/process"
	"github.com/eolinker/eosc/worker"
)

func init() {
	process.Register("eoserver: worker", worker.Work)
	process.Register("eoserver: master", master.Master)
	process.Register("eoserver: helper", helper.Helper)
}
func main() {
	if process.Run() {
		return
	}

	process.Start("eoserver: master", os.Args[1:], nil)
}
