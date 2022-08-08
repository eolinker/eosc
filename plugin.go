package eosc

import "reflect"

type IPluginReset interface {
	Reset(originVal reflect.Value, targetVal reflect.Value, params map[string]string, configTypes map[string]reflect.Type) ([]string, error)
}
