/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package process_listener

import (
	"context"
	"errors"
	"fmt"
	"github.com/eolinker/eosc/service"
	"io"
	"log"
	"net"
	"os"
	"syscall"
)

type Listener struct {

	cli service.Master_AcceptClient
	addr net.Addr
}

func NewListener(masterCli service.MasterClient,port int) (*Listener,error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil,err
	}

	cli, err := masterCli.Accept(context.Background(),&service.ListenerRequest{
		Port:int32(port),
	})
	if err != nil {
		return nil,err
	}

	return &Listener{cli:cli,addr: addr},nil
}

func (l *Listener) Accept() (net.Conn, error) {
	for {
		m,err:=l.cli.Recv()
		if err != nil{
			if errors.As(err, &io.EOF){
				return nil,err
			}
			log.Printf("accept recv:%v\n",err)
			continue
			//return nil,err
		}
		scms,err:=syscall.ParseSocketControlMessage(m.FD)
		if err != nil{
			log.Printf("accept ParseSocketControlMessage :%v\n",err)

		}
		if len(scms) >0{
			//fd:=uintptr(m.FD)

			fds,err:= syscall.ParseUnixRights(&(scms[0]))
			if err!= nil{
				log.Printf("accept ParseUnixRights:%v\n",err)

				continue
			}
			log.Println("fds:",fds)
			f:=os.NewFile(uintptr(fds[0]),"")
			if f==nil{
				continue
 			}

			n,err:= net.FileConn(f)
			if err!= nil{
				log.Printf("accept fileconn:%v" ,err)
				continue
 			}
			return n,nil
		}

	}

}

func (l *Listener) Close() error {
	return l.cli.CloseSend()
}

func (l *Listener) Addr() net.Addr {
	return l.addr
}

