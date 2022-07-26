/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package traffic

import (
	"github.com/eolinker/eosc/log"
	"io"
	"net"
	"os"
)

var (
	_ IController = (*Controller)(nil)
)

type IController interface {
	ITraffic
	Shutdown()
	//Reset(ports []int) (isCreate bool, err error)
	Export(int) ([]*PbTraffic, []*os.File)
}

type Controller struct {
	*Traffic
}

func (c *Controller) Export(startIndex int) ([]*PbTraffic, []*os.File) {
	log.Debug("traffic controller: Export:")
	ts := c.all()
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

func (c *Controller) Shutdown() {
	c.locker.Lock()
	list := c.data.list()
	c.data = newTTrafficData()
	c.locker.Unlock()
	for _, it := range list {
		it.shutdown()
	}
}

func (c *Controller) reset(addrs []*net.TCPAddr) error {
	c.locker.Lock()
	defer c.locker.Unlock()

	newData := newTTrafficData()

	old := c.data.clone()

	for _, addr := range addrs {

		name := addrToName(addr)
		if o, has := old.Del(name); has {
			log.Debug("move traffic:", name)
			newData.add(o)
		} else {
			log.Debug("create traffic:", name)
			l, err := net.ListenTCP("tcp", addr)
			if err != nil {
				log.Error("listen tcp:", err)
				return err
			}
			newData.add(newTTcpListener(l))
		}
	}
	for n, o := range old.All() {
		log.Debug("close old : ", n)
		o.shutdown()
		log.Debug("close old done:", n)
	}
	c.data = newData
	return nil
}

func (c *Controller) all() []*tListener {

	c.locker.Lock()
	list := c.data.list()
	c.locker.Unlock()

	return list
}

func ReadController(r io.Reader, addr ...*net.TCPAddr) (IController, error) {
	c := &Controller{
		Traffic: NewTraffic(),
	}
	if r != nil {
		c.data.Read(r)
	}
	err := c.reset(addr)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Controller) ListenTcp(ip string, port int, trafficType TrafficType) net.Listener {
	tcpAddr := ResolveTCPAddr(ip, port)

	c.locker.Lock()
	defer c.locker.Unlock()
	tcp := c.Traffic.get(ip, port)

	if tcp == nil {
		log.Warn("get listen tcp not exist")
		l, err := net.ListenTCP("tcp", tcpAddr)
		if err != nil {
			log.Warn("listen tcp:", err)
			return nil
		}
		ln := newTTcpListener(l)
		c.data.add(ln)
		tcp = ln
	}
	if tcp != nil {
		return c.match(tcp, trafficType)
	}
	return tcp
}
