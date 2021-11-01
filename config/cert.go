package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/eolinker/eosc/log"
)

type Cert struct {
	certs map[string]*tls.Certificate
}

func NewCert(certs []*Certificate, dir string) (*Cert, error) {
	cs := make(map[string]*tls.Certificate)
	for _, c := range certs {
		if c.Key != "" && c.Cert != "" {
			cert, certificate, err := loadCert(c.Cert, c.Key, dir)
			if err != nil {
				log.Error("load certificate error: ", err, " pem is ", &c.Cert, " key is ", c.Key)
				continue
			}
			cs[certificate.Subject.CommonName] = cert
			for _, dnsName := range certificate.DNSNames {
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
			infos, err := ioutil.ReadDir(dir)
			if err != nil {
				return nil, err
			}
			certMap := make(map[string]*Certificate)
			for _, fInfo := range infos {
				name := fInfo.Name()
				if strings.HasSuffix(name, ".pem") {
					key := strings.Replace(name, ".pem", "", -1)
					if _, ok := certMap[key]; !ok {
						certMap[key] = &Certificate{}
					}
					certMap[key].Cert = name
				} else if strings.HasSuffix(name, ".key") {
					key := strings.Replace(name, ".key", "", -1)
					if _, ok := certMap[key]; !ok {
						certMap[key] = &Certificate{}
					}
					certMap[key].Key = name
				}
			}
			for _, c := range certMap {
				cert, certificate, err := loadCert(c.Cert, c.Key, dir)
				if err != nil {
					log.Error("load certificate error: ", err, " pem is ", &c.Cert, " key is ", c.Key)
					continue
				}
				cs[certificate.Subject.CommonName] = cert
				for _, dnsName := range certificate.DNSNames {
					cs[dnsName] = cert
				}
			}
		}
	}
	return &Cert{
		certs: cs,
	}, nil
}

func loadCert(pem string, key string, dir string) (*tls.Certificate, *x509.Certificate, error) {
	if !filepath.IsAbs(pem) {
		pem = fmt.Sprintf("%s/%s", strings.TrimSuffix(dir, "/"), strings.TrimPrefix(pem, "/"))
	}
	if !filepath.IsAbs(key) {
		key = fmt.Sprintf("%s/%s", strings.TrimSuffix(dir, "/"), strings.TrimPrefix(key, "/"))
	}
	cert, err := tls.LoadX509KeyPair(pem, key)
	if err != nil {
		return nil, nil, err

	}
	certificate, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, nil, err
	}
	return &cert, certificate, nil
}

func (c *Cert) GetCertificate(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
	if c.certs == nil {
		return nil, errorCertificateNotExit
	}
	certificate, has := c.Get(strings.ToLower(info.ServerName))
	if !has {
		return nil, errorCertificateNotExit
	}

	return certificate, nil
}

//Get 获取证书
func (c *Cert) Get(hostName string) (*tls.Certificate, bool) {
	if c == nil || len(c.certs) == 0 {
		return nil, true
	}
	cert, has := c.certs[hostName]
	if has {
		return cert, true
	}
	hs := strings.Split(hostName, ".")
	if len(hs) < 1 {
		return nil, false
	}

	cert, has = c.certs[fmt.Sprintf("*.%s", strings.Join(hs[1:], "."))]
	return cert, has
}
