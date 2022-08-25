package eosc

import (
	"reflect"
)

type IExtenderDriverFactory interface {
	Render() interface{}
	Create(profession string, name string, label string, desc string, params map[string]interface{}) (IExtenderDriver, error)
}

type IExtenderConfigChecker interface {
	Check(v interface{}, workers map[RequireId]IWorker) error
}

type IExtenderDriver interface {
	ConfigType() reflect.Type
	Create(id, name string, v interface{}, workers map[RequireId]IWorker) (IWorker, error)
}

type ISetting interface {
	Render() interface{}
	ConfigType() reflect.Type
	Set(conf interface{}) error
	Get() interface{}
	ReadOnly() bool
}

type ISettings interface {
	GetDriver(name string) (ISetting, bool)
	Set(name string, org []byte, variable IVariable) (format interface{}, err error)
	Update(name string, variable IVariable) (err error)
	CheckVariable(name string, variable IVariable) (err error)
	GetConfig(name string) interface{}
}
