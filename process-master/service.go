/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package process_master

import (
	"os"
	"syscall"

	"github.com/eolinker/eosc"

	"github.com/eolinker/eosc/process-master/cli"

	"github.com/eolinker/eosc/service"

	grpc_unixsocket "github.com/eolinker/eosc/grpc-unixsocket"
	"github.com/eolinker/eosc/log"
	"google.golang.org/grpc"
)

// startMasterGRPCService 开启master
func (m *Master) startMasterGRPCService() error {

	addr := service.ServerAddr(os.Getpid(), eosc.ProcessMaster)
	// 移除unix socket
	syscall.Unlink(addr)

	log.Info("start Master :", addr)
	l, err := grpc_unixsocket.Listener(addr)
	if err != nil {
		log.Error("start service error: ", err)
		return err
	}
	var opts = []grpc.ServerOption{
		grpc.MaxRecvMsgSize(64 * 1024 * 1024),
		grpc.MaxSendMsgSize(64 * 1024 * 1024),
	}
	grpcServer := grpc.NewServer(opts...)
	m.dispatcherServe = NewDispatcherServer()
	service.RegisterCtiServiceServer(grpcServer, cli.NewMasterCliServer(m.etcdServer))
	service.RegisterMasterDispatcherServer(grpcServer, m.dispatcherServe)

	service.RegisterMasterEventsServer(grpcServer, NewEventService(m.etcdServer))
	go func() {
		err := grpcServer.Serve(l)
		if err != nil {
			log.Error("listen serve error: ", err)
		}
	}()
	log.Info("start service successful")
	m.masterSrv = grpcServer
	return nil
}

func (m *Master) stopService() {
	if m.masterSrv == nil {
		return
	}
	m.masterSrv.GracefulStop()
	//syscall.Unlink(addr)
}
