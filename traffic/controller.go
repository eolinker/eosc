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

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
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
	locker sync.Mutex
	data   *tTrafficData
}

func (c *Controller) Expire(ports []int) {
	c.Reset(ports)
}

func (c *Controller) Close() {
	c.locker.Lock()
	list := c.data.list()
	c.data = newTTrafficData()
	c.locker.Unlock()
	for _, it := range list {
		err := it.Close()
		if err != nil {
			log.Info("close traffic listener:", err)
		}
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
				log.Warn("listen tcp:", err)
				return false, err
			}
			newData.add(l)
			isCreate = true
		}
	}
	for n, o := range old.All() {

		//l, ok := o.(*net.TCPListener)
		//if !ok {
		//	log.Warn("unknown error while try close  listener:", n)
		//	continue
		//}
		log.Debug("close old : ", n)
		if err := o.Close(); err != nil {
			log.Warn("close listener:", err, " ", o.Addr())
		}

		log.Debug("close old done:", n)
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
	c.locker.Lock()
	defer c.locker.Unlock()
	tcpAddr := ResolveTCPAddr(ip, port)
	tcp, has := c.data.get(addrToName(tcpAddr))

	if !has {
		log.Warn("get listen tcp not exist")

		//tcpAddr := ResolveTCPAddr(ip, port)

		l, err := net.ListenTCP("tcp", tcpAddr)
		if err != nil {
			log.Warn("listen tcp:", err)
			return nil, err
		}

		c.data.add(l)
		tcp = l
	}
	return tcp, nil
}
