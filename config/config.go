package config

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/utils"

	"google.golang.org/protobuf/proto"

	"github.com/ghodss/yaml"

	"github.com/eolinker/eosc/env"
)

var (
	errorCertificateNotExit = errors.New("not exist cert")
)

type Config struct {
	Listen         []int           `json:"listen" yaml:"listen"`
	SSL            *SSLConfig      `json:"ssl" yaml:"ssl"`
	Admin          *AdminConfig    `json:"admin" yaml:"admin"`
	CertificateDir *CertificateDir `json:"certificate" yaml:"certificate"`
}

func (c *Config) Ports() []int {
	portLen := len(c.Listen)
	if c.SSL != nil {
		portLen += len(c.SSL.Listen)
	}
	ports := make([]int, 0, portLen)
	ports = append(ports, c.Listen...)
	for _, p := range c.SSL.Listen {
		ports = append(ports, p.Port)
	}
	return ports
}

func (c *Config) Encode(startIndex int) ([]byte, []*os.File, error) {
	data, err := c.encode()
	if err != nil {
		return nil, nil, err
	}
	return utils.EncodeFrame(data), nil, nil
}

func (c *Config) encode() ([]byte, error) {
	portNum := len(c.Listen)
	if c.SSL != nil {
		portNum += len(c.SSL.Listen)
	}
	cfg := &ListensMsg{Listens: make([]*ListenMsg, 0, portNum), Dir: c.CertificateDir.Dir}
	for _, p := range c.Listen {
		cfg.Listens = append(cfg.Listens, &ListenMsg{
			Port:        int32(p),
			Scheme:      "http",
			Certificate: nil,
		})
	}
	if c.SSL != nil {
		for _, info := range c.SSL.Listen {
			cfg.Listens = append(cfg.Listens, &ListenMsg{
				Port:        int32(info.Port),
				Scheme:      "https",
				Certificate: info.Certificate,
			})
		}
	}

	data, err := proto.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	return data, nil
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

type CertificateDir struct {
	Dir string `json:"dir" yaml:"dir"`
}

var defaultPath = "/etc/%s/config.yml"

func GetConfig() (*Config, error) {
	path := env.ConfigPath()
	if path == "" {
		path = fmt.Sprintf(defaultPath, env.AppName())
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if cfg.Admin == nil {
		cfg.Admin = &AdminConfig{Listen: 9400, IP: "", Scheme: "http"}
	}
	if cfg.CertificateDir == nil {
		cfg.CertificateDir = &CertificateDir{
			Dir: fmt.Sprintf("/etc/%s/cert", env.AppName()),
		}
	}
	err = cfg.checkPort()
	if err != nil {
		return nil, err
	}
	return &cfg, err
}

func (c *Config) checkPort() error {
	usedPorts := map[int]bool{c.Admin.Listen: true}
	for _, p := range c.Listen {
		if _, ok := usedPorts[p]; ok {
			return errors.New(fmt.Sprintf("repeated port %d in listen config, please check config", p))
		}
	}
	if c.SSL != nil {
		for _, cfg := range c.SSL.Listen {
			if _, ok := usedPorts[cfg.Port]; ok {
				return errors.New(fmt.Sprintf("repeated port %d in ssl listen config, please check config", cfg.Port))
			}
		}
	}

	return nil
}

func ReadHttpTrafficConfig(r io.Reader) *ListensMsg {
	frame, err := utils.ReadFrame(r)
	if err != nil {
		log.Warn("read  workerIds frame :", err)

		return nil
	}

	msg := new(ListensMsg)
	if e := proto.Unmarshal(frame, msg); e != nil {
		log.Warn("unmarshal workerIds data :", e)
		return nil
	}

	return msg
}
