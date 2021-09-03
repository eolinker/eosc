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
	"io"
	"net"
	"os"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/utils"
	"google.golang.org/protobuf/proto"
)

var (
	ErrorInvalidFiles = errors.New("invalid files")
)

type Traffics []*TrafficOut

//IController 端口管理器接口
type IController interface {
	Listener(network string, addr string) error
	All() Traffics
	Close()
}

type Controller struct {
	data eosc.IUntyped
}

func (c *Controller) All() Traffics {
	list := c.data.List()
	ts := make(Traffics, 0, len(list))
	for _, it := range list {
		tf, ok := it.(*TrafficOut)
		if !ok {
			continue
		}
		ts = append(ts, tf)
	}
	return ts
}

//NewController 新建端口管理器（流量入口）
func NewController() *Controller {
	return &Controller{
		data: eosc.NewUntyped(),
	}
}

func (ts Traffics) WriteTo(w io.Writer) ([]*os.File, error) {

	pts := new(PbTraffics)
	files := make([]*os.File, 0, len(ts))
	pts.Traffic = make([]*PbTraffic, 0, len(ts))
	for i, it := range ts {

		pt := &PbTraffic{
			FD:      uint64(i),
			Addr:    it.Addr.String(),
			Network: it.Addr.Network(),
		}
		pts.Traffic = append(pts.Traffic, pt)
		files = append(files, it.File)
	}

	data, err := proto.Marshal(pts)
	if err != nil {
		return nil, err
	}
	fmt.Println(data)
	err = utils.WriteFrame(w, data)
	if err != nil {
		return nil, err
	}
	return files, nil
}

//Listener 设置端口监听器，如果地址已经被监听，则报错
func (c *Controller) Listener(network string, addr string) error {
	tcpAddr, err := net.ResolveTCPAddr(network, addr)
	if err != nil {
		return err
	}

	l, err := net.ListenTCP(network, tcpAddr)
	if err != nil {
		return err
	}
	file, err := l.File()
	if err != nil {
		return err
	}
	tf := &TrafficOut{
		Addr: tcpAddr,
		File: file,
	}

	name := fmt.Sprintf("%s://%s", tcpAddr.Network(), tcpAddr.String())
	c.data.Set(name, tf)
	return nil
}

//Close 关闭文件监听
func (c *Controller) Close() {
	list := c.data.List()
	c.data = eosc.NewUntyped()

	for _, it := range list {
		tf, ok := it.(*TrafficOut)
		if !ok {
			continue
		}
		tf.File.Close()
	}
}
