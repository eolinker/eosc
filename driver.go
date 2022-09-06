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
type SettingMode int

const (
	SettingModeReadonly SettingMode = iota
	SettingModeSingleton
	SettingModeBatch
)

type ISetting interface {
	ConfigType() reflect.Type
	Set(conf ...interface{}) (update []*WorkerConfig, delete []string, err error)
	Get() interface{}
	Mode() SettingMode
}

type ISettings interface {
	GetDriver(name string) (ISetting, bool)
	SettingWorker(id string, config []byte, variable IVariable) error
	Set(name string, org []byte, variable IVariable) (format interface{}, update []*WorkerConfig, delete []string, err error)
	Update(id string, variable IVariable) (err error)
	CheckVariable(name string, variable IVariable) (err error)
	GetConfig(name string) interface{}
}
