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
	"log"
	"os"
	"syscall"

	"github.com/eolinker/eosc/process"

	grpc_unixsocket "github.com/eolinker/eosc/grpc-unixsocket"
	"google.golang.org/grpc"
)

var pidSuffix = "pid"

//StartMaster 开启master
func StartMaster(addr string) (*grpc.Server, error) {
	// 先检查是否有pid，如果pid不存在，则unlink socket文件
	path := fmt.Sprintf("%s.%s", process.AppName(), pidSuffix)
	err := process.CheckPIDFILEAlreadyExists(path)
	if err != nil {
		// 存在，则报错开启失败
		return nil, fmt.Errorf("the master is running")
	}
	// 移除unix socket
	syscall.Unlink(addr)

	log.Println("start Master :", addr)
	l, err := grpc_unixsocket.Listener(addr)
	if err != nil {
		return nil, err
	}
	err = process.CreatePidFile(path)
	if err != nil {
		// 创建pid文件失败，则报错
		os.Exit(1)
	}
	grpcServer := grpc.NewServer()

	go func() {
		grpcServer.Serve(l)
	}()
	return grpcServer, nil
}
