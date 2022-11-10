package config

import (
	"fmt"
	"github.com/eolinker/eosc/log"
	"net"
	"net/url"
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
		ip := strings.TrimSuffix(u.Host, fmt.Sprintf(":%s", port))
		if port == "" {
			switch u.Scheme {
			case "http":
				port = "80"
			case "https":
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

//func splitPeerUrl(u string) (scheme, ip string, port int) {
//	url, err := url.Parse(u)
//	if err != nil {
//		return "", "", 0
//	}
//	url.Port()
//	if index := strings.Index(url, "://"); index > -1 {
//		scheme = url[:index]
//		url = url[index+3:]
//	} else
//	if index := strings.Index(url, ":"); index > -1 {
//		ip = url[:index]
//		url = url[index+1:]
//	}
//
//}

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
