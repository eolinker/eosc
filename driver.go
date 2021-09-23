package eosc

import (
	"reflect"

	"github.com/eolinker/eosc/process-worker/worker"
)

type DriverInfo struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Label      string `json:"label"`
	Desc       string `json:"desc"`
	Profession string `json:"profession"`
}
type DriverDetail struct {
	DriverInfo
	Group   string            `json:"group"`
	Project string            `json:"project"`
	Param   map[string]string `json:"param"`
	//	TODO: 待加字段：预期版本、实际版本
}

type IProfessionDriverFactory interface {
	Create(profession string, name string, label string, desc string, params map[string]string) (IProfessionDriver, error)
}
type IProfessionDriverCheckConfig interface {
	Check(v interface{}, workers map[RequireId]interface{}) error
}
type IProfessionDriver interface {
	ConfigType() reflect.Type
	Create(id, name string, v interface{}, workers map[RequireId]interface{}) (worker.IWorker, error)
}

type IProfessionDriverInfo interface {
	IProfessionDriver
	DriverInfo() DriverInfo
}
