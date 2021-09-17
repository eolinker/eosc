package admin

import "github.com/eolinker/eosc"

type TWorker map[string]interface{}
type IWorkers interface {
	GetWork(id string) (TWorker, error)
	GetList(profession string) ([]TWorker, error)
	CheckerSkill(id string, skill string) (bool, error)
	Delete(id string) (*WorkerInfo, bool)
	Set(profession, name, driver string, data eosc.IData) error
}
