/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package worker

import (
	"context"
	"fmt"
	"github.com/eolinker/eosc/traffic"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type TestHandler string

func (t TestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(t))
	n:=time.Now()
	w.Write([]byte(n.Format(time.RFC3339Nano)))

	log.Println("handler:")
}

func Process() {
	log.SetPrefix("[worker]")
	log.Println("start work")

	trfs, err := traffic.Reader(os.Stdin,3)
	if err!= nil{
		//if e:=utils.WriteFrame(os.Stdout,[]byte(err.Error()));e!= nil{
		//	return
		//}
		log.Println("read:",err)
		return
	}
	if len(trfs) ==0{
		return
	}
	sers:= make( []*http.Server,0,len(trfs))

	for i,t:=range trfs{
		l:=t.Listener
		ser:=&http.Server{
			Handler: TestHandler(fmt.Sprintf("test:%d",i)),
		}
		sers = append(sers, ser)
		go ser.Serve(l)
	}
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Wait for a SIGINT or SIGKILL:
	sig := <- sigc
	log.Printf("Caught signal %s: shutting down.", sig)
	// Stop listening (and unlink the socket if unix type):
	ctx:=context.Background()
	for _,ser:=range sers{
		ser.Shutdown(ctx)
	}



}