/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package main

import (
	"github.com/eolinker/eosc/process"
	"log"
	"os"
	"time"
)

func init() {
	process.Register("worker:eoServer", func() {
		log.Print("run worker:eoServer")
		time.Sleep(time.Minute)
	})
	process.Register("master:eoServer", func() {
		log.Print("run master:eoServer")
		log.Print("call worker:eoServer")
		process.Start("worker:eoServer",nil)

		time.Sleep(time.Minute)


	})
}
func main() {
	if process.Run(){
		return
	}
	process.Start("master:eoServer",os.Args[1:])
}
