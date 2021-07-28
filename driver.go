package eosc

import "reflect"

type ExtendInfo struct {
	ID      string `json:"id"`
	Group   string `json:"group"`
	Project string `json:"project"`
	Name    string `json:"name"`
}

type DriverInfo struct {
	Id         string            `json:"id"`
	Name       string            `json:"name"`
	Label      string            `json:"label"`
	Desc       string            `json:"desc"`
	Profession string            `json:"profession"`
	Params     map[string]string `json:"params"`
}
type DriverDetail struct {
	DriverInfo
	Extends ExtendInfo `json:"extends"`
}
type IProfessionDriverFactory interface {
	ExtendInfo() ExtendInfo
	Create(profession string, name string, label string, desc string, params map[string]string) (IProfessionDriver, error)
}

type IProfessionDriverConfigCheck interface {
	Check(id, name string, v interface{}, workers map[RequireId]interface{}) error
}

type IProfessionDriver interface {
	ConfigType() reflect.Type
	Create(id, name string, v interface{}, workers map[RequireId]interface{}) (IWorker, error)
}

type IProfessionDriverInfo interface {
	IProfessionDriver
	ExtendInfo() ExtendInfo
	DriverInfo() DriverInfo
}
