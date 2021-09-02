/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package traffic

import (
	"errors"
	"fmt"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/utils"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
	"os"
)

var (
	ErrorInvalidFiles = errors.New("invalid files")
)
type Traffics []*TrafficOut

type IController interface {
	Listener(network string,addr string)error
 	All()Traffics
	Close()
}

type Controller struct {
	data eosc.IUntyped
}

func (c *Controller) All() Traffics {
	list := c.data.List()
	ts:=make(Traffics,0,len(list))
	for _,it:=range list{
		tf,ok:= it.(*TrafficOut)
		if !ok{
			continue
		}
		ts = append(ts, tf)
	}
	return ts
}

func NewController() *Controller {
	return &Controller{
		data: eosc.NewUntyped(),
	}
}

func (ts Traffics) WriteTo(w io.Writer) ([]*os.File,error){

	pts:=new(PbTraffics)
	files:=make([]*os.File,0,len(ts))
	pts.Traffic = make([]*PbTraffic,0,len(ts))
	for i,it:=range ts{

		pt:=&PbTraffic{
			FD: uint64(i),
			Addr:    it.Addr.String(),
			Network: it.Addr.Network(),
		}
		pts.Traffic = append(pts.Traffic, pt)
		files = append(files, it.File)
	}

	data, err := proto.Marshal(pts)
	if err!= nil{
		return nil, err
	}
	fmt.Println(data)
	err = utils.WriteFrame(w, data)
	if err != nil {
		return nil, err
	}
	return files,nil
}
func (c *Controller) Listener(network string, addr string) error {
	tcpAddr, err := net.ResolveTCPAddr(network, addr)
	if err != nil {
		return err
	}

	l,err:=net.ListenTCP(network,tcpAddr)
	if err!= nil{
		return err
	}
	file, err := l.File()
	if err != nil {
		return err
	}
	tf:=&TrafficOut{

		Addr: tcpAddr,
		File: file,
	}
	name:=fmt.Sprintf("%s://%s",tcpAddr.Network(),tcpAddr.String())
	c.data.Set(name,tf)
	return nil
}


func (c *Controller) Close() {
	list := c.data.List()
	c.data = eosc.NewUntyped()

	for _,it:=range list{
		tf,ok:=it.(*TrafficOut)
		if !ok{
			continue
		}
		tf.File.Close()
	}
}

