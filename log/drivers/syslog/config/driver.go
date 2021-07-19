package config

import (
	"encoding/json"
	"errors"

	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/log/dlog"
)

const driverName = "syslog"
const driverTitle = "Syslog"
const defaultLineFormatter = "json"

type encoder interface {
	encode() (string, error)
}

func NewSysLogConfigDriver() dlog.ConfigDriver {
	return &SysLogConfigDriver{
		FullFieldsDriver: dlog.NewFullFieldsDriver(fields),
	}
}

func ToConfig(v interface{}) (*Config, error) {
	if v == nil {
		return nil, errors.New("config is nil")
	}
	if c, ok := v.(*Config); ok {
		return c, nil
	}
	return nil, errors.New("unknown type")
}

type SysLogConfigDriver struct {
	*dlog.FullFieldsDriver
}

func (c *SysLogConfigDriver) Name() string {
	return driverName
}

func (c *SysLogConfigDriver) Title() string {
	return driverTitle
}
func (c *SysLogConfigDriver) Format(v string) (interface{}, error) {
	configConfig := &ConfigEncode{}
	err := json.Unmarshal([]byte(v), configConfig)
	if err != nil {
		return nil, err
	}
	return configConfig, nil
}

func (c *SysLogConfigDriver) Decode(v string) (interface{}, error) {
	configConfig := &ConfigEncode{}
	err := json.Unmarshal([]byte(v), configConfig)
	if err != nil {
		return nil, err
	}

	level, err := log.ParseLevel(configConfig.Level)
	if err != nil {
		level = log.InfoLevel
	}

	config := Config{
		RAddr:   configConfig.URL,
		Network: configConfig.Network,
		Level:   level,
	}

	return config, nil
}

func (c *SysLogConfigDriver) Encode(v interface{}) (string, error) {
	if e, ok := v.(encoder); ok {
		return e.encode()
	}
	return "", errors.New("unknown config type")
}

//
//func (c *SysLogConfigDriver) ConfigFields() []dlog.Field {
//	return fields
//}
