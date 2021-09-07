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

	"github.com/eolinker/eosc/process"
	"github.com/eolinker/eosc/log"
	grpc_unixsocket "github.com/eolinker/eosc/grpc-unixsocket"
	"google.golang.org/grpc"
)

var pidSuffix = "pid"

//StartMaster 开启master
func StartMaster(addr string) (*grpc.Server, error) {

	if process.CheckPIDFILEAlreadyExists() {
		// 存在，则报错开启失败
		return nil, fmt.Errorf("the master is running")
	}
	// 移除unix socket
	syscall.Unlink(addr)

	log.Info("start Master :", addr)
	l, err := grpc_unixsocket.Listener(addr)
	if err != nil {
		return nil, err
	}
	err = process.CreatePidFile()
	if err != nil {
		// 创建pid文件失败，则报错
		return nil,err
	}
	grpcServer := grpc.NewServer()

	go func() {
		grpcServer.Serve(l)
	}()
	return grpcServer, nil
}
