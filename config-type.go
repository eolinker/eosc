package eosc

import (
	"reflect"
)

func NewConfigType(alias map[string]string, cfgType map[string]reflect.Type) *ConfigType {
	return &ConfigType{alias: alias, cfgType: cfgType}
}

type ConfigType struct {
	alias   map[string]string
	cfgType map[string]reflect.Type
}

func (c *ConfigType) Get(id string) (reflect.Type, bool) {
	return c.get(id)
}

func (c *ConfigType) get(id string) (reflect.Type, bool) {
	v, ok := c.cfgType[id]
	return v, ok
}

func (c *ConfigType) GetByAlias(alias string) (reflect.Type, bool) {
	id, ok := c.alias[alias]
	if !ok {
		return nil, false
	}
	return c.get(id)
}
