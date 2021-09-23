package worker

import "github.com/eolinker/eosc"

type IWorker interface {
	Id() string
	Start() error
	Reset(conf interface{}, workers map[eosc.RequireId]interface{}) error
	Stop() error
	CheckSkill(skill string) bool
}

type tWorker struct {
	worker IWorker
}

func newTWorker(worker IWorker) *tWorker {
	return &tWorker{worker: worker}
}
