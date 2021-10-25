package config

import (
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"

	"github.com/eolinker/eosc/env"
)

type Config struct {
	Listen         []int           `json:"listen" yaml:"listen"`
	SSL            *SSLConfig      `json:"ssl" yaml:"ssl"`
	Admin          *AdminConfig    `json:"admin" yaml:"admin"`
	CertificateDir *CertificateDir `json:"certificate" yaml:"certificate"`
}

type SSLConfig struct {
	Listen *ListenConfig `json:"listen"`
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

type Certificate struct {
	Key string `json:"key" yaml:"key"`
	Pem string `json:"pem" yaml:"pem"`
}

type CertificateDir struct {
	Dir string `json:"dir" yaml:"dir"`
}

var defaultPath = "/etc/%s/config.yml"

func GetConfig() (*Config, error) {
	path, has := env.GetEnv("config")
	if !has {
		path = fmt.Sprintf(defaultPath, env.AppName())
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	return &cfg, err
}
