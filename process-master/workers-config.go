package process_master

import (
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/process-master/workers"
)

type WorkerConfigs struct {
	workers.ITypedWorkers
}

func NewWorkerConfigs() *WorkerConfigs {
	return &WorkerConfigs{ITypedWorkers: workers.NewTypedWorkers()}
}

func (w *WorkerConfigs) export() []*eosc.WorkerConfig {
	values := w.ITypedWorkers.All()

	wds := make([]*eosc.WorkerConfig, len(values))

	for i, v := range values {

		wds[i] = v.WorkerConfig
	}
	return wds
}

func (w *WorkerConfigs) reset(vs []*eosc.WorkerConfig) {
	w.ITypedWorkers.Reset(vs)
}
