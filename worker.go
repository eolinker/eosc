package eosc

import (
	"time"
)

type RequireId string
type IWorker interface {
	Id() string
	Start() error
	Reset(conf interface{}, workers map[RequireId]interface{}) error
	Stop() error
	CheckSkill(skill string) bool
}
type IWorkers interface {
	Get(id string) (IWorker, bool)
}

type TWorker struct {
	Id         string      `json:"id,omitempty" yaml:"id"`
	Name       string      `json:"name,omitempty" yaml:"name"`
	Driver     string      `json:"driver,omitempty" yaml:"driver"`
	Profession string      `json:"profession,omitempty" yaml:"profession"`
	Create     time.Time   `json:"create" yaml:"create"`
	Update     time.Time   `json:"update" yaml:"update"`
	Data       interface{} `json:"data,omitempty" yaml:"data"`
}

//type IWorkersData interface {
//	GetWork(id string) (TWorker, error)
//	GetList(profession string) ([]TWorker, error)
//	Delete(id string) (TWorker, error)
//	Set(profession, name, driver string, data []byte) (TWorker, error)
//}

type IWorkerResources interface {
	Ports() []int
}
