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
	"time"

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

	t := time.NewTicker(time.Millisecond)
	defer t.Stop()
	for {
		conn, err := net.DialUnix("unix", nil, unixAddress)
		if err == nil {
			return conn, nil
		}
		//log.Info("dail unix:", err)
		select {
		case <-ctx.Done():
			return nil, err
		case <-t.C:

		}
	}

}
func Connect(addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithContextDialer(UnixConnect))
	if err != nil {

		return nil, fmt.Errorf("did not connect: %w", err)
	}
	return conn, nil
}

func Listener(addr string) (net.Listener, error) {
	defer utils.Timeout(fmt.Sprint("port-reqiure unix:", addr))()
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
	log.Debug("unix port-reqiure close:", u.Addr().String())
	err := u.UnixListener.Close()
	if err != nil {
		log.Warn("close unix port-reqiure:", err)
		return err
	}
	//err = syscall.Unlink(u.Addr().String())
	//if err != nil {
	//	log.Warn("Unlink unix port-reqiure:", err)
	//	fmt.Println("Unlink unix port-reqiure:", err)
	//	return err
	//}
	return nil
}
