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

	"github.com/eolinker/eosc/pidfile"
	"github.com/eolinker/eosc/service"

	"github.com/eolinker/eosc/store"

	"github.com/eolinker/eosc"

	eosc_args "github.com/eolinker/eosc/eosc-args"
	"github.com/eolinker/eosc/log"

	"github.com/eolinker/eosc/traffic"

	"google.golang.org/grpc"

	"fmt"

	"os"
	"os/signal"
	"syscall"

	"github.com/eolinker/eosc/process"
)

func Process() {
	file, err := pidfile.New()
	if err != nil {
		log.Errorf("the master is running:%v", err)
		return
	}
	master := NewMasterHandle(file)
	if err := master.Start(); err != nil {
		master.close()
		log.Errorf("master start faild:%v", err)
		return
	}

	if _, has := eosc_args.GetEnv("MASTER_CONTINUE"); has {
		syscall.Kill(syscall.Getppid(), syscall.SIGQUIT)
	}
	master.Wait()
}

type Master struct {
	service.UnimplementedMasterServer
	service.UnimplementedCtiServiceServer

	masterTraffic traffic.IController
	workerTraffic traffic.IController
	store         eosc.IStore
	masterSrv     *grpc.Server
	ctx           context.Context
	cancelFunc    context.CancelFunc

	PID *pidfile.PidFile
}

func (m *Master) Hello(ctx context.Context, request *service.HelloRequest) (*service.HelloResponse, error) {
	return &service.HelloResponse{
		Name: request.GetName(),
	}, nil
}

func (m *Master) Start() error {
	m.InitLogTransport()
	m.masterTraffic = traffic.NewController(os.Stdin)
	m.workerTraffic = traffic.NewController(os.Stdin)

	// 设置存储操作
	s, err := store.NewStore()
	if err != nil {
		log.Error("new store error: ", err.Error())
		return err
	}
	m.store = s

	log.Info("master start grpc service")
	err = m.startService()
	if err != nil {
		log.Error("master start  grpc server error: ", err.Error())
		return err
	}

	ip := os.Getenv(fmt.Sprintf("%s_%s", process.AppName(), eosc_args.IP))
	port := os.Getenv(fmt.Sprintf("%s_%s", process.AppName(), eosc_args.Port))
	log.Info(fmt.Sprintf("%s:%s", ip, port))
	// 监听master监听地址，用于接口处理
	_, err = m.masterTraffic.ListenTcp("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		log.Error(err)
		os.Exit(1)
		return err
	}
	return nil

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

			}
		case syscall.SIGUSR1:
			{

				// TODO: 平滑重启操作
				m.Fork() //传子进程需要的内容

			}
		default:
			continue
		}
	}
}

func (m *Master) close() {

	m.cancelFunc()
	m.stopService()
	m.PID.Remove()
}

func NewMasterHandle(pid *pidfile.PidFile) *Master {

	cancel, cancelFunc := context.WithCancel(context.Background())
	return &Master{
		PID:                           pid,
		cancelFunc:                    cancelFunc,
		ctx:                           cancel,
		UnimplementedMasterServer:     service.UnimplementedMasterServer{},
		UnimplementedCtiServiceServer: service.UnimplementedCtiServiceServer{},
	}
}
