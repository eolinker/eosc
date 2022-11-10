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
	cmuxMatch "github.com/eolinker/eosc/traffic/cmux-match"
	"net"
	"strconv"
	"strings"
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
	data   *MatcherData

	stop bool
}

func (t *Traffic) Listen(addr string) net.Listener {
	//TODO implement me
	panic("implement me")
}

func (t *Traffic) IsStop() bool {
	return t.stop
}

func NewTraffic(traffics []*PbTraffic) *Traffic {
	listeners, err := toListeners(traffics)
	log.Debug("read listeners: ", len(listeners))
	if err != nil {
		log.Warn("read listeners:", err)
	}
	data := NewMatcherData(listeners)

	tf := &Traffic{
		data:   data,
		locker: sync.Mutex{},
	}
	return tf
}

func (t *Traffic) ListenTcp(port int, trafficType TrafficType) net.Listener {
	log.Debug("traffic try ListenTcp for:", port)

	t.locker.Lock()
	defer t.locker.Unlock()
	l := t.data.GetByPort(port)
	if l == nil {
		log.Warn("listen to un open port: ", port, " for :", trafficType)
		return nil
	}

	return l.Match(trafficType)
}

type ITraffic interface {
	ListenTcp(port int, trafficType TrafficType) net.Listener
	Listen(addr string) net.Listener
	IsStop() bool
	Close()
}

func (t *Traffic) Close() {
	t.locker.Lock()
	list := t.data.All()
	for _, it := range list {
		it.Close()
	}
	t.data = NewMatcherData(nil)
	t.locker.Unlock()
}

func readAddr(ipPort string) (string, int) {

	i := strings.LastIndex(ipPort, ":")
	ip := ipPort[:i]

	port := ipPort[i+1:]
	pv, _ := strconv.Atoi(port)
	if ip == "" {
		ip = "0.0.0.0"
	}
	return ip, pv
}

type EmptyTraffic struct {
}

func (e *EmptyTraffic) Listen(addr string) net.Listener {
	//TODO implement me
	panic("implement me")
}

func NewEmptyTraffic() *EmptyTraffic {
	return &EmptyTraffic{}
}

func (e *EmptyTraffic) ListenTcp(port int, trafficType TrafficType) net.Listener {
	return nil
}

func (e *EmptyTraffic) IsStop() bool {
	return true
}

func (e *EmptyTraffic) Close() {
	return
}
