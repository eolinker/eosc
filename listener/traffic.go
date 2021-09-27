package listener

import (
	"errors"
	"net"

	"github.com/eolinker/eosc/traffic"
)

var (
	defaultTraffic traffic.ITraffic
)

func SetTraffic(t traffic.ITraffic) {
	defaultTraffic = t
}
func ListenTcp(ip string, port int) (net.Listener, error) {
	if defaultTraffic == nil {
		return nil, errors.New("traffic not init")
	}

	return defaultTraffic.ListenTcp(ip, port)
}
