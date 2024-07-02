package config

import (
	"fmt"
	"github.com/eolinker/eosc/env"
	"github.com/eolinker/eosc/log"
	"github.com/ghodss/yaml"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	lastVersion = 2
)

var (
	defaultPort = map[string]int{
		"https": 443,
		"http":  80,
		"tcp":   80,
		"ssl":   443,
		"tls":   443,
	}
)

type VersionConfig struct {
	Version int `json:"version" yaml:"version"`
}
type CertConfig struct {
	Cert string `json:"cert" yaml:"cert"`
	Key  string `json:"key" yaml:"key"`
}
type UrlConfig struct {
	ListenUrl
	Certificate []CertConfig `json:"certificate,omitempty" yaml:"certificate,omitempty"`
}
type ListenUrl struct {
	ListenUrls    []string `json:"listen_urls" yaml:"listen_urls"`
	AdvertiseUrls []string `json:"advertise_urls,omitempty" yaml:"advertise_urls,omitempty"`
}

type CertificateDir struct {
	Dir string `json:"dir" yaml:"dir"`
}
type NConfig struct {
	Version        int             `json:"version" yaml:"version"`
	CertificateDir *CertificateDir `json:"certificate" yaml:"certificate"`
	Peer           UrlConfig       `json:"peer"`
	Client         UrlConfig       `json:"client"`
	Gateway        ListenUrl       `json:"gateway" yaml:"gateway"`
}
type Certificate struct {
	Key  string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Cert string `protobuf:"bytes,2,opt,name=cert,proto3" json:"cert,omitempty"`
}

func GetListens(ucs ...ListenUrl) []string {

	urls := make([]string, 0, len(ucs))
	for _, uc := range ucs {
		urls = append(urls, uc.ListenUrls...)
	}
	return FormatListenUrl(urls...)
}

func FormatListenUrl(urls ...string) []string {
	addrs := make(map[string]struct{})
	for _, lu := range urls {
		u, err := url.Parse(lu)
		if err != nil {
			continue
		}

		port, _ := strconv.Atoi(u.Port())
		if port == 0 {
			port = defaultPort[strings.ToLower(u.Scheme)]
		}
		addr := net.TCPAddr{
			IP:   net.ParseIP(u.Hostname()),
			Port: port,
			Zone: "",
		}

		addrs[addr.String()] = struct{}{}

	}
	rs := make([]string, 0, len(addrs))
	for u := range addrs {
		rs = append(rs, u)
	}
	return rs
}
func readConfigData() ([]byte, string, error) {
	paths := env.ConfigPath()

	var err error
	var data []byte

	for _, path := range paths {
		data, err = os.ReadFile(path)
		if err == nil {
			return data, path, nil
		}

	}

	if err != nil {
		return nil, "", fmt.Errorf("read config fail in:[%s]", strings.Join(paths, ","))
	}

	return nil, "", fmt.Errorf("need config")

}
func Load() NConfig {

	var config *NConfig
	var upGradle = false
	data, path, err := readConfigData()
	if err != nil {
		config = new(NConfig)
	} else {
		config, upGradle, err = readConfig(data)
		if err != nil {
			log.Warn("read config:", err)
			config = new(NConfig)
			upGradle = true
		}
	}

	if upGradle {
		rebuild, _ := yaml.Marshal(config)
		os.WriteFile(path, rebuild, 0644)
	}
	initial(config)
	return *config
}

func readConfig(data []byte) (config *NConfig, upGrade bool, err error) {
	version := new(VersionConfig)
	err = yaml.Unmarshal(data, version)
	if err != nil {
		return
	}
	config = new(NConfig)
	switch version.Version {
	case 2:

		err = yaml.Unmarshal(data, config)
		if err != nil {
			return
		}
	default:
		o := new(OConfig)
		err = yaml.Unmarshal(data, o)
		if err != nil {
			return
		}
		upGrade = true
		config.Version = lastVersion
		config.CertificateDir = o.CertificateDir
		var ssl []*ListenConfig
		if o.SSL != nil {
			ssl = o.SSL.Listen
		}
		config.Gateway = toGateway(o.Listen, ssl)

		config.Peer, config.Client = fromAdmin(o.Admin)
	}
	return
}

func initial(c *NConfig) {

	if len(c.Peer.ListenUrls) == 0 {
		c.Peer.ListenUrls = []string{"http://0.0.0.0:9401"}
		c.Peer.Certificate = nil
	}
	if len(c.Client.ListenUrls) == 0 {
		c.Client.ListenUrls = []string{"http://0.0.0.0:9400"}
		c.Client.Certificate = nil
	}

	if len(c.Gateway.ListenUrls) == 0 {
		c.Gateway.ListenUrls = []string{"http://0.0.0.0:80"}
	}

	if len(c.Peer.AdvertiseUrls) == 0 {
		peerAdvertiseUrls, has := os.LookupEnv("INIT_PEER_ADVERTISE_URLS")
		if has {
			c.Peer.AdvertiseUrls = strings.Split(peerAdvertiseUrls, ",")

		} else {
			c.Peer.AdvertiseUrls = createAdvertiseUrls(c.Peer.ListenUrls)
		}
	}
	c.Peer.AdvertiseUrls = parseHostName(c.Peer.AdvertiseUrls)
	if len(c.Client.AdvertiseUrls) == 0 {
		c.Client.AdvertiseUrls = createAdvertiseUrls(c.Client.ListenUrls)
	}
	c.Client.AdvertiseUrls = parseHostName(c.Client.AdvertiseUrls)
	if len(c.Gateway.AdvertiseUrls) == 0 {
		c.Gateway.AdvertiseUrls = createAdvertiseUrls(c.Gateway.ListenUrls)
	}

	c.Gateway.AdvertiseUrls = parseHostName(c.Gateway.AdvertiseUrls)
}
