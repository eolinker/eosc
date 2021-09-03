/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package master

import (
	"bytes"
	"github.com/eolinker/eosc/process"
	"github.com/eolinker/eosc/traffic"
	"io"
	"github.com/eolinker/eosc/log"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

type Param struct {

}

type IWorkerParam interface {
	Encode()(io.ReadCloser,[]*os.File)
}

type WorkerController struct {
	cmd *exec.Cmd
	locker sync.Mutex
	tc traffic.IController
}

func NewWorkerController(tc traffic.IController) *WorkerController {
	return &WorkerController{tc: tc}
}
func (wc *WorkerController) Stop() {
	wc.locker.Lock()
	defer wc.locker.Unlock()
	if wc.cmd != nil{
		err:=wc.cmd.Process.Signal(syscall.SIGUSR1)
		if err!= nil{
			log.Warn("Signal error:",err)
			return
		}
		err =wc.cmd.Wait()
		if err!=nil{
			log.Warn("stop worker:",err)
			return
		}
	}
}
func (wc *WorkerController) Start() {
	wc.Restart()
}
func (wc *WorkerController)Restart()  {
	wc.locker.Lock()
	defer wc.locker.Unlock()

	if wc.cmd != nil{
		err:=wc.cmd.Process.Signal(syscall.SIGUSR1)
		if err!= nil{
			log.Warn("Signal error:",err)
		}

	}
}

func (wc *WorkerController) new() ( *exec.Cmd,error) {
	worker,err :=process.Cmd("worker",nil)
	if err != nil{
		log.Warn(err)
		return nil,err
	}
	traffics:=wc.tc.All()
	buf:=&bytes.Buffer{}
	files,err := traffics.WriteTo(buf)
	if err != nil {
		return nil,err
	}
	worker.Stdin = buf
	worker.Stdout = os.Stdout
	worker.Stderr = os.Stderr
	worker.ExtraFiles = files
	err = worker.Start()
	if err != nil {
		return nil,err
	}
	return worker,nil
}