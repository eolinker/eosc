package worker

import "github.com/eolinker/eosc"

type tWorker struct {
	worker eosc.IWorker
}

func newTWorker(worker eosc.IWorker) *tWorker {
	return &tWorker{worker: worker}
}
