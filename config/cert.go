package config

import (
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tjfoc/gmsm/gmtls"

	"github.com/eolinker/eosc/log"
)

var (
	errorCertificateNotExit = errors.New("not exist cert")
)

type Cert[T any] struct {
	certs map[string]*T
}

func NewCert[T any](certs map[string]*T) *Cert[T] {
	return &Cert[T]{certs: certs}
}

func LoadCert(certs []CertConfig, dir string) (*Cert[tls.Certificate], error) {
	cs := make(map[string]*tls.Certificate)
	for _, c := range certs {
		if c.Key != "" && c.Cert != "" {
			cert, err := loadCert(c.Cert, c.Key, dir)
			if err != nil {
				log.Error("load certificate error: ", err, " pem is ", &c.Cert, " key is ", c.Key)
				continue
			}
			cs[cert.Leaf.Subject.CommonName] = cert
			for _, dnsName := range cert.Leaf.DNSNames {
				cs[dnsName] = cert
			}
		}
	}
	if len(cs) < 1 {
		info, err := os.Stat(dir)
		if err != nil {
			return nil, err
		}
		if info.IsDir() {
			infos, err := os.ReadDir(dir)
			if err != nil {
				return nil, err
			}
			certMap := make(map[string]*CertConfig)
			for _, fInfo := range infos {
				name := fInfo.Name()
				if strings.HasSuffix(name, ".pem") {
					key := strings.Replace(name, ".pem", "", -1)
					if _, ok := certMap[key]; !ok {
						certMap[key] = &CertConfig{}
					}
					certMap[key].Cert = name
				} else if strings.HasSuffix(name, ".key") {
					key := strings.Replace(name, ".key", "", -1)
					if _, ok := certMap[key]; !ok {
						certMap[key] = &CertConfig{}
					}
					certMap[key].Key = name
				}
			}
			for _, c := range certMap {
				cert, err := loadCert(c.Cert, c.Key, dir)
				if err != nil {
					log.Error("load certificate error: ", err, " pem is ", &c.Cert, " key is ", c.Key)
					continue
				}
				cs[cert.Leaf.Subject.CommonName] = cert
				for _, dnsName := range cert.Leaf.DNSNames {
					cs[dnsName] = cert
				}
			}
		}
	}
	return NewCert[tls.Certificate](cs), nil
}

func loadCert(pem string, key string, dir string) (*tls.Certificate, error) {
	if !filepath.IsAbs(pem) {
		pem = fmt.Sprintf("%s/%s", strings.TrimSuffix(dir, "/"), strings.TrimPrefix(pem, "/"))
	}
	if !filepath.IsAbs(key) {
		key = fmt.Sprintf("%s/%s", strings.TrimSuffix(dir, "/"), strings.TrimPrefix(key, "/"))
	}
	cert, err := tls.LoadX509KeyPair(pem, key)
	return &cert, err
}

func (c *Cert[T]) GetCertificate(clientHello interface{}) (*T, error) {
	if c.certs == nil {
		return nil, errorCertificateNotExit
	}
	name := ""
	switch t := clientHello.(type) {
	case *tls.ClientHelloInfo:
		name = strings.ToLower(t.ServerName)
	case *gmtls.ClientHelloInfo:
		name = strings.ToLower(t.ServerName)
	default:
		return nil, fmt.Errorf("unsupported type %T for GetCertificate", clientHello)
	}
	if len(c.certs) == 1 {
		// There's only one choice, so no point doing any work.
		for _, cert := range c.certs {
			return cert, nil
		}
	}
	if cert, ok := c.certs[name]; ok {
		return cert, nil
	}
	if len(name) > 0 {
		labels := strings.Split(name, ".")
		labels[0] = "*"
		wildcardName := strings.Join(labels, ".")
		if cert, ok := c.certs[wildcardName]; ok {
			return cert, nil
		}
	}

	return nil, errorCertificateNotExit
}

func GetCertificateFunc(cert *Cert[tls.Certificate]) func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
		return cert.GetCertificate(info)
	}
}
