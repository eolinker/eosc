/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package grpc_unixsocket

import (
	"context"
	"fmt"
	"github.com/eolinker/eosc/log"
	"google.golang.org/grpc"

	"net"
)

func UnixConnect(context.Context, string) (net.Conn, error) {
	unixAddress, err := net.ResolveUnixAddr("unix", "/tmp/default.sock")
	if err != nil{
		return nil,err
	}
	conn, err := net.DialUnix("unix", nil, unixAddress)
	return conn, err
}
func Connect(addr string)(*grpc.ClientConn,error) {
	conn, err := grpc.Dial("/tmp/default.sock", grpc.WithContextDialer(UnixConnect))
	if err != nil {
		return nil,fmt.Errorf("did not connect: %w", err)
	}
	return conn,nil
}

func Listener(addr string)(net.Listener,error ){
	serverAddress, err := net.ResolveUnixAddr("unix", addr)
	if err!= nil{
		return nil, err
	}
	listen, listenErr := net.ListenUnix("unix", serverAddress)
	if listenErr != nil {
		log.Errorf("listenErr: %v", listenErr)
		return nil, listenErr
	}
   	return listen,nil
}