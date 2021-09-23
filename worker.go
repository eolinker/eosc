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
