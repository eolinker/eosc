package eosc

import (
	"reflect"
)

type IExtenderDriverFactory interface {
	Create(profession string, name string, label string, desc string, params map[string]interface{}) (IExtenderDriver, error)
}

type IExtenderConfigChecker interface {
	Check(v interface{}, workers map[RequireId]IWorker) error
}

type IExtenderDriver interface {
	Render() interface{}
	ConfigType() reflect.Type
	Create(id, name string, v interface{}, workers map[RequireId]IWorker) (IWorker, error)
}
