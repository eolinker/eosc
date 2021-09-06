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
	"io"
	"net"
	"sync"
)
var(
	ErrorInvalidListener = errors.New("invalid listener")
)
type Traffic struct {
	locker sync.Mutex
	data eosc.IUntyped
}

func NewTraffic() *Traffic {
	return &Traffic{
		data: eosc.NewUntyped(),
		locker: sync.Mutex{},
	}
}

func (t *Traffic) ListenTcp(network, addr string) (*net.TCPListener, error) {
	t.locker.Lock()
	defer t.locker.Unlock()

	tcpAddr, err := net.ResolveTCPAddr(network, addr)
	if err != nil {
		return nil, err
	}
	name:=fmt.Sprintf("%s://%s",tcpAddr.Network(),tcpAddr.String())

	if o, has := t.data.Get(name);has{
		listener,ok := o.(*net.TCPListener)
		if !ok{
			return nil, ErrorInvalidListener
		}

		return listener,nil
	}

	return nil,nil
}

type ITraffic interface {
	ListenTcp(network,addr string)(*net.TCPListener,error)
}

func (t *Traffic) Read(r io.Reader) {
	t.locker.Lock()
	defer t.locker.Unlock()

	listeners, err := Reader(r)
	if err != nil {
		return
	}
	for _,ln:=range listeners{
		t.add(ln)
	}

}

func (t *Traffic) add(ln *net.TCPListener)  {
	tcpAddr := ln.Addr()
	name:=fmt.Sprintf("%s://%s",tcpAddr.Network(),tcpAddr.String())
	t.data.Set(name,ln)
}
