/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package worker

import (
	"log"
	"net/http"
	"time"
)

func Work() {
	log.Println("start work")






	//
	//conn, err := grpc_unixsocket.Connect("/tmp/eoserver.master.sock")
	//if err!= nil{
	//	log.Println(err)
	//	return
	//}
	//defer conn.Close()
	//masterClient:=service.NewMasterClient(conn)
	//listener, err := process_listener.NewListener(masterClient, 1900)
	//if err!= nil{
	//	log.Println(err)
	//	return
	//}
	//err=http.Serve(listener,new(httpTest))
	//if err!= nil{
	//	log.Println(err)
	//}
}
type httpTest struct {

}

func (h *httpTest) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	n:=time.Now()
	w.Write([]byte(n.Format(time.RFC3339Nano)))

}

