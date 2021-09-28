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
	"sync"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
)

var (
	ErrorInvalidListener = errors.New("invalid listener")
)
var _ ITraffic = (*Traffic)(nil)

type Traffic struct {
	locker sync.Mutex
	data   eosc.IUntyped
}

func (t *Traffic) remove(name string) {
	t.data.Del(name)
}

func NewTraffic() *Traffic {
	return &Traffic{
		data:   eosc.NewUntyped(),
		locker: sync.Mutex{},
	}
}

func (t *Traffic) ListenTcp(ip string, port int) (net.Listener, error) {

	tcpAddr := ResolveTCPAddr(ip, port)
	t.locker.Lock()
	defer t.locker.Unlock()

	name := fmt.Sprintf("%s://%s", tcpAddr.Network(), tcpAddr.String())
	log.Debug("traffic listen:", name)
	if o, has := t.data.Get(name); has {
		listener, ok := o.(net.Listener)
		if !ok {
			return nil, ErrorInvalidListener
		}

		return listener, nil
	}

	return nil, nil
}

type ITraffic interface {
	ListenTcp(ip string, port int) (net.Listener, error)
	Close()
	remove(name string)
}

func (t *Traffic) Read(r io.Reader) {
	t.locker.Lock()
	defer t.locker.Unlock()

	listeners, err := readListener(r)
	if err != nil {
		log.Warn("read listeners:", err)
		return
	}
	for _, ln := range listeners {
		t.add(ln)
	}

}

func (t *Traffic) add(ln *net.TCPListener) {
	tcpAddr := ln.Addr()
	name := toName(tcpAddr)
	log.Info("traffic add:", name)
	t.data.Set(name, ln)
}

func (t *Traffic) Close() {
	t.locker.Lock()
	list := t.data.List()
	t.data = eosc.NewUntyped()
	t.locker.Unlock()
	for _, it := range list {
		tf, ok := it.(*net.TCPListener)
		if !ok {
			continue
		}
		tf.Close()
	}
}

func resolve(value string) net.IP {
	ip := net.ParseIP(value)
	if ip == nil {
		return net.IPv6zero
	}
	if ip.Equal(net.IPv4zero) {
		return net.IPv6zero
	}
	return ip
}

func ResolveTCPAddr(ip string, port int) *net.TCPAddr {

	return &net.TCPAddr{
		IP:   resolve(ip),
		Port: port,
		Zone: "",
	}
}
func toName(addr net.Addr) string {
	return fmt.Sprintf("%s://%s", addr.Network(), addr.String())

}
