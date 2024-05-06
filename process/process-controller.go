package process

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/utils"
)

type IProcessUpdates []IProcessUpdate

func (I IProcessUpdates) Update(cmd *exec.Cmd) {
	for _, i := range I {
		i.Update(cmd)
	}
}

type IProcessUpdate interface {
	Update(cmd *exec.Cmd)
}

type IConfigBuild interface {
	Config() StartArgs
}

type StartArgs struct {
	Data       []byte
	ExtraFiles []*os.File
}

type ProcessController struct {
	ctx         context.Context
	cancel      context.CancelFunc
	name        string
	current     *ProcessCmd
	callback    IProcessUpdates
	restartChan chan *StartArgs
	locker      sync.Mutex
	isStop      bool
	isShutDown  int32
	logWriter   io.Writer
}

func NewProcessController(ctx context.Context, name string, logWriter io.Writer, callback ...IProcessUpdate) *ProcessController {

	newCtx, cancel := context.WithCancel(ctx)
	c := &ProcessController{
		callback:    IProcessUpdates(callback),
		name:        name,
		ctx:         newCtx,
		cancel:      cancel,
		restartChan: make(chan *StartArgs),
		logWriter:   logWriter,
	}
	atomic.StoreInt32(&c.isShutDown, 1)
	go c.doControl()
	return c
}

func (pc *ProcessController) Shutdown() {
	pc.locker.Lock()
	defer pc.locker.Unlock()
	atomic.StoreInt32(&pc.isShutDown, 1)
	if pc.current != nil {
		pc.current.Close()
		pc.current = nil
	}
}
func (pc *ProcessController) Stop() {
	pc.locker.Lock()
	defer pc.locker.Unlock()

	if pc.isStop {
		return
	}
	if pc.cancel != nil {
		pc.cancel()
		pc.cancel = nil
	}

	pc.isStop = true
}
func newProcess(name string, data []byte, logWriter io.Writer, extraFiles []*os.File) (*ProcessCmd, error) {
	cmd, err := Cmd(name, nil)
	if err != nil {
		return nil, err
	}

	reader, writer, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	cmd.Stdin = bytes.NewReader(utils.EncodeFrame(data))
	cmd.Stdout = writer
	cmd.Stderr = logWriter
	cmd.ExtraFiles = extraFiles

	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	pc := NewProcessCmd(name, cmd, reader)
	go pc.Read()
	return pc, nil
}

func (pc *ProcessController) check(w *ProcessCmd, configData []byte, extraFiles []*os.File) {
	err := w.Wait()
	if err != nil {
		log.Warn("worker exit:", err)
	}
	if atomic.LoadInt32(&pc.isShutDown) == 1 {
		return
	}
	pc.locker.Lock()
	defer pc.locker.Unlock()
	if pc.current == w {
		// 连接断开
		err = pc.create(configData, extraFiles)
		if err != nil {
			log.Error("worker create:", err)
		}
	}
}

func (pc *ProcessController) run(configData []byte, extraFiles []*os.File) error {
	log.DebugF("create %s process start...\n", pc.name)

	p, err := newProcess(pc.name, configData, pc.logWriter, extraFiles)
	if err != nil {
		log.Warnf("new %s process: %s", pc.name, err.Error())
		return err
	}

	old := pc.current
	pc.current = p

	go pc.check(pc.current, configData, extraFiles)

	if old != nil {
		old.Close()
	}

	return nil
}
func (pc *ProcessController) create(configData []byte, extraFiles []*os.File) error {
	err := pc.run(configData, extraFiles)
	if err != nil {
		log.Warn("new process[", pc.name, "]:", err)
		return err
	}

	ticker := time.NewTicker(time.Millisecond * 5)
	defer ticker.Stop()
	defer utils.TimeSpend(fmt.Sprint("wait [", pc.name, "] process start:"))()
	log.Debug(pc.name, " controller ping...")
	defer func() {
		log.Debug(pc.name, " controller ping done")
	}()
	for {
		select {
		case <-pc.ctx.Done():
			log.Debug(pc.name, " end")
			return nil
		case <-ticker.C:
			if pc.current == nil {
				pc.callback.Update(nil)
				return errors.New("process not exist")
			}

			switch status := pc.current.Status(); status {
			case StatusRunning:
				pc.callback.Update(pc.current.Cmd())
				return nil
			case StatusExit, StatusError:
				pc.callback.Update(nil)
				return errors.New("fail to start process " + pc.name + " " + strconv.Itoa(status))
			case StatusStart:
				// continue
			}
		}
	}
}

func (pc *ProcessController) Start(configData []byte, extraFiles []*os.File) error {
	pc.locker.Lock()
	defer pc.locker.Unlock()
	atomic.StoreInt32(&pc.isShutDown, 0)
	return pc.create(configData, extraFiles)
}

func (pc *ProcessController) TryRestart(configData []byte, extraFiles []*os.File) {

	pc.restartChan <- &StartArgs{
		Data:       configData,
		ExtraFiles: extraFiles,
	}
}

func (pc *ProcessController) restart(configData []byte, extraFiles []*os.File) {
	pc.locker.Lock()
	defer pc.locker.Unlock()

	err := pc.create(configData, extraFiles)
	if err != nil {
		log.Error("restart error: ", err)
	}

}

func (pc *ProcessController) doControl() {
	t := time.NewTimer(time.Second)
	t.Stop()
	defer t.Stop()
	var lastConfig = new(StartArgs)

	for {
		select {
		case <-pc.ctx.Done():
			pc.Shutdown()
			return
		case arg, ok := <-pc.restartChan:
			if ok {
				lastConfig = arg
				t.Reset(time.Second)
			}

		case <-t.C:
			if atomic.LoadInt32(&pc.isShutDown) == 0 {
				pc.restart(lastConfig.Data, lastConfig.ExtraFiles)
			}
		}
	}
}

//func (pc *ProcessController) Sleep() error {
//	client := pc.getClient()
//	if client != nil {
//		err := client.Sleep()
//		if err != nil {
//			pc.callback.Update(nil)
//			return err
//		}
//		pc.callback.Update(client.Cmd())
//		return nil
//	}
//	pc.callback.Update(nil)
//	return errors.Start("process not exist")
//}
