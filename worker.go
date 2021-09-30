package eosc

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

type TWorker map[string]interface{}
type IWorkersData interface {
	GetWork(id string) (TWorker, error)
	GetList(profession string) ([]TWorker, error)
	Delete(id string) (TWorker, error)
	Set(profession, name, driver string, data []byte) (TWorker, error)
}

type IWorkerResources interface {
	Ports() []int
}
