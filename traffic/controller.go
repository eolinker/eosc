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
	ms := c.data.All()
	pts := make([]*PbTraffic, 0, len(ms))
	files := make([]*os.File, 0, len(ms))
	i := 0
	for _, matcher := range ms {
		ts := matcher.Listeners()
		for _, ln := range ts {

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
			i++
		}
	}

	return pts, files
}

func (c *Controller) Shutdown() {
	c.locker.Lock()
	list := c.data.All()
	c.data = NewMatcherData()
	c.locker.Unlock()
	for _, it := range list {
		it.Close()
	}
}

func (c *Controller) reset(addrs []*net.TCPAddr) error {

	c.locker.Lock()
	defer c.locker.Unlock()

	newData := NewMatcherData()

	old := c.data.Clone()
	datas := make(map[string]*net.TCPListener)
	for _, ms := range old.All() {
		for _, l := range ms.Listeners() {
			datas[l.Addr().String()] = l
		}
	}

	for _, ad := range addrs {
		name := addrName(ad)

		o, has := datas[name]
		if has {
			delete(datas, name)
		} else {
			log.Debug("create traffic:", name)

			l, err := net.ListenTCP("tcp", ad)
			if err != nil {
				log.Error("listen tcp:", err)
				return err
			}
			o = l
		}
		om := old.Get(ad.Port)
		if om == nil {
			om = newData.Get(ad.Port)
		}
		if om == nil {
			om = NewMatcher(ad.Port, o)
		}
		newData.Set(ad.Port, om)

	}
	for n, o := range old.All() {
		log.Debug("close old : ", n)
		o.Close()
		log.Debug("close old done:", n)
	}
	c.data = newData
	return nil
}
func addrName(addr *net.TCPAddr) string {
	add := *addr
	if add.IP == nil {
		add.IP = net.IPv6zero
	}
	return add.String()
}
func rebuildAddr(addrs []*net.TCPAddr) map[int][]net.IP {
	addrMap := make(map[int][]*net.TCPAddr)
	for _, ad := range addrs {
		addrMap[ad.Port] = append(addrMap[ad.Port], ad)
	}
	newAddr := make(map[int][]net.IP)
	for p, ads := range addrMap {
		var ipv4zero bool = false
		var ipv6zero bool = false

		var ipv4s []net.IP
		var ipv6s []net.IP
		for _, ad := range ads {
			if ad.IP == nil {
				ipv4zero = false
				ipv6zero = false
				ipv4s = nil
				ipv6s = nil
				break
			}
			if ad.IP.Equal(net.IPv4zero) {
				ipv4zero = true
				continue
			}
			if ad.IP.Equal(net.IPv6zero) {
				ipv6zero = true
				continue
			}
			if len(ad.IP) == net.IPv4len {
				ipv4s = append(ipv4s, ad.IP)
			} else if len(ad.IP) == net.IPv6len {
				ipv6s = append(ipv6s, ad.IP)
			}

		}
		if ipv4zero {
			ipv4s = []net.IP{net.IPv4zero}
		}
		if ipv6zero {
			ipv6s = []net.IP{net.IPv6zero}
		}

		newAddr[p] = append(newAddr[p], ipv4s...)
		newAddr[p] = append(newAddr[p], ipv6s...)
	}
	return newAddr
}

func ReadController(r io.Reader, addr ...*net.TCPAddr) (IController, error) {
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
	err := c.reset(addr)
	if err != nil {
		return nil, err
	}
	return c, nil
}
