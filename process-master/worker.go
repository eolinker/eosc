/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package process_master

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/process"
	"github.com/eolinker/eosc/traffic"
)

type Param struct {
}

type IWorkerParam interface {
	Encode() (io.ReadCloser, []*os.File)
}

type WorkerController struct {
	cmd    *exec.Cmd
	locker sync.Mutex
	tc     traffic.IController
}

func NewWorkerController(tc traffic.IController) *WorkerController {
	return &WorkerController{tc: tc}
}
func (wc *WorkerController) Stop() {
	wc.locker.Lock()
	defer wc.locker.Unlock()
	if wc.cmd != nil {
		err := wc.cmd.Process.Signal(syscall.SIGUSR1)
		if err != nil {
			log.Warn("Signal error:", err)
			return
		}
		err = wc.cmd.Wait()
		if err != nil {
			log.Warn("stop workers:", err)
			return
		}
	}
}

func (wc *WorkerController) Start() {

	wc.Restart()
	go func() {

	}()
}
func (wc *WorkerController) Restart() {
	wc.locker.Lock()
	defer wc.locker.Unlock()

	if wc.cmd != nil {
		err := wc.cmd.Process.Signal(syscall.SIGUSR1)
		if err != nil {
			log.Warn("Signal error:", err)
		}
		wc.cmd = nil
	}
}

func (wc *WorkerController) new() (*exec.Cmd, error) {
	worker, err := process.Cmd("workers", nil)
	if err != nil {
		log.Warn(err)
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

func getWaitCmd(cmd *exec.Cmd) <-chan error {
	c := make(chan error, 0)
	go func() {
		c <- cmd.Wait()
		close(c)
	}()
	return c
}
