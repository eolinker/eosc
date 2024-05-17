package proxy

import (
	"crypto/tls"
	"errors"
	"github.com/eolinker/eosc/log"
	"net"
	"strings"
)

type IRaftLeader interface {
	IsLeader() (bool, []string)
}

func ProxyToLeader(leader IRaftLeader, proxy *UnixProxy) func(conn net.Conn) {
	return func(conn net.Conn) {
		isLeader, peers := leader.IsLeader()
		if isLeader {
			proxy.ProxyToUnix(conn)
			return
		}

		target, err := dialPeer(peers)
		if err != nil {
			log.Info(err)
			return
		}
		doProxy(conn, target)
	}
}

var (
	tlsConf = &tls.Config{
		InsecureSkipVerify: true,
	}
)

func dialPeer(peers []string) (net.Conn, error) {
	for _, peer := range peers {
		addr := peer
		if strings.HasPrefix(addr, "https://") {
			conn, err := tls.Dial("tcp", addr[8:], tlsConf)
			if err != nil {
				log.DebugF("dial peer %s error:%s", addr, err.Error())
				continue
			}
			return conn, nil
		}
		if strings.HasPrefix(addr, "http://") {
			conn, err := net.Dial("tcp", addr[7:])
			if err != nil {
				log.DebugF("dial peer %s error:%s", addr, err.Error())
				continue
			}
			return conn, nil
		}

		conn, err := net.Dial("tcp", addr)
		if err != nil {
			log.DebugF("dial peer %s error:%s", addr, err.Error())
			continue
		}
		return conn, nil
	}
	return nil, errors.New("no peer")
}
