package config

import (
	"fmt"
	"strings"
)

type OConfig struct {
	Listen         []int           `json:"listen" yaml:"listen"`
	SSL            *SSLConfig      `json:"ssl" yaml:"ssl"`
	Admin          *AdminConfig    `json:"admin" yaml:"admin"`
	CertificateDir *CertificateDir `json:"certificate" yaml:"certificate"`
}
type SSLConfig struct {
	Listen []*ListenConfig `json:"listen"`
}

type ListenConfig struct {
	Port        int            `json:"port" yaml:"port"`
	Certificate []*Certificate `json:"certificate" yaml:"certificate"`
}

type AdminConfig struct {
	Scheme      string       `json:"scheme" yaml:"scheme"`
	Listen      int          `json:"listen" yaml:"listen"`
	IP          string       `json:"ip" yaml:"ip"`
	Certificate *Certificate `json:"certificate" yaml:"certificate"`
}

func fromAdmin(admin *AdminConfig) (UrlConfig, UrlConfig) {

	scheme := strings.ToLower(admin.Scheme)
	if scheme != "https" {
		scheme = "http"
	}

	port := admin.Listen
	if port == 0 {
		port = 9400
	}
	ip := admin.IP
	certificate := admin.Certificate
	return createDefault(scheme, port+1, ip, certificate), createDefault(scheme, port, ip, certificate)
}
func createDefault(scheme string, port int, ip string, certificate *Certificate) UrlConfig {
	config := UrlConfig{}
	if len(ip) == 0 {
		ip = "0.0.0.0"
	}
	ssl := (scheme == "https") && (certificate != nil)
	if !ssl {
		scheme = "http"
	}
	config.ListenUrls = []string{fmt.Sprintf("%s://%s:%d", scheme, ip, port)}

	if ip != "0.0.0.0" {
		config.AdvertiseUrls = []string{fmt.Sprintf("%s://%s:%d", scheme, ip, port)}
	}
	if ssl {
		config.Certificate = make([]CertConfig, 0, 1)
		config.Certificate = append(config.Certificate, CertConfig{
			Cert: certificate.Cert,
			Key:  certificate.Key,
		})
	}
	return config
}
func toGateway(ports []int, ssl []*ListenConfig) ListenUrl {

	config := ListenUrl{}

	config.ListenUrls = make([]string, 0, len(ports)+len(ssl))

	for _, p := range ports {
		config.ListenUrls = append(config.ListenUrls, fmt.Sprintf("http://0.0.0.0:%d", p))
	}
	certs := make(map[string]string)

	for _, sl := range ssl {
		config.ListenUrls = append(config.ListenUrls, fmt.Sprintf("https://0.0.0.0:%d", sl.Port))
		for _, cert := range sl.Certificate {
			certs[cert.Cert] = cert.Key
		}
	}
	//
	//config.Certificate = make([]CertConfig, 0, len(certs))
	//for p, k := range certs {
	//	config.Certificate = append(config.Certificate, CertConfig{
	//		Cert: p,
	//		Key:  k,
	//	})
	//}

	return config

}
