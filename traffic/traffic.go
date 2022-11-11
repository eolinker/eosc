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
	"github.com/eolinker/eosc/log"
	"net"
)

var (
	ErrorInvalidListener = errors.New("invalid port-reqiure")
)

type ITraffic interface {
	Listen(addrs ...string) (tcp []net.Listener, ssl []net.Listener)
	IsStop() bool
	Close()
}

type Traffic struct {
	*TrafficData
}

func (t *Traffic) Listen(addrs ...string) (tcp []net.Listener, ssl []net.Listener) {
	//TODO implement me
	panic("implement me")
}

func NewTraffic(trafficData *TrafficData) ITraffic {
	return &Traffic{TrafficData: trafficData}
}
func TrafficFromArg(traffics []*PbTraffic) ITraffic {
	listeners := toListeners(traffics)
	log.Debug("read listeners: ", len(listeners))

	data := NewTrafficData(listeners)
	return NewTraffic(data)
}

type EmptyTraffic struct {
}

func (e *EmptyTraffic) Listen(addrs ...string) (tcp []net.Listener, ssl []net.Listener) {
	return nil, nil
}

func (e *EmptyTraffic) IsStop() bool {
	return false
}

func (e *EmptyTraffic) Close() {

}

func NewEmptyTraffic() ITraffic {
	return &EmptyTraffic{}
}
