package process_worker

import "github.com/eolinker/eosc"

var _ IWorker = (*Worker)(nil)

type IWorker interface {
	eosc.IWorker
	//GetProfession() *Profession
}
type Worker struct {
	eosc.IWorker
	id         string
	Profession string
	Name       string
	Driver     string
	body       []byte

	//profession *Profession
	driver eosc.IExtenderDriver
}

func NewWorker(id, professionName, name, driverName string, body []byte, target eosc.IWorker, profession *Profession, driver eosc.IExtenderDriver) *Worker {
	return &Worker{
		IWorker:    target,
		id:         id,
		Profession: professionName,
		Name:       name,
		Driver:     driverName,
		//profession: profession,
		body:   body,
		driver: driver,
	}
}

//func (w *Worker) GetProfession() *Profession {
//	return w.profession
//}

func (w *Worker) Id() string {
	return w.id
}
