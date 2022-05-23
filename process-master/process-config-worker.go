package process_master

import (
	"github.com/eolinker/eosc/process"
	"github.com/eolinker/eosc/process-master/extender"
	"github.com/eolinker/eosc/service"
	"github.com/eolinker/eosc/traffic"
	"google.golang.org/protobuf/proto"
	"os"
	"sync"

	"github.com/eolinker/eosc/config"
)

type WorkerController struct {
	workerProcess *process.ProcessController
	extends       map[string]string
	locker        sync.Mutex
	traffics      []*traffic.PbTraffic
	trafficFiles  []*os.File
	listensMsg    *config.ListensMsg
	isRunning     bool
}

func (wc *WorkerController) Stop() {
	wc.workerProcess.Stop()
}
func (wc *WorkerController) Update(status []*extender.Status, success bool) {
	if success {
		extends := make(map[string]string)
		for _, s := range status {
			extends[s.Name()] = s.Version
		}
		wc.locker.Lock()
		wc.extends = extends
		wc.locker.Unlock()

		args := &service.ProcessLoadArg{
			Traffic:    wc.traffics,
			ListensMsg: wc.listensMsg,
			Extends:    extends,
		}
		data, _ := proto.Marshal(args)
		if wc.isRunning {
			wc.workerProcess.TryRestart(data, wc.trafficFiles)
		} else {
			wc.isRunning = true
			wc.workerProcess.Start(data, wc.trafficFiles)
		}
	}
}

func NewWorkerController(traffic traffic.IController, config *config.Config, workerProcess *process.ProcessController) *WorkerController {
	traffics, files := traffic.Export(3)
	wc := &WorkerController{
		traffics:      traffics,
		trafficFiles:  files,
		listensMsg:    config.Export(),
		workerProcess: workerProcess,
	}
	return wc
}
