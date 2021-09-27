package process_master

import (
	"bytes"
	"os"
	"sync"

	"github.com/eolinker/eosc/log"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/service"
)

var _ service.WorkerServiceClient = (*WorkerController)(nil)

type WorkerController struct {
	locker sync.Mutex
	dms    []eosc.IDataMarshaller

	current *WorkerProcess

	expireWorkers []*WorkerProcess

	isStop bool
}

func NewWorkerController(dms ...eosc.IDataMarshaller) *WorkerController {
	return &WorkerController{dms: dms}
}
func (wc *WorkerController) Stop() {
	wc.locker.Lock()
	defer wc.locker.Unlock()

	if wc.isStop {
		return
	}
	wc.isStop = true
	if wc.current != nil {
		wc.current.Close()
		wc.expireWorkers = append(wc.expireWorkers, wc.current)
		wc.current = nil
	}

}
func (wc *WorkerController) check(w *WorkerProcess) {
	err := w.cmd.Wait()
	if err != nil {
		log.Warn("worker exit:", err)

	}
	wc.locker.Lock()
	defer wc.locker.Unlock()
	if wc.current == w {
		err := wc.new()
		if err != nil {
			log.Error("worker create:", err)
		}
	} else {

		for i, v := range wc.expireWorkers {
			if v == w {
				wc.expireWorkers = append(wc.expireWorkers[:i], wc.expireWorkers[i+1:]...)
			}
		}
	}
}
func (wc *WorkerController) Start() {
	wc.locker.Lock()
	defer wc.locker.Unlock()
	wc.new()
}
func (wc *WorkerController) NewWorker() error {
	wc.locker.Lock()
	defer wc.locker.Unlock()
	return wc.new()
}
func (wc *WorkerController) new() error {

	buf := bytes.NewBuffer(nil)
	var fileAll []*os.File
	index := 3
	for _, dm := range wc.dms {
		data, files, err := dm.Encode(index)

		if err != nil {
			return err
		}
		index += len(files)
		fileAll = append(fileAll, files...)
		buf.Write(data)

	}

	workerProcess, err := wc.newWorkerProcess(buf, fileAll)
	if err != nil {
		return err
	}

	if wc.current != nil {
		wc.expireWorkers = append(wc.expireWorkers, wc.current)
	}
	wc.current = workerProcess
	go wc.check(wc.current)
	return nil
}

func (wc *WorkerController) getClient() service.WorkerServiceClient {
	wc.locker.Lock()
	defer wc.locker.Unlock()
	return wc.current
}
