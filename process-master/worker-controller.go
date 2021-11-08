package process_master

import (
	"context"
	"os"
	"sync"
	"time"

	raft_service "github.com/eolinker/eosc/raft/raft-service"

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

	ctx        context.Context
	cancelFunc context.CancelFunc

	restartChan chan int
}

func NewWorkerController(traffic traffic.IController, config *config.Config, extenderSetting extenders.ITypedExtenderSetting, professions eosc.IProfessions, workers *WorkerConfigs, workerServiceProxy *WorkerServiceProxy) *WorkerController {
	traffics, files := traffic.Export(3)

	ctx, cancelFunc := context.WithCancel(context.Background())

	wc := &WorkerController{
		workerServiceProxy: workerServiceProxy,
		traffics:           traffics,
		trafficFiles:       files,
		listensMsg:         config.Export(),
		extenderSetting:    extenderSetting,
		professions:        professions,
		workers:            workers,
		ctx:                ctx,
		cancelFunc:         cancelFunc,
		restartChan:        make(chan int, 1),
	}
	go wc.doControl()
	return wc
}

func (wc *WorkerController) Stop() {
	wc.locker.Lock()
	defer wc.locker.Unlock()

	if wc.isStop {
		return
	}
	if wc.cancelFunc != nil {
		wc.cancelFunc()
		wc.cancelFunc = nil
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

func (wc *WorkerController) restart() bool {

	process := wc.getClient()
	if process == nil {
		err := wc.NewWorker()
		if err != nil {
			return true
		}
		return true
	}
	oldExtenderSetting := process.extenderSetting
	deletedExtenderSetting := process.extendersDeleted
	extenderSetting := wc.extenderSetting.All()

	for id, oVersion := range oldExtenderSetting {
		v, has := extenderSetting[id]
		if !has {
			if _, ok := deletedExtenderSetting[id]; !ok {
				deletedExtenderSetting[id] = oVersion
			}
			continue
		}
		delete(extenderSetting, id)
		if v != oVersion {
			// 存在不同版本的，直接重启
			err := wc.NewWorker()
			if err != nil {
				log.Errorf("restart worker process:", err)
				return true
			}
			return true
		}
		// 版本一致，不做操作
	}

	for id, version := range extenderSetting {
		if dv, has := deletedExtenderSetting[id]; has {
			if version != dv {
				// 该项目版本已经在被删除列表中， 且版本不一致，需要重启
				err := wc.NewWorker()
				if err != nil {
					log.Errorf("restart worker process:", err)
					return true
				}
				return true
			}
			// 以删除的版本一致
			delete(deletedExtenderSetting, id)
		}
		oldExtenderSetting[id] = version
	}

	if len(extenderSetting) > 0 {
		_, err := wc.workerServiceProxy.AddExtender(wc.ctx, &service.WorkerAddExtender{
			Extenders: extenderSetting,
		})
		if err != nil {
			log.Error("call addExtender:", err)
		}
	}

	return false

}
func (wc *WorkerController) reset() {

	resetArg := &service.ResetRequest{
		Professions: wc.professions.All(),
		Workers:     wc.workers.export(),
	}

	wc.workerServiceProxy.Reset(wc.ctx, resetArg)
}

func (wc *WorkerController) tryRestart() {
	wc.restartChan <- 1
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

		_, err := wc.workerServiceProxy.Ping(context.TODO(), &service.WorkerHelloRequest{Hello: "hello"})
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

	wc.workerServiceProxy.SetWorkerProcess(wc.current.client)

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
	log.Debug("worker controller get event:", event)
	switch event {
	case raft_service.EventComplete:
		wc.reset()
		wc.tryRestart()
	}
	return
}
func (wc *WorkerController) raftCommitEvent(namespace, cmd string) {
	log.Debug("worker controller get comit event:", namespace, ":", cmd)

	switch namespace {
	case extenders.NamespaceExtenders:
		switch cmd {
		case extenders.CommandSet:
			wc.tryRestart()
		case extenders.CommandDelete:
			wc.tryRestart()
		}
	}
}

func (wc *WorkerController) doControl() {
	t := time.NewTimer(time.Second)
	t.Stop()
	defer t.Stop()
	for {
		select {
		case <-wc.ctx.Done():

			return
		case <-wc.restartChan:
			t.Reset(time.Second)
		case <-t.C:
			wc.restart()
		}
	}
}
