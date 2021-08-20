package config

import (
	"encoding/json"
	"github.com/eolinker/goku-standard/common/log"
)

type encoder interface {
	encode()(string,error)
}
type ConfigEncode struct {

	Dir string `json:"dir"`
	File string `json:"file"`
	Level string `json:"level"`
	Period string `json:"period"`
	Expire int `json:"expire"`
}

func (c * ConfigEncode)encode()(string,error)  {
	data, err := json.Marshal(c)
	if err!=nil{
		return  "",err
	}
	return string(data),err
}

type Config struct {
	Dir           string
	File          string
	Expire        int
	Period        LogPeriod
	Level         log.Level
}

func (c *Config) encode() (string,error) {
	en:= &ConfigEncode{
		Dir:c.Dir,
		File:c.File,
		Level:c.Level.String(),
		Period:c.Period.String(),
		Expire:c.Expire,
	}
	return en.encode()
}