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
	"github.com/soheilhy/cmux"
	"net"
	"net/url"
	"strconv"
	"strings"
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

const (
	bitTCP = 1 << iota
	bitSSL

	bitBoth = bitTCP | bitSSL
)

func (t *Traffic) Listen(addrs ...string) (tcp []net.Listener, ssl []net.Listener) {

	schemes := make(map[string]int)
	for _, addr := range addrs {
		addrValue, isSSl := readAddr(addr)
		if isSSl {
			schemes[addrValue] = schemes[addrValue] | bitSSL
		} else {
			schemes[addrValue] = schemes[addrValue] | bitTCP
		}
	}
	for addr, v := range schemes {

		listener, has := t.data[addr]
		if !has {
			continue
		}
		switch v {
		case bitBoth:
			{
				cMux := cmux.New(listener)
				ssl = append(ssl, cMux.Match(cmux.TLS()))
				tcp = append(tcp, cMux.Match(cmux.Any()))

			}
		case bitTCP:
			tcp = append(tcp, listener)
		case bitSSL:
			ssl = append(ssl, listener)
		}

	}
	return tcp, ssl
}
func readAddr(addr string) (string, bool) {
	u, err := url.Parse(addr)
	if err != nil {
		u = &url.URL{Scheme: "tcp", Host: addr}
	}
	ssl := false
	port, _ := strconv.Atoi(u.Port())

	switch strings.ToLower(u.Scheme) {
	case "https", "ssl", "tls":
		ssl = true
		if port == 0 {
			port = 443
		}
	default:
		ssl = false
		if port == 0 {
			port = 80
		}
	}
	parseIP := net.ParseIP(u.Hostname())

	tcpAddr := net.TCPAddr{IP: parseIP, Port: port}

	return tcpAddr.String(), ssl
}

func NewTraffic(trafficData *TrafficData) ITraffic {
	return &Traffic{TrafficData: trafficData}
}
func FromArg(traffics []*PbTraffic) ITraffic {
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
