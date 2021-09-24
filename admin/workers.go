package admin

import "github.com/eolinker/eosc"

type TWorker map[string]interface{}
type IWorkers interface {
	GetWork(id string) (TWorker, error)
	GetList(profession string) ([]TWorker, error)
	Delete(id string) (*eosc.WorkerInfo, error)
	Set(profession, name, driver string, data []byte) error
}
