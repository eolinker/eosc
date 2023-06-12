package env

import (
	"bytes"
	"go.etcd.io/etcd/client/pkg/v3/fileutil"
	"go.uber.org/zap"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/eolinker/eosc"

	"github.com/eolinker/eosc/log"
)

type Config struct {
	path string
	args eosc.Untyped[string, string]
}

func NewConfig(path string) *Config {
	return &Config{path: path, args: eosc.BuildUntyped[string, string]()}
}

func (c *Config) ReadFile(paths ...string) {
	for _, path := range paths {
		// 参数配置文件格式：分行获取

		_ = fileutil.CreateDirAll(zap.NewNop(), filepath.Dir(path))

		data, err := os.ReadFile(path)
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

	return c.args.Get(name)

}

func (c *Config) GetDefault(name string, value string) string {
	v, has := c.args.Get(name)
	if !has {
		return value
	}
	return v
}

func (c *Config) Save() error {
	builder := strings.Builder{}
	for key, value := range c.args.All() {
		builder.WriteString(key)
		builder.WriteString("=")
		builder.WriteString(value)
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
