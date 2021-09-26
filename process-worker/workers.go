package process_worker

import (
	"encoding/json"
	"errors"
	"io"
	"reflect"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/utils"
	"github.com/golang/protobuf/proto"
)

var _ eosc.IWorkers = (*WorkerManager)(nil)

func ReadWorkers(r io.Reader) []*eosc.WorkerData {
	frame, err := utils.ReadFrame(r)
	if err != nil {
		log.Warn("read  workers frame :", err)

		return nil
	}

	wd := new(eosc.WorkersData)
	if e := proto.Unmarshal(frame, wd); e != nil {
		log.Warn("unmarshal workers data :", e)
		return nil
	}
	return wd.Data
}

type WorkerManager struct {
	professions IProfessions
	data        ITypedWorkers
}

func (wm *WorkerManager) Get(id string) (eosc.IWorker, bool) {
	return wm.data.Get(id)
}

func NewWorkerManager(professions IProfessions) *WorkerManager {
	return &WorkerManager{professions: professions, data: NewTypedWorkers()}
}

func (wm *WorkerManager) Init(wdl []*eosc.WorkerData) error {

	ps := wm.professions.Sort()

	pm := make(map[string][]*eosc.WorkerData)
	for _, wd := range wdl {
		pm[wd.Profession] = append(pm[wd.Profession], wd)
	}

	for _, p := range ps {
		for _, wd := range pm[p.Name] {
			if err := wm.Set(wd); err != nil {
				return err
			}
		}
	}
	return nil
}

func (wm *WorkerManager) Set(wd *eosc.WorkerData) error {
	p, has := wm.professions.Get(wd.Profession)
	if !has {
		return errors.New("profession not exist")
	}
	driver, has := p.GetDriver(wd.Driver)
	if !has {
		return errors.New("driver not exist")
	}
	configType := driver.ConfigType()
	conf := reflect.New(configType).Interface()

	err := json.Unmarshal(wd.Body, conf)
	if err != nil {
		return err
	}
	requires, err := eosc.CheckConfig(conf, wm)
	if err != nil {
		return err
	}
	if dc, ok := driver.(eosc.IProfessionDriverCheckConfig); ok {
		if err := dc.Check(conf, requires); err != nil {
			return err
		}
	}
	// if update
	o, has := wm.data.Get(wd.Id)
	if has {
		return o.Reset(conf, requires)
	}
	// create
	w, err := driver.Create(wd.Id, wd.Name, conf, requires)
	if err != nil {
		return err
	}
	// start
	e := w.Start()
	if e != nil {
		return e
	}
	// cache
	wm.data.Set(wd.Id, NewWorker(wd, w, p))
	return nil
}
