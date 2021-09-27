package eosc

import (
	"reflect"
)

type IProfessionDriverFactory interface {
	Create(profession string, name string, label string, desc string, params map[string]string) (IProfessionDriver, error)
}
type IProfessionDriverCheckConfig interface {
	Check(v interface{}, workers map[RequireId]interface{}) error
}
type IProfessionDriver interface {
	ConfigType() reflect.Type
	Create(id, name string, v interface{}, workers map[RequireId]interface{}) (IWorker, error)
}
