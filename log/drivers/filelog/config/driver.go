package config

import (
	"encoding/json"
	"errors"
	"github.com/eolinker/goku-standard/common/log"
	"github.com/eolinker/goku-standard/common/log/dlog"
)

const driverName  = "filelog"
const driverTitle  = "文件日志"
const defaultLineFormatter  = "line"
type FileConfigDriver struct {
*dlog.FullFieldsDriver
}

func (c *FileConfigDriver) Format(v string) (interface{}, error) {
	configConfig := &ConfigEncode{}
	err := json.Unmarshal([]byte(v), configConfig)
	if err!= nil{
		return  nil,err
	}
	return configConfig,nil
}

func NewFileConfigDriver() dlog.ConfigDriver {
	return &FileConfigDriver{
		FullFieldsDriver: dlog.NewFullFieldsDriver(fields),
	}
}

func (c *FileConfigDriver) Name() string {
	return driverName
}

func (c *FileConfigDriver) Title() string {
	return driverTitle
}

func (c *FileConfigDriver) Decode(v string) (interface{}, error) {

	configConfig := &ConfigEncode{}
	err := json.Unmarshal([]byte(v), configConfig)
	if err!= nil{
		return  nil,err
	}

	return ToConfig(configConfig)
}

func (c *FileConfigDriver) Encode(v interface{}) (string, error) {

	if e,ok:=v.(encoder);ok{
		return e.encode()
	}
	return "", errors.New("unknown config type")

}



func ToConfig(v interface{}) (*Config,error) {
	if v== nil{
		return nil,errors.New("config is nil")
	}
	if c,ok:=v.(*Config);ok{
		return c,nil
	}
	if c,ok:=v.(*ConfigEncode);ok{
		period ,err:= ParsePeriod(c.Period)
		if err!=nil{
			return nil,err
		}
		level, err := log.ParseLevel(c.Level)
		if err!= nil{
			level = log.InfoLevel
		}

		config:= Config{
			Dir:    c.Dir,
			File:   c.File,
			Expire: c.Expire,
			Period:period,
			Level:level,

		}

		return &config,nil
	}
	return nil,errors.New("unknown type")
}