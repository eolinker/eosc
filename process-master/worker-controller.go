package process_master

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/eolinker/eosc/config"

	"github.com/eolinker/eosc/process-master/extenders"

	"github.com/eolinker/eosc/utils"

	"github.com/eolinker/eosc/traffic"

	"github.com/eolinker/eosc/log"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/service"
)

type WorkerController struct {
	locker             sync.Mutex
	workerServiceProxy *WorkerServiceProxy

	current       *WorkerProcess
	expireWorkers []*WorkerProcess

	traffics     []*traffic.PbTraffic
	trafficFiles []*os.File
	listensMsg   *config.ListensMsg
	isStop       bool

	extenderSetting extenders.ITypedExtenderSetting
	professions     eosc.IProfessions
	workers         *WorkerConfigs
}

func NewWorkerController(traffic traffic.IController, config *config.Config, extenderSetting extenders.ITypedExtenderSetting, professions eosc.IProfessions, workers *WorkerConfigs, workerServiceProxy *WorkerServiceProxy) *WorkerController {
	traffics, files := traffic.Export(3)
	return &WorkerController{
		workerServiceProxy: workerServiceProxy,
		traffics:           traffics,
		trafficFiles:       files,
		listensMsg:         config.Export(),
		extenderSetting:    extenderSetting,
		professions:        professions,
		workers:            workers,
	}
}

func (wc *WorkerController) Stop() {
	wc.locker.Lock()
	defer wc.locker.Unlock()

	if wc.isStop {
		return
	}

	wc.isStop = true
	if wc.current != nil {
		wc.workerServiceProxy.SetWorkerProcess(nil)
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

	if wc.getClient() == w {
		err := wc.NewWorker()
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

func (wc *WorkerController) restart() {

	process := wc.getClient()
	if process == nil {
		err := wc.NewWorker()
		if err != nil {
			return
		}
		return
	}
}

func (wc *WorkerController) NewWorker() error {

	wc.locker.Lock()
	defer wc.locker.Unlock()
	err := wc.new()
	if err != nil {
		log.Warn("new worker:", err)
		return err
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	defer utils.Timeout("wait worker process start:")()

	for {
		_, err := wc.current.Ping(context.TODO(), &service.WorkerHelloRequest{Hello: "hello"})
		if err == nil {
			return nil
		}

		log.Debug("work controller ping: ", err)
		<-ticker.C
	}

}
func (wc *WorkerController) config() (*service.WorkerLoadArg, []*os.File) {

	return &service.WorkerLoadArg{
		Traffic:         wc.traffics,
		ListensMsg:      wc.listensMsg,
		ExtenderSetting: wc.extenderSetting.All(),
		Professions:     wc.professions.All(),
		Workers:         wc.workers.export(),
	}, wc.trafficFiles
}
func (wc *WorkerController) new() error {
	log.Debug("create worker process start")

	arg, files := wc.config()

	workerProcess, err := newWorkerProcess(arg, files)
	if err != nil {
		log.Warn("new worker process:", err)
		return err
	}
	workerProcess.createClient()
	old := wc.current
	wc.current = workerProcess

	wc.workerServiceProxy.SetWorkerProcess(wc.current)

	go wc.check(wc.current)

	if old != nil {
		old.Close()
	}

	return nil
}

func (wc *WorkerController) getClient() *WorkerProcess {
	wc.locker.Lock()
	defer wc.locker.Unlock()

	return wc.current
}

func (wc *WorkerController) raftEvent(event string) {

}
func (wc *WorkerController) raftCommitEvent(namespace, cmd string) {

}
