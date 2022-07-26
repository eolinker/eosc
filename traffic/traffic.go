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
	cmuxMatch "github.com/eolinker/eosc/traffic/cmux-match"
	"net"
	"sync"

	"github.com/eolinker/eosc/log"
)

var (
	ErrorInvalidListener          = errors.New("invalid port-reqiure")
	_                    ITraffic = (*Traffic)(nil)
	_                    ITraffic = (*EmptyTraffic)(nil)
)

type TrafficType = cmuxMatch.MatchType

const (
	Any TrafficType = iota
	Http1
	Https
	Http2
	Websocket
	GRPC
)

type Traffic struct {
	locker sync.Mutex
	data   *tTrafficData
	matchs map[string]cmuxMatch.CMuxMatch
	stop   bool
}

func (t *Traffic) IsStop() bool {
	return t.stop
}

func NewTraffic() *Traffic {
	return &Traffic{
		data:   newTTrafficData(),
		locker: sync.Mutex{},
		matchs: map[string]cmuxMatch.CMuxMatch{},
	}
}
func (t *Traffic) Read(tfConf []*PbTraffic) error {
	t.locker.Lock()
	defer t.locker.Unlock()
	data := newTTrafficData()
	data.setListener(tfConf)
	t.data = data
	return nil
}
func (t *Traffic) get(ip string, port int) net.Listener {
	tcpAddr := ResolveTCPAddr(ip, port)
	name := addrToName(tcpAddr)
	if o, has := t.data.get(name); has {
		log.Debug("traffic ListenTcp:", ip, ":", port, ", ok")
		return o
	}
	ipv := resolve(ip)
	if ipv.Equal(net.IPv4zero) && ipv.Equal(net.IPv6zero) {
		return nil
	}
	return t.get(ipv.String(), port)
}
func (t *Traffic) ListenTcp(ip string, port int, trafficType TrafficType) net.Listener {
	log.Debug("traffic try ListenTcp:", ip, ":", port)
	//tcpAddr := ResolveTCPAddr(ip, port)
	//name := addrToName(tcpAddr)
	t.locker.Lock()
	defer t.locker.Unlock()
	l := t.get(ip, port)
	if l == nil {
		return nil
	}
	return t.match(l, trafficType)
}
func (t *Traffic) match(l net.Listener, trafficType TrafficType) net.Listener {
	name := l.Addr().String()
	matcher, has := t.matchs[name]
	if !has {
		matcher = cmuxMatch.NewMatch(l)
		t.matchs[name] = matcher
	}
	return matcher.Match(trafficType)
}

type ITraffic interface {
	ListenTcp(ip string, port int, trafficType TrafficType) net.Listener
	IsStop() bool
	Close()
}

func (t *Traffic) Close() {
	t.locker.Lock()
	list := t.data.list()
	t.data = newTTrafficData()
	t.locker.Unlock()
	for _, it := range list {
		err := it.Close()
		if err != nil {
			log.Info("close traffic port-reqiure:", err)
		}
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

func toName(ln net.Listener) string {
	addr := ln.Addr()
	return addrToName(addr)
}

func addrToName(addr net.Addr) string {
	return fmt.Sprintf("%s://%s", addr.Network(), addr.String())

}

type EmptyTraffic struct {
}

func NewEmptyTraffic() *EmptyTraffic {
	return &EmptyTraffic{}
}

func (e *EmptyTraffic) ListenTcp(ip string, port int, trafficType TrafficType) net.Listener {
	return nil
}

func (e *EmptyTraffic) IsStop() bool {
	return true
}

func (e *EmptyTraffic) Close() {
	return
}
