package config

import (
	"fmt"
	"github.com/eolinker/eosc/log"
	"net"
	"net/url"
	"os"
	"strings"
	"sync"
)

func createAdvertiseUrls(listenUrls []string) []string {

	urls := make(map[string]struct{})
	for _, lUrl := range listenUrls {
		u, err := url.Parse(lUrl)
		if err != nil {
			continue
		}
		port := u.Port()

		ip := u.Hostname()
		if port == "" {
			switch u.Scheme {
			case "http", "tcp":
				port = "80"
			case "https", "ssl":
				port = "443"
			}
		}
		if ip == "0.0.0.0" {

			ips := getIps()
			for _, i := range ips {
				urls[fmt.Sprintf("%s://%s:%s", u.Scheme, i, port)] = struct{}{}
			}

		} else {
			urls[fmt.Sprintf("%s://%s:%s", u.Scheme, ip, port)] = struct{}{}
		}
	}
	advertise := make([]string, 0, len(urls))
	for u := range urls {
		advertise = append(advertise, u)
	}
	return advertise

}

var (
	onceGetIp sync.Once
	ipsCache  []string
)

func getIps() []string {
	onceGetIp.Do(func() {
		addrs, err := net.InterfaceAddrs()

		if err != nil {
			log.Warn("get Ip fail:", err)
			return
		}

		ipsCache = make([]string, 0, len(addrs))
		for _, address := range addrs {

			// 检查ip地址判断是否回环地址
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ipsCache = append(ipsCache, ipnet.IP.String())
				}

			}
		}

	})

	return ipsCache
}

func parseHostName(advertiseUrls []string) []string {

	hostname, _ := os.Hostname()
	rs := make([]string, 0, len(advertiseUrls))
	for _, advertiseUrl := range advertiseUrls {
		rs = append(rs, strings.Replace(advertiseUrl, "${hostname}", hostname, -1))
	}
	return rs

}
