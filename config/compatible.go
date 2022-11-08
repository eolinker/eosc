package config

import (
	"errors"
	"fmt"
	"github.com/eolinker/eosc/env"
	"github.com/ghodss/yaml"
	"net"
	"os"
	"strings"
)

const (
	lastVersion = 1
)

type VersionConfig struct {
	Version int `json:"version" yaml:"version"`
}
type CertConfig struct {
	Cert string `json:"cert" yaml:"cert"`
	Key  string `json:"key" yaml:"key"`
}
type UrlConfig struct {
	ListenUrls    []string     `json:"listen_urls" yaml:"listen_urls"`
	Certificate   []CertConfig `json:"certificate"`
	AdvertiseUrls []string     `json:"advertise_urls" yaml:"advertise_urls"`
}

type NConfig struct {
	Version        int             `json:"version" yaml:"version"`
	CertificateDir *CertificateDir `json:"certificate" yaml:"certificate"`
	Peer           UrlConfig       `json:"peer"`
	Client         UrlConfig       `json:"client"`
	Gateway        UrlConfig       `json:"gateway" yaml:"gateway"`
}

type OConfig struct {
	Listen         []int           `json:"listen" yaml:"listen"`
	SSL            *SSLConfig      `json:"ssl" yaml:"ssl"`
	Admin          *AdminConfig    `json:"admin" yaml:"admin"`
	CertificateDir *CertificateDir `json:"certificate" yaml:"certificate"`
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
func Load() (*NConfig, error) {
	data, path, err := readConfigData()
	if err != nil {
		return nil, err
	}
	config, upGradle, err := read(data)
	if err != nil {
		return nil, err
	}
	initial(config)
	if upGradle {
		rebuild, _ := yaml.Marshal(config)
		os.WriteFile(path, rebuild, 0644)
	}
	return config, nil
}

func read(data []byte) (config *NConfig, upGrade bool, err error) {
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
func fromAdmin(admin *AdminConfig) (UrlConfig, UrlConfig) {
	peer := UrlConfig{}
	scheme := strings.ToLower(admin.Scheme)
	if scheme != "https" {
		scheme = "http"
	}
	ssl := (scheme == "https") && (admin.Certificate != nil)
	if !ssl {
		scheme = "http"
	}
	peer.ListenUrls = []string{fmt.Sprintf("%s://%s:%d", scheme, admin.IP, admin.Listen)}

	if admin.IP == "0.0.0.0" || admin.IP == "" {
		ips, _ := getIps()
		peer.AdvertiseUrls = make([]string, 0, len(ips))
		for _, ip := range ips {
			peer.AdvertiseUrls = append(peer.AdvertiseUrls, fmt.Sprintf("%s://%s:%d", scheme, ip, admin.Listen))
		}
	} else {
		peer.AdvertiseUrls = []string{fmt.Sprintf("%s://%s:%d", scheme, admin.IP, admin.Listen)}
	}
	if ssl {
		peer.Certificate = make([]CertConfig, 0, 1)
		peer.Certificate = append(peer.Certificate, CertConfig{
			Cert: admin.Certificate.Cert,
			Key:  admin.Certificate.Key,
		})
	}

	client := peer
	return peer, client
}
func toGateway(ports []int, ssl []*ListenConfig) UrlConfig {

	config := UrlConfig{}
	ips, _ := getIps()
	config.ListenUrls = make([]string, 0, len(ports)+len(ssl))
	config.AdvertiseUrls = make([]string, 0, (len(ports)+len(ssl))*len(ips))

	for _, p := range ports {
		config.ListenUrls = append(config.ListenUrls, fmt.Sprintf("http://0.0.0.0:%d", p))
		for _, ip := range ips {
			config.AdvertiseUrls = append(config.AdvertiseUrls, fmt.Sprintf("http://%s:%d", ip, p))
		}
	}
	certs := make(map[string]string)

	for _, sl := range ssl {
		config.ListenUrls = append(config.ListenUrls, fmt.Sprintf("https://0.0.0.0:%d", sl.Port))

		for _, ip := range ips {
			config.AdvertiseUrls = append(config.AdvertiseUrls, fmt.Sprintf("http://%s:%d", ip, sl.Port))
		}
		for _, cert := range sl.Certificate {
			certs[cert.Cert] = cert.Key
		}
	}

	config.Certificate = make([]CertConfig, 0, len(certs))
	for p, k := range certs {
		config.Certificate = append(config.Certificate, CertConfig{
			Cert: p,
			Key:  k,
		})
	}

	return config

}
func getIps() ([]string, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return nil, err
	}

	ips := make([]string, 0, len(addrs))
	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}

		}
	}
	if len(ips) == 0 {

		return nil, errors.New("not find valid ip")
	}
	return ips, nil
}

func initial(c *NConfig) {

}
