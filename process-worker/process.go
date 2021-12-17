/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package process_worker

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/eolinker/eosc/service"
	"google.golang.org/protobuf/proto"

	"github.com/eolinker/eosc/utils"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/common/bean"

	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/traffic"
)

func Process() {
	arg := readArg()
	level, _ := log.ParseLevel(arg.Level)
	utils.InitStdTransport(eosc.ProcessWorker, level)
	//log.Debug("load plugin env...")
	log.Info("worker process start...")

	log.Debug("create worker...")
	w, err := NewProcessWorker(arg)
	if err != nil {
		log.Error("new process worker error: ", err)
		return
	}

	w.Start()

	w.wait()
	log.Info("worker process end")
}

type ProcessWorker struct {
	tf traffic.ITraffic

	workers IWorkers

	once         sync.Once
	workerServer *WorkerServer
}

func (w *ProcessWorker) wait() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	for {
		sig := <-sigc
		log.Infof("Caught signal pid:%d ppid:%d signal %s: .\n", os.Getpid(), os.Getppid(), sig.String())
		switch sig {
		case os.Interrupt, os.Kill:
			{
				w.close()
				return
			}
		case syscall.SIGQUIT:
			{
				w.close()
				return
			}
		case syscall.SIGUSR1:
			{

			}
		default:
			continue
		}
	}

}

//NewProcessWorker 创建新的worker进程
//启动时通过stdin传输配置信息
func NewProcessWorker(arg *service.WorkerLoadArg) (*ProcessWorker, error) {

	register := loadPluginEnv(arg.ExtenderSetting)

	tf := createTraffic(arg.Traffic)
	professions := NewProfessions()
	professions.Reset(arg.Professions, register)
	wm := NewWorkerManager(professions)
	workerServer, err := NewWorkerServer(wm, register, professions)
	if err != nil {
		return nil, err
	}
	w := &ProcessWorker{
		workerServer: workerServer,
		workers:      wm,
		tf:           tf,
	}
	var extenderDrivers eosc.IExtenderDrivers = register
	bean.Injection(&extenderDrivers)
	bean.Injection(&tf)
	bean.Injection(&professions)

	var iWorkers eosc.IWorkers = w.workers
	bean.Injection(&iWorkers)
	bean.Injection(&arg.ListensMsg)
	bean.Check()

	err = wm.Reset(arg.Workers)
	if err != nil {
		log.Warn("worker configs error:", err)
		return nil, err
	}

	w.workers = wm
	//ports32 := wm.portsRequire.All()
	//ports := make([]int, len(ports32))
	//for i, v := range ports32 {
	//	ports[i] = int(v)
	//}
	//w.tf.Expire(ports)
	return w, nil
}

func (w *ProcessWorker) close() {

	w.once.Do(func() {
		w.tf.Close()
		w.workerServer.Stop()
	})

}

func (w *ProcessWorker) Start() error {

	return nil
}

func readArg() *service.WorkerLoadArg {
	arg := new(service.WorkerLoadArg)
	frame, err := utils.ReadFrame(os.Stdin)
	if err != nil {
		log.Warn("read arg fail:", err)
		return arg
	}
	err = proto.Unmarshal(frame, arg)
	if err != nil {
		log.Warn("unmarshal arg fail:", err)
		return arg
	}
	log.Debug("read arg: ", arg)
	return arg
}
func createTraffic(tfConf []*traffic.PbTraffic) traffic.ITraffic {
	t := traffic.NewTraffic()

	err := t.Read(tfConf)
	if err != nil {
		log.Error("read traffic :", err)
		return t
	}
	return t
}
