/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package master

import (
	"context"
	"github.com/eolinker/eosc/service"
	"time"

	"github.com/eolinker/eosc/store"

	"github.com/eolinker/eosc"

	eosc_args "github.com/eolinker/eosc/eosc-args"
	"github.com/eolinker/eosc/log"

	"github.com/eolinker/eosc/traffic"

	"github.com/eolinker/eosc/log/filelog"
	"google.golang.org/grpc"

	"fmt"

	"os"
	"os/signal"
	"syscall"

	"github.com/eolinker/eosc/process"
)

func Process() {

	if process.CheckPIDFILEAlreadyExists() {
		// 存在，则报错开启失败
		log.Error("the master is running")
		return
	}

	master := NewMasterHandle()
	master.Start()

	if _, has := eosc_args.GetEnv("MASTER_CONTINUE"); has {
		syscall.Kill(syscall.Getppid(), syscall.SIGQUIT)
	}
	err := process.CreatePidFile()
	if err != nil {
		// 创建pid文件失败，则报错

		return
	}
	master.Wait()
}

type Master struct {
	service.UnimplementedMasterServer
	service.UnimplementedCtiServiceServer

	masterTraffic traffic.IController
	workerTraffic traffic.IController
	store     eosc.IStore
	masterSrv *grpc.Server
}

func (m *Master) Hello(ctx context.Context, request *service.HelloRequest) (*service.HelloResponse, error) {
	return &service.HelloResponse{
		Name: request.GetName(),
	},nil
}


func (m *Master) InitLogTransport() {
	writer := filelog.NewFileWriteByPeriod()
	writer.Set(fmt.Sprintf("/var/log/%s", process.AppName()), "error.log", filelog.PeriodDay, 7*24*time.Hour)
	writer.Open()
	transport := log.NewTransport(writer, log.InfoLevel)
	formater := &log.LineFormatter{
		TimestampFormat:  "[2006-01-02 15:04:05]",
		CallerPrettyfier: nil,
	}
	transport.SetFormatter(formater)
	log.NewStdTransport(formater)
	log.Reset(transport, log.NewStdTransport(formater))
}

func (m *Master) Start() {

	m.masterTraffic = traffic.NewController(os.Stdin)
	m.workerTraffic = traffic.NewController(os.Stdin)

	m.InitLogTransport()
	// 设置存储操作
	s, err := store.NewStore()
	if err != nil {
		log.Error("new store error: ", err.Error())
		return
	}
	m.store = s

	log.Info("start master")
	srv, err := m.StartMaster()
	if err != nil {
		log.Error("start master grpc server error: ", err.Error())
		return
	}
	m.masterSrv = srv
	log.Debug("RegisterCtiServiceServer")



	ip := os.Getenv(fmt.Sprintf("%s_%s", process.AppName(), eosc_args.IP))
	port := os.Getenv(fmt.Sprintf("%s_%s", process.AppName(), eosc_args.Port))
	log.Info(fmt.Sprintf("%s:%s", ip, port))
	// 监听master监听地址，用于接口处理
	_, err = m.masterTraffic.ListenTcp("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		log.Error(err)
		os.Exit(1)
		return
	}

}

func (m *Master) Wait() error {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	for {
		sig := <-sigc
		switch sig {
		case os.Interrupt, os.Kill:
			{
				log.Infof("Caught signal %s: shutting down.\n", sig.String())

				m.close()
				return nil
			}
		case syscall.SIGQUIT:
			{
				m.close()
				process.ClearPid()
			}
		case syscall.SIGUSR1:
			{
				// TODO: 平滑重启操作
				m.Fork() //传子进程需要的内容
				process.GetPidByFile()
			}
		default:
			continue
		}
	}
}

func (m *Master) close() {

	syscall.Unlink(fmt.Sprintf("/tmp/%s.master.sock", process.AppName()))

	m.masterSrv.GracefulStop()

}

func NewMasterHandle() *Master {
	return &Master{}
}