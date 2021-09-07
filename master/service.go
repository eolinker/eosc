/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package master

import (
	"fmt"
	"syscall"

	"github.com/eolinker/eosc/service"

	grpc_unixsocket "github.com/eolinker/eosc/grpc-unixsocket"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/process"
	"google.golang.org/grpc"
)

//startService 开启master
func (m *Master) startService() error {

	addr := fmt.Sprintf("/tmp/%s.master.sock", process.AppName())
	// 移除unix socket
	syscall.Unlink(addr)

	log.Info("start Master :", addr)
	l, err := grpc_unixsocket.Listener(addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	service.RegisterCtiServiceServer(grpcServer, m)
	service.RegisterMasterServer(grpcServer, m)
	go func() {
		grpcServer.Serve(l)
	}()

	m.masterSrv = grpcServer
	return nil
}

func (m *Master) stopService() {
	m.masterSrv.GracefulStop()
	addr := fmt.Sprintf("/tmp/%s.master.sock", process.AppName())
	syscall.Unlink(addr)
}
