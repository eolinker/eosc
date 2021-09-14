package eosc_args

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/eolinker/eosc"

	"github.com/eolinker/eosc/log"
)

type Config struct {
	path string
	args eosc.IUntyped
}

func NewConfig(path string) (*Config, error) {
	// 参数配置文件格式：分行获取
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("load args error: %s", err.Error())
		return nil, err
	}
	var args = map[string]string{}
	lineData := strings.Split(string(data), "\n")
	for _, d := range lineData {
		d = strings.TrimSpace(d)
		index := strings.Index(d, "=")
		if index == -1 {
			args[d] = ""
			continue
		}
		args[d[:index]] = d[index+1:]
	}
	cfg := &Config{path: path, args: eosc.NewUntyped()}
	for key, value := range args {
		cfg.Set(key, value)
	}
	return cfg, nil
}

func (c *Config) Set(name string, value string) {
	if name != "" {
		c.args.Set(EnvName(name), value)
	}
}

func (c *Config) Get(name string) (string, bool) {
	vl, has := c.args.Get(EnvName(name))
	if !has {
		return "", false
	}
	v, ok := vl.(string)
	if !ok {
		return "", false
	}
	return v, true
}

func (c *Config) GetDefault(name string, value string) string {
	vl, has := c.args.Get(EnvName(name))
	if !has {
		return value
	}
	v, ok := vl.(string)
	if !ok {
		return value
	}
	return v
}

func (c *Config) Save() error {
	builder := strings.Builder{}
	for key, value := range c.args.All() {
		v, ok := value.(string)
		if !ok {
			continue
		}
		builder.WriteString(key)
		builder.WriteString("=")
		builder.WriteString(v)
		builder.WriteString("\n")
	}
	f, err := os.OpenFile(c.path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	// offset
	_, err = f.Write([]byte(builder.String()))
	if err != nil {
		return err
	}
	log.Info("write args file succeed!")

	return nil
}
