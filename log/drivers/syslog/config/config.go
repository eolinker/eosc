package config

import (
	"encoding/json"

	"github.com/eolinker/eosc/log"
)

type ConfigEncode struct {
	Network string `json:"network"`
	URL     string `json:"url"`
	Level   string `json:"level"`
}

func (c *ConfigEncode) encode() (string, error) {

	data, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(data), err
}

type Config struct {
	Network string
	RAddr   string
	Level   log.Level
}
