/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 *
 */
/*

io 通信控制模块
管理所有的需要热重启的监听管理（端口监听）， 只允许master执行新增， 序列化成描述信息+文件描述符列表，在fork worker时传递给worker，
worker只允许使用传入进来的端口

 */
package traffic

import (
	"fmt"
	"github.com/eolinker/eosc/utils"
	"go.etcd.io/etcd/Godeps/_workspace/src/github.com/golang/protobuf/proto"
	"io"
	"log"
	"net"
	"os"
)

type TrafficIn struct {
	Addr net.Addr
	Listener net.Listener
	File *os.File
}
type TrafficOut struct {
	Addr net.Addr
	File *os.File
}

func Reader(r io.Reader,start int)([]*TrafficIn,error) {

	frame, err := utils.ReadFrame(r)
	if err != nil {
		return nil, err
	}
	log.Println("readed frame:",len(frame))

	pts:=new(PbTraffics)
	err = proto.Unmarshal(frame, pts)
	if err != nil {
		return nil,err
	}
	log.Println("triv:",pts)
	tfs:=make([]*TrafficIn,0,len(pts.GetTraffic()))
	for _,pt:=range pts.GetTraffic(){
		name:=fmt.Sprintf("%s:/%s",pt.Network,pt.Addr)

		addr, err := net.ResolveTCPAddr(pt.GetNetwork(), pt.GetAddr())
		if err != nil {
			return nil, err
		}
		f:=os.NewFile(uintptr(uint64(start)+pt.GetFD()),name)
		switch pt.Network{
		//case "udp","udp4","udp8":
		//
		//	c,err:=net.FilePacketConn(f)
		case "tcp","tcp4","tcp6":
			l,err:= net.FileListener(f)
			if err!= nil{
				log.Println("error to read listener:",err)
				return nil, err
			}
			tfs = append(tfs, &TrafficIn{
				Addr: addr,
				Listener: l,
				File: f,
			})
		}
	}

	return tfs, nil
}