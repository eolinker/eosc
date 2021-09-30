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

	"github.com/eolinker/eosc/log"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/utils"
	"google.golang.org/protobuf/proto"
)

type IController interface {
	eosc.IDataMarshaller
	ITraffic
	Close()
	Reset(ports []int) (isCreate bool, err error)
}

type Controller struct {
	Traffic
}

func (c *Controller) Reset(ports []int) (bool, error) {
	c.locker.Lock()
	defer c.locker.Unlock()
	isCreate := false
	newData := eosc.NewUntyped()

	old := c.data.Clone()

	for _, p := range ports {
		addr := ResolveTCPAddr("", p)
		name := toName(addr)
		if o, has := old.Del(name); has {
			newData.Set(name, o)
		} else {
			l, err := net.ListenTCP("tcp", addr)
			if err != nil {
				log.Warn("listen tcp:", err)
				return false, err
			}
			newData.Set(name, l)
			isCreate = true
		}
	}
	for n, o := range old.All() {

		l, ok := o.(*net.TCPListener)
		if !ok {
			log.Warn("unknown error while try close  listener:", n)
			continue
		}
		log.Debug("close old address: ", l.Addr())
		if err := l.Close(); err != nil {
			log.Warn("close listener:", err, " ", l.Addr())
		}
	}
	c.data = newData
	return isCreate, nil
}

func (c *Controller) Encode(startIndex int) ([]byte, []*os.File, error) {
	log.Debug("traffic controller: encode:")
	ts := c.All()
	pts := new(PbTraffics)
	files := make([]*os.File, 0, len(ts))
	pts.Traffic = make([]*PbTraffic, 0, len(ts))
	for i, ln := range ts {
		file, err := ln.File()
		if err != nil {
			continue
		}
		ln.Close()
		addr := ln.Addr()
		pt := &PbTraffic{
			FD:      uint64(i + startIndex),
			Addr:    addr.String(),
			Network: addr.Network(),
		}
		pts.Traffic = append(pts.Traffic, pt)
		files = append(files, file)
	}

	data, err := proto.Marshal(pts)
	if err != nil {
		return nil, nil, err
	}

	return utils.EncodeFrame(data), files, nil

}

func (c *Controller) All() []*net.TCPListener {

	c.locker.Lock()
	list := c.data.List()
	c.locker.Unlock()

	ts := make([]*net.TCPListener, 0, len(list))
	for _, it := range list {
		tf, ok := it.(*net.TCPListener)
		if !ok {
			continue
		}
		ts = append(ts, tf)
	}

	return ts
}

func NewController(r io.Reader) IController {
	c := &Controller{
		Traffic: Traffic{
			data: eosc.NewUntyped(),
		},
	}
	if r != nil {
		c.Read(r)
	}
	return c
}

func (c *Controller) ListenTcp(ip string, port int) (net.Listener, error) {

	tcp, err := c.Traffic.ListenTcp(ip, port)
	if err != nil {
		log.Warn("get listen tcp from traffic :", err)
		return nil, err
	}
	if tcp == nil {
		log.Warn("get listen tcp not exist")
		c.locker.Lock()
		defer c.locker.Unlock()
		tcpAddr := ResolveTCPAddr(ip, port)

		l, err := net.ListenTCP("tcp", tcpAddr)
		if err != nil {
			log.Warn("listen tcp:", err)
			return nil, err
		}

		c.Traffic.add(l)
		tcp = l
	}
	return tcp, nil
}
