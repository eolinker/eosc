/*
 * Copyright (c) 2021-2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package service

import (
	"context"
	"fmt"
	"github.com/eolinker/eosc/listener"
	"github.com/eolinker/eosc/service"
	"log"
	"net"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

type MasterServer struct {
	service.UnimplementedMasterServer
	sign string
}

func NewMasterServer() *MasterServer {
	return &MasterServer{sign: fmt.Sprintf("%d",time.Now().UnixNano())}
}

func (m *MasterServer) Accept(request *service.ListenerRequest, server service.Master_AcceptServer) error {
	lr, err := listener.ListenTCP(int(request.Port), m.sign)
	if err!= nil{
		return err
	}
	for{
		c,e:=lr.Accept()
		if e!=nil{
			return e
		}
		
		f,e:=c.(*net.TCPConn).File()
		if e!= nil{
			log.Println("tcp connect to file :",e.Error())
			continue
		}

		fd:=f.Fd()

		log.Println("master send:fd=",fd)
		er:=server.Send(&service.AcceptResponse{
			Status:  0,
			FD:      syscall.UnixRights(int(fd)),
			Scheme:  request.Scheme,
			Port:    request.Port,
			Message: "",
		})
		if er!= nil{
			log.Println(err)
			return er
		}
	}
	return nil
}

func (m *MasterServer) Open(ctx context.Context, open *service.FileOpen) (*service.FileHandler, error) {
	path,_:= filepath.Abs(open.FilePath)
	file, err := os.OpenFile(path, int(open.Flag), os.FileMode(open.FileMode))
	if err!= nil{
		return nil,err
	}

	fd:=file.Fd()
	defer file.Close()
	return &service.FileHandler{
		FD:    uint64(fd),
		Code:  0,
		Name:  filepath.Base(open.FilePath),
		Path:  path,
		Error: "",
	},nil
}

