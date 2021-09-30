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
	"errors"
	"os"
	"syscall"
	"time"

	"github.com/eolinker/eosc/env"

	"github.com/eolinker/eosc"

	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/pidfile"
	"github.com/eolinker/eosc/process"
)

var runningMasterForked bool

//Fork Master fork 子进程，入参为子进程需要的内容
func (m *Master) Fork() error {
	if runningMasterForked {
		return errors.New("Another process already forked. Ignoring this one.")
	}
	err := m.PID.TryFork()
	if err != nil {
		return err
	}
	runningMasterForked = true

	// 子进程的环境变量加入MASTER_CONTINUE字段，用于新的Master启动后给父Master传送中断信号
	env := append(os.Environ(), env.GenEnv("MASTER_CONTINUE", "1"))
	dataMasterTraffic, filesMaster, err := m.masterTraffic.Encode(3)
	if err != nil {
		return err
	}

	dataWorkerTraffic, filesWorker, err := m.workerTraffic.Encode(len(filesMaster) + 3)
	if err != nil {
		return err
	}

	data := make([]byte, len(dataMasterTraffic)+len(dataWorkerTraffic))
	copy(data, dataMasterTraffic)
	copy(data[len(dataMasterTraffic):], dataWorkerTraffic)

	cmd, err := process.Cmd(eosc.ProcessMaster, os.Args[1:])
	if err != nil {
		return err
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = bytes.NewBuffer(data)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}
	cmd.ExtraFiles = append(filesMaster, filesWorker...)
	cmd.Env = env

	err = cmd.Start()
	if err != nil {
		log.Fatalf("Restart: Failed to launch, error: %v", err)
		return err
	}
	log.Debug("fork new process: ", cmd.String())
	// check cmd
	go m.waitFork(cmd.Process.Pid)

	return nil
}

func (m *Master) waitFork(pid int) {
	t := time.NewTicker(time.Millisecond * 100)
	defer t.Stop()
	for {
		select {
		case <-m.ctx.Done():
			return
		case <-t.C:
			if !pidfile.ProcessExists(pid) {
				m.PID.UnFork()
				return
			}
		}
	}

}
