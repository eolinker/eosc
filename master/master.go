/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package master

import (
	"time"

	eosc_args "github.com/eolinker/eosc/eosc-args"

	"github.com/eolinker/eosc/log"

	"github.com/eolinker/eosc/traffic"

	"github.com/eolinker/eosc/log/filelog"
	"google.golang.org/grpc"

	"github.com/eolinker/eosc/master/service"

	"fmt"

	"os"
	"os/signal"
	"syscall"

	"github.com/eolinker/eosc/process"
)

func Process() {
	master := NewMasterHandle()
	master.Start()
	master.Wait()
}

type Master struct {

	masterTraffic traffic.IController
	workerTraffic traffic.IController

	srv *grpc.Server

}

func (m *Master) InitLogTransport() {
	writer := filelog.NewFileWriteByPeriod()
	writer.Set(fmt.Sprintf("/var/log/%s", process.AppName()), "error.log", filelog.PeriodDay, 7*24*time.Hour)
	writer.Open()
	transport := log.NewTransport(writer, log.InfoLevel)
	transport.SetFormatter(&log.LineFormatter{
		TimestampFormat:  "[2006-01-02 15:04:05]",
		CallerPrettyfier: nil,
	})
	log.Reset(transport)
}

func (m *Master) Start() {

	m.masterTraffic = traffic.NewController(os.Stdin)
	m.workerTraffic = traffic.NewController(os.Stdin)

	m.InitLogTransport()

	log.Info("start master")
	srv, err := service.StartMaster(fmt.Sprintf("/tmp/%s.master.sock", process.AppName()))
	if err != nil {
		log.Error(err)
		os.Exit(1)
		return
	}

	m.srv = srv

	ip := os.Getenv(fmt.Sprintf("%s_%s", process.AppName(), eosc_args.IP))
	port := os.Getenv(fmt.Sprintf("%s_%s", process.AppName(), eosc_args.Port))
	log.Info(fmt.Sprintf("%s:%s", ip, port))
	// 监听master监听地址，用于接口处理
	_,err = m.masterTraffic.ListenTcp("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		log.Error(err)
		os.Exit(1)
		return
	}

	//TODO 若该进程是master的子进程，则给父进程一个退出信号
	pEnv := fmt.Sprintf("%s_%s",process.AppName(),"IS_MASTER_CHILD")
	if  os.Getenv(pEnv) != "" {
		syscall.Kill(syscall.Getppid(), syscall.SIGQUIT)
	}

}

func (m *Master) Wait() error {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	for {
		sig := <-sigc
		switch sig {
		case os.Interrupt, os.Kill, syscall.SIGQUIT:
			{
				log.Infof("Caught signal %s: shutting down.\n", sig.String())
				m.srv.GracefulStop()
				m.close()
				return nil
			}
		case syscall.SIGUSR1:
			{
				// TODO: 平滑重启操作
				process.Fork()  //传子进程需要的内容
			}
		default:
			continue
		}
	}
}

func (m *Master) close() {
	syscall.Unlink(fmt.Sprintf("/tmp/%s.master.sock", process.AppName()))

}

func NewMasterHandle() *Master {
	return &Master{}
}
