package config

import (
	"encoding/json"
	"errors"

	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/log/dlog"
)

const driverName = "httplog"
const driverTitle = "http日志"
const defaultLineFormatter = "json"

type httpLogConfigDriver struct {
	*dlog.FullFieldsDriver
}

func NewHttpLogConfigDriver() dlog.ConfigDriver {
	return &httpLogConfigDriver{
		FullFieldsDriver: dlog.NewFullFieldsDriver(fields),
	}
}
func (c *httpLogConfigDriver) Format(v string) (interface{}, error) {
	configConfig := &ConfigEncode{}
	err := json.Unmarshal([]byte(v), configConfig)
	if err != nil {
		return nil, err
	}
	return configConfig, nil
}
func (c *httpLogConfigDriver) Name() string {
	return driverName
}

func (c *httpLogConfigDriver) Title() string {
	return driverTitle
}

func (c *httpLogConfigDriver) Decode(v string) (interface{}, error) {
	configConfig := &ConfigEncode{}
	err := json.Unmarshal([]byte(v), configConfig)
	if err != nil {
		return nil, err
	}
	return ToConfig(configConfig)

}

func (c *httpLogConfigDriver) Encode(v interface{}) (string, error) {
	if e, ok := v.(encoder); ok {
		return e.encode()
	}
	return "", errors.New("unknown config type")
}

func ToConfig(v interface{}) (*Config, error) {
	if v == nil {
		return nil, errors.New("config is nil")
	}
	if c, ok := v.(*Config); ok {
		return c, nil
	}
	if c, ok := v.(*ConfigEncode); ok {
		level, err := log.ParseLevel(c.Level)
		if err != nil {
			level = log.InfoLevel
		}

		config := &Config{
			Method:       c.Method,
			Url:          c.Url,
			Headers:      toHeader(c.Headers),
			Level:        level,
			HandlerCount: 5, // 默认值， 以后可能会改成配置
		}

		return config, nil
	}
	return nil, errors.New("unknown type")
}
