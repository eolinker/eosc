/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package process_worker

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/eolinker/eosc/config"
	grpc_unixsocket "github.com/eolinker/eosc/grpc-unixsocket"
	"github.com/eolinker/eosc/router"

	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/eolinker/eosc/process"

	"github.com/eolinker/eosc/extends"

	"github.com/eolinker/eosc/service"
	"google.golang.org/protobuf/proto"

	"github.com/eolinker/eosc/utils"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/common/bean"

	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/traffic"
)

func Process() {

	//utils.InitStdTransport(eosc.ProcessWorker)
	arg := readArg()
	log.Info("worker process start...")

	log.Debug("create worker...")
	w, err := NewProcessWorker(arg)
	if err != nil {
		writeOutput(process.StatusExit, err.Error())
		log.Error("new process worker error: ", err)
		return
	}

	w.Start()
	writeOutput(process.StatusRunning, "")

	w.wait()
	log.Info("worker process end")
}

func writeOutput(status int, msg string) {
	data := new(eosc.ProcessStatus)
	data.Status = int32(status)
	data.Msg = msg
	d, _ := proto.Marshal(data)
	err := utils.WriteFrame(os.Stdout, d)
	if err != nil {
		log.Error("write output error: ", err)
	}
}

type ProcessWorker struct {
	tf traffic.ITraffic

	once          sync.Once
	server        *WorkerServer
	metricsServer *http.Server
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
				w.close()
			}
		default:
			continue
		}
	}

}

// NewProcessWorker 创建新的worker进程
// 启动时通过stdin传输配置信息
func NewProcessWorker(arg *service.ProcessLoadArg) (*ProcessWorker, error) {

	register := extends.InitRegister()
	tf := createTraffic(arg.Traffic)
	bean.Injection(&tf)
	var listenUrl = new(config.ListenUrl)
	*listenUrl = arg.ListensMsg
	bean.Injection(&listenUrl)

	extends.LoadPlugins(arg.Extends, register)
	var extenderDrivers eosc.IExtenderDrivers = register
	bean.Injection(&extenderDrivers)

	server, err := NewWorkerServer(os.Getppid(), register, func() {
		bean.Check()
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	w := newProcessWorker(tf, server)

	return w, nil
}

func newProcessWorker(tf traffic.ITraffic, server *WorkerServer) *ProcessWorker {
	w := &ProcessWorker{
		tf:            tf,
		server:        server,
		metricsServer: &http.Server{},
	}

	w.metricsServer.Handler = router.GetHandler()

	return w
}

func (w *ProcessWorker) close() {
	w.once.Do(func() {
		w.tf.Close()
		w.server.Stop()

		timeout, _ := context.WithTimeout(context.Background(), time.Second*3)
		w.metricsServer.Shutdown(timeout)
	})
}

func (w *ProcessWorker) Start() error {
	w.OpenMetricsServer()
	return nil
}

func (w *ProcessWorker) OpenMetricsServer() error {
	addr := service.ServerUnixAddr(os.Getpid(), eosc.ProcessWorker)
	syscall.Unlink(addr)
	log.Info("start worker unix server: ", addr)
	l, err := grpc_unixsocket.Listener(addr)
	if err != nil {
		return err
	}

	go func() {
		err := w.metricsServer.Serve(l)
		if err != nil {
			log.Info("http server error: ", err)
		}
		return
	}()

	return nil
}

func readArg() *service.ProcessLoadArg {
	arg := new(service.ProcessLoadArg)
	frame, err := utils.ReadFrame(os.Stdin)
	if err != nil {
		log.Warn("read arg fail:", err)
		return arg
	}
	err = json.Unmarshal(frame, arg)
	if err != nil {
		log.Warn("unmarshal arg fail:", err)
		return arg
	}
	log.Debug("read arg: ", arg)
	return arg
}

func createTraffic(tfConf []*traffic.PbTraffic) traffic.ITraffic {
	t := traffic.FromArg(tfConf)

	return t
}
