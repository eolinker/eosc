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

	"github.com/eolinker/eosc/utils"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/common/bean"

	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/traffic"
)

func Process() {

	utils.InitLogTransport(eosc.ProcessWorker)
	//log.Debug("load plugin env...")
	log.Info("worker process start...")

	register := loadPluginEnv()
	log.Debug("create worker...")
	w, err := NewProcessWorker(register)
	if err != nil {
		log.Error("new process worker error: ", err)
		return
	}

	w.Start()

	w.wait()
	log.Info("worker process end")
}

type ProcessWorker struct {
	tf          traffic.ITraffic
	professions IProfessions
	workers     IWorkers

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
func NewProcessWorker(extends eosc.IExtenders) (*ProcessWorker, error) {
	workerServer, err := NewWorkerServer()
	if err != nil {
		return nil, err
	}
	w := &ProcessWorker{
		workerServer: workerServer,
	}

	tf := traffic.NewTraffic()
	w.tf = tf
	ps := NewProfessions()
	w.professions = ps
	wm := NewWorkerManager(w.professions)
	w.workers = wm

	tf.Read(os.Stdin)

	bean.Injection(&w.tf)
	bean.Injection(&w.professions)
	var iWorkers eosc.IWorkers = w.workers
	bean.Injection(&iWorkers)
	bean.Check()

	psData, err := ReadProfessionData(os.Stdin)
	if err != nil {
		log.Warn("profession configs error:", err)
		return nil, err
	}
	ps.init(psData, extends)
	workersData := ReadWorkers(os.Stdin)

	err = wm.Init(workersData)
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
	w.workerServer.SetTraffic(w.tf)
	w.workerServer.SetWorkers(w.workers)

	return nil
}
