package eosc_args

import (
	"bytes"
	"io"
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

func NewConfig(path string) *Config {
	return &Config{path: path, args: eosc.NewUntyped()}
}

func (c *Config) ReadFile(paths ...string) {
	for _, path := range paths {
		// 参数配置文件格式：分行获取
		data, err := ioutil.ReadFile(path)
		if err != nil {
			continue
		}
		buf := bytes.NewBuffer(data)
		for {
			line, err := buf.ReadString('\n')
			if err != nil && err != io.EOF {
				break
			}
			line = strings.TrimSpace(line)
			if len(line) > 0 {
				index := strings.Index(line, "=")
				if index == -1 {
					continue
				}
				c.Set(line[:index], line[index+1:])
			}
			if err != nil {
				return
			}
		}
	}
}

func (c *Config) Set(name string, value string) {
	if name != "" {
		c.args.Set(name, value)
	}
}

func (c *Config) Get(name string) (string, bool) {
	vl, has := c.args.Get(name)
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
	vl, has := c.args.Get(name)
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
