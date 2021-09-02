/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package service

import (
	grpc_unixsocket "github.com/eolinker/eosc/grpc-unixsocket"
	"github.com/eolinker/eosc/service"
	"google.golang.org/grpc"
	"log"
)

func StartMaster(addr string) (*grpc.Server,error){
	log.Println("start Master :",addr)
	l, err := grpc_unixsocket.Listener(addr)
	if err!= nil{
		return nil,err
	}
	grpcServer := grpc.NewServer()
	service.RegisterMasterServer(grpcServer, NewMasterServer())
	go func() {
		grpcServer.Serve(l)
	}()
	return grpcServer,nil
}
