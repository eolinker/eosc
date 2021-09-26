package process_worker

import "github.com/eolinker/eosc"

var _ IWorker = (*Worker)(nil)

type IWorker interface {
	eosc.IWorker
	Profession() *Profession
}
type Worker struct {
	*eosc.WorkerData
	target     eosc.IWorker
	profession *Profession
	driver     eosc.IProfessionDriver
}

func NewWorker(workerData *eosc.WorkerData, target eosc.IWorker, profession *Profession) *Worker {
	return &Worker{WorkerData: workerData, target: target, profession: profession}
}

func (w *Worker) Profession() *Profession {
	return w.profession
}

func (w *Worker) Id() string {
	return w.WorkerData.Id
}

func (w *Worker) Start() error {
	return w.target.Start()
}

func (w *Worker) Reset(conf interface{}, workers map[eosc.RequireId]interface{}) error {
	return w.target.Start()
}

func (w *Worker) Stop() error {
	return w.target.Stop()
}

func (w *Worker) CheckSkill(skill string) bool {
	return w.target.CheckSkill(skill)
}
