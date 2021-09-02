/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package master

import (
	"bytes"
	"github.com/eolinker/eosc/master/service"
	"github.com/eolinker/eosc/process"
	"github.com/eolinker/eosc/traffic"
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
	log.SetPrefix("[master]")

	log.Println("start master")

	srv, err := service.StartMaster("/tmp/eoserver.master.sock")
 	if err!= nil{
 		log.Println(err)
 		os.Exit(1)
		return
	}


	masterSrv = srv
	defer 	masterSrv.Stop()

	trafficController := traffic.NewController()
	defer trafficController.Close()
	err = trafficController.Listener("tcp", ":1900")
	if err != nil {
		return
	}
	err = trafficController.Listener("tcp", ":1901")
	if err != nil {
		return
	}

	worker,err :=process.Cmd("worker",nil)
	if err != nil{
		log.Println(err)
		return
	}
	traffics:=trafficController.All()

 	buf:=&bytes.Buffer{}

	files,err := traffics.WriteTo(buf)
	if err != nil {
		return
	}
 	worker.Stdin = nil
	pipe, err := worker.StdinPipe()
	if err != nil {
		log.Println(err)
		return
	}
	worker.Stdout = os.Stdout
	worker.Stderr = os.Stderr
	worker.ExtraFiles = files
	log.Println("file:",len(files))
	err = worker.Start()
	if err != nil {
		log.Println(err)
		return
	}

	bsize, err := buf.WriteTo(pipe)
	if err != nil {
		return
	}
	log.Println("write to worker:",bsize)


	wait()

}

func wait()  {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Wait for a SIGINT or SIGKILL:
	sig := <- sigc
	log.Printf("Caught signal %s: shutting down.", sig)
	// Stop listening (and unlink the socket if unix type):

	syscall.Unlink("/tmp/eoserver.master.sock")

	//os.Exit(0)
}