/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package master

import (
	"github.com/eolinker/eosc/process"
	"google.golang.org/grpc"
	"log"
	"os"
	"os/signal"
	"syscall"
)
var(
	masterSrv *grpc.Server
)
func Master() {
	log.Println("start master")
	//srv, err := service.StartMaster("/tmp/eoserver.master.sock")
 	//if err!= nil{
 	//	log.Println(err)
	//	return
	//}
	//masterSrv = srv


	work,err:=process.Start("worker",nil,nil)
	if err != nil{
		log.Println(err)
		return
	}
	os.Pipe()
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Wait for a SIGINT or SIGKILL:
	sig := <- sigc
	log.Printf("Caught signal %s: shutting down.", sig)
	// Stop listening (and unlink the socket if unix type):
	masterSrv.Stop()
	// And we're done:
	os.Exit(0)

}