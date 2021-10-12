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
	"net"

	"github.com/eolinker/eosc/utils"

	"github.com/eolinker/eosc/log"
	"google.golang.org/grpc"
)

func UnixConnect(ctx context.Context, addr string) (net.Conn, error) {
	log.Debug("UnixConnect:", addr)
	unixAddress, err := net.ResolveUnixAddr("unix", addr)
	if err != nil {
		log.Debug("ResolveUnixAddr:", addr, ":", err)
		return nil, err
	}
	conn, err := net.DialUnix("unix", nil, unixAddress)
	if err != nil {
		log.Info("dail unix:", err)
	}
	return conn, err
}
func Connect(addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithContextDialer(UnixConnect))
	if err != nil {

		return nil, fmt.Errorf("did not connect: %w", err)
	}
	return conn, nil
}

func Listener(addr string) (net.Listener, error) {
	defer utils.Timeout(fmt.Sprint("listener unix:", addr))()
	serverAddress, err := net.ResolveUnixAddr("unix", addr)
	if err != nil {
		return nil, err
	}

	listen, listenErr := net.ListenUnix("unix", serverAddress)
	if listenErr != nil {
		log.Errorf("listenErr: %v", listenErr)
		return nil, listenErr
	}
	return &unixListener{listen}, nil
}

type unixListener struct {
	*net.UnixListener
}

func (u *unixListener) Close() error {
	log.Debug("unix listener close:", u.Addr().String())
	err := u.UnixListener.Close()
	if err != nil {
		log.Warn("close unix listener:", err)
		return err
	}
	//err = syscall.Unlink(u.Addr().String())
	//if err != nil {
	//	log.Warn("Unlink unix listener:", err)
	//	fmt.Println("Unlink unix listener:", err)
	//	return err
	//}
	return nil
}
