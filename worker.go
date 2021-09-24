package eosc

type IWorker interface {
	Id() string
	Start() error
	Reset(conf interface{}, workers map[RequireId]interface{}) error
	Stop() error
	CheckSkill(skill string) bool
}
type IWorkers interface {
	Set(id string, w IWorker)
	Del(id string) (IWorker, bool)
	Get(id string) (IWorker, bool)
}

type TWorker map[string]interface{}
type IWorkersData interface {
	GetWork(id string) (TWorker, error)
	GetList(profession string) ([]TWorker, error)
	Delete(id string) (*WorkerInfo, error)
	Set(profession, name, driver string, data []byte) error
}
