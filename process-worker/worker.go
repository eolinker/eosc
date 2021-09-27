package process_worker

import "github.com/eolinker/eosc"

var _ IWorker = (*Worker)(nil)

type IWorker interface {
	eosc.IWorker
	GetProfession() *Profession
}
type Worker struct {
	id         string
	Profession string
	Name       string
	Driver     string
	body       []byte
	target     eosc.IWorker
	profession *Profession
	driver     eosc.IProfessionDriver
}

func NewWorker(id, professionName, name, driverName string, body []byte, target eosc.IWorker, profession *Profession, driver eosc.IProfessionDriver) *Worker {
	return &Worker{
		id:         id,
		Profession: professionName,
		Name:       name,
		Driver:     driverName,
		target:     target,
		profession: profession,
		body:       body,
		driver:     driver,
	}
}

func (w *Worker) GetProfession() *Profession {
	return w.profession
}

func (w *Worker) Id() string {
	return w.id
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
