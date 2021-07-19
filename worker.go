package eosc

type IWorker interface {
	Id() string
	Start() error
	Reset(conf interface{}, workers map[RequireId]interface{}) error
	Stop() error
	CheckSkill(skill string) bool
}

type tWorker struct {
	worker IWorker
}

func newTWorker(worker IWorker) *tWorker {
	return &tWorker{worker: worker}
}
