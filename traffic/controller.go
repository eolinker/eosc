/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package traffic

import (
	"io"
	"net"
	"os"
	"sync"

	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/utils"
	"google.golang.org/protobuf/proto"
)

var (
	_ IController = (*Controller)(nil)
)

type IController interface {
	ITraffic
	Close()
	Reset(ports []int) (isCreate bool, err error)
	Export(int) ([]*PbTraffic, []*os.File)
}

type Controller struct {
	locker sync.Mutex
	data   *tTrafficData
}

func (c *Controller) IsStop() bool {
	return false
}

func (c *Controller) Export(startIndex int) ([]*PbTraffic, []*os.File) {
	log.Debug("traffic controller: Export:")
	ts := c.All()
	pts := make([]*PbTraffic, 0, len(ts))
	files := make([]*os.File, 0, len(ts))
	for i, ln := range ts {
		file, err := ln.File()
		if err != nil {
			continue
		}
		addr := ln.Addr()
		pt := &PbTraffic{
			FD:      uint64(i + startIndex),
			Addr:    addr.String(),
			Network: addr.Network(),
		}
		pts = append(pts, pt)
		files = append(files, file)
	}
	return pts, files
}

func (c *Controller) Close() {
	c.locker.Lock()
	list := c.data.list()
	c.data = newTTrafficData()
	c.locker.Unlock()
	for _, it := range list {
		it.shutdown()
	}
}

func (c *Controller) Reset(ports []int) (bool, error) {
	c.locker.Lock()
	defer c.locker.Unlock()

	isCreate := false
	newData := newTTrafficData()

	old := c.data.clone()

	for _, p := range ports {
		addr := ResolveTCPAddr("", p)
		name := addrToName(addr)
		if o, has := old.Del(name); has {
			log.Debug("move traffic:", name)
			newData.add(o)
		} else {
			log.Debug("create traffic:", name)
			l, err := net.ListenTCP("tcp", addr)
			if err != nil {
				log.Error("listen tcp:", err)
				return false, err
			}
			newData.add(newTTcpListener(l))
			isCreate = true
		}
	}
	for n, o := range old.All() {
		log.Debug("close old : ", n)
		o.shutdown()
		log.Debug("close old done:", n)
	}
	c.data = newData
	return isCreate, nil
}

func (c *Controller) Encode(startIndex int) ([]byte, []*os.File, error) {

	pt, files := c.Export(startIndex)
	pts := &PbTraffics{
		Traffic: pt,
	}
	data, err := proto.Marshal(pts)
	if err != nil {
		return nil, nil, err
	}

	return utils.EncodeFrame(data), files, nil

}

func (c *Controller) All() []*tListener {

	c.locker.Lock()
	list := c.data.list()
	c.locker.Unlock()

	return list
}

func NewController(r io.Reader) IController {
	c := &Controller{
		data: newTTrafficData(),
	}
	if r != nil {
		c.data.Read(r)
	}
	return c
}

func (c *Controller) ListenTcp(ip string, port int) (net.Listener, error) {
	tcpAddr := ResolveTCPAddr(ip, port)
	c.locker.Lock()
	defer c.locker.Unlock()
	tcp, has := c.data.get(addrToName(tcpAddr))
	if !has {
		log.Warn("get listen tcp not exist")

		//tcpAddr := ResolveTCPAddr(ip, port)

		l, err := net.ListenTCP("tcp", tcpAddr)
		if err != nil {
			log.Warn("listen tcp:", err)
			return nil, err
		}
		ln := newTTcpListener(l)
		c.data.add(ln)
		tcp = ln
	}
	return tcp, nil
}
