/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package traffic

import (
	"fmt"
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
	log.Debug("traffic controller: Export: begin ", startIndex)
	ms := c.data.All()
	pts := make([]*PbTraffic, 0, len(ms))
	files := make([]*os.File, 0, len(ms))
	i := 0
	for addr, ln := range ms {

		file, err := ln.File()
		if err != nil {
			continue
		}
		pt := &PbTraffic{
			FD:      uint64(i + startIndex),
			Addr:    addr,
			Network: ln.Addr().Network(),
		}
		pts = append(pts, pt)
		files = append(files, file)
		i++

	}
	log.Debug("traffic controller: Export: size ", len(files))

	return pts, files
}

func (c *Controller) Shutdown() {
	c.locker.Lock()
	list := c.data.All()
	c.data = NewMatcherData(nil)
	c.locker.Unlock()
	for _, it := range list {
		it.Close()
	}
}

func (c *Controller) reset(addrs []string) error {

	c.locker.Lock()
	defer c.locker.Unlock()

	old := c.data.clone()
	datas := make(map[string]*net.TCPListener)

	for _, ad := range addrs {

		v, has := old[ad]
		if has {
			delete(datas, ad)
		} else {
			log.Debug("create traffic:", ad)

			l, err := net.Listen("tcp", ad)
			if err != nil {
				log.Error("listen tcp:", err)
				return err
			}
			v = l.(*net.TCPListener)
		}
		datas[ad] = v
	}
	for n, o := range old {
		log.Debug("close old : ", n)
		o.Close()
		log.Debug("close old done:", n)
	}
	c.data = NewMatcherData(datas)
	return nil
}

func ReadController(r io.Reader, addrs ...string) (IController, error) {
	c := &Controller{
		Traffic: nil,
	}
	if r != nil {
		traffics, err := readTraffic(r)
		if err != nil {
			return nil, err
		}
		c.Traffic = NewTraffic(traffics)
	} else {
		c.Traffic = NewTraffic(nil)
	}

	err := c.reset(unrepeated(addrs...))
	if err != nil {
		return nil, err
	}
	return c, nil
}

func unrepeated(addrs ...string) []string {

	zeros := make(map[int]struct{})

	for _, ad := range addrs {
		ip, port := readAddr(ad)
		if ip == "" || ip == "0.0.0.0" {
			zeros[port] = struct{}{}
		}
	}

	data := make(map[string]struct{})
	for port := range zeros {
		data[fmt.Sprintf(":%d", port)] = struct{}{}
	}

	for _, ad := range addrs {
		_, port := readAddr(ad)
		if _, has := zeros[port]; !has {
			data[ad] = struct{}{}
		}

	}

	rs := make([]string, len(data))
	for addr := range data {
		rs = append(rs, addr)
	}
	return rs
}
