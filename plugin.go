package eosc

import "reflect"

type IPluginReset interface {
	Reset(originVal reflect.Value, targetVal reflect.Value, params map[string]string, configType *ConfigType) ([]string, error)
}
