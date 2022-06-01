package eosc

import (
	"reflect"
)

type IExtenderDriverFactory interface {
	Render() interface{}
	Create(profession string, name string, label string, desc string, params map[string]interface{}) (IExtenderDriver, error)
}

type IExtenderConfigChecker interface {
	Check(v interface{}, workers map[RequireId]interface{}) error
}

type IExtenderDriver interface {
	ConfigType() reflect.Type
	Create(id, name string, v interface{}, workers map[RequireId]interface{}) (IWorker, error)
}
