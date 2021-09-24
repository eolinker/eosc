package process_master

import (
	"bytes"
	"os"
	"os/exec"
	"sync"

	"github.com/eolinker/eosc/admin"
	"github.com/eolinker/eosc/process"
	"github.com/eolinker/eosc/service"
	"github.com/eolinker/eosc/traffic"
)

type iWorkerClientPoll interface {
	GetWorkerServiceClient() service.WorkerServiceClient
	Close()
	Start() error
}

type WorkerController struct {
	locker     sync.Mutex
	tc         traffic.IController
	profession admin.IProfessions
}

func NewWorkerController(tc traffic.IController) *WorkerController {
	return &WorkerController{tc: tc}
}
func (wc *WorkerController) Stop() {
	wc.locker.Lock()
	defer wc.locker.Unlock()

}

func (wc *WorkerController) Start() {

}

func (wc *WorkerController) new() (*exec.Cmd, error) {
	worker, err := process.Cmd("workers", nil)
	if err != nil {

		return nil, err
	}
	data, files, err := wc.tc.Encode(3)

	if err != nil {
		return nil, err
	}
	worker.Stdin = bytes.NewBuffer(data)
	worker.Stdout = os.Stdout
	worker.Stderr = os.Stderr
	worker.ExtraFiles = files
	err = worker.Start()
	if err != nil {
		return nil, err
	}
	return worker, nil
}
