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
	"net"
	"net/http"
	"strconv"

	raft_service "github.com/eolinker/eosc/raft/raft-service"

	eosc_args "github.com/eolinker/eosc/eosc-args"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/pidfile"
	"github.com/eolinker/eosc/raft"
	"github.com/eolinker/eosc/service"
	"github.com/eolinker/eosc/store"

	"github.com/eolinker/eosc/traffic"

	"google.golang.org/grpc"

	"os"
	"os/signal"
	"syscall"
)

func Process() {
	InitLogTransport()
	file, err := pidfile.New()
	if err != nil {
		log.Errorf("the master is running:%v by:%d", err, os.Getpid())
		return
	}
	master := NewMasterHandle(file)
	if err := master.Start(); err != nil {
		master.close()
		log.Errorf("master[%d] start faild:%v", os.Getpid(), err)
		return
	}
	if _, has := eosc_args.GetEnv("MASTER_CONTINUE"); has {
		syscall.Kill(syscall.Getppid(), syscall.SIGQUIT)
	}
	log.Info("master start grpc service")
	err = master.startService()
	if err != nil {
		log.Error("master start  grpc server error: ", err.Error())
		return
	}
	master.Wait()
}

type Master struct {
	service.UnimplementedMasterServer
	service.UnimplementedCtiServiceServer
	node          *raft.Node
	masterTraffic traffic.IController
	workerTraffic traffic.IController
	raftService   raft.IService
	//store         eosc.IStore
	masterSrv  *grpc.Server
	ctx        context.Context
	cancelFunc context.CancelFunc
	PID        *pidfile.PidFile

	httpserver *http.Server
}

func (m *Master) Start() error {
	// 设置存储操作
	store, err := store.NewStore()
	if err != nil {
		log.Error("new store error: ", err.Error())
		return err
	}
	//m.store = s
	m.raftService = raft_service.NewService(store)

	m.node, _ = raft.NewNode(m.raftService)

	ip := eosc_args.GetDefault(eosc_args.IP, "")
	port, _ := strconv.Atoi(eosc_args.GetDefault(eosc_args.Port, "9400"))
	// 监听master监听地址，用于接口处理
	l, err := m.masterTraffic.ListenTcp(ip, port)
	if err != nil {
		log.Error(err)
		return err
	}

	m.startHttp(l)

	return nil

}
func (m *Master) startHttp(l net.Listener) {
	m.httpserver = &http.Server{
		Handler: m.handler(),
	}
	go func() {
		err := m.httpserver.Serve(l)
		if err != nil {
			log.Warn(err)
		}
	}()
}
func (m *Master) handler() http.Handler {
	sm := http.NewServeMux()
	sm.Handle("/raft/", m.node.Handler())

	return sm
}
func (m *Master) Wait() error {

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	for {
		sig := <-sigc
		log.Infof("Caught signal pid:%d ppid:%d signal %s: .\n", os.Getpid(), os.Getppid(), sig.String())

		switch sig {
		case os.Interrupt, os.Kill:
			{
				m.close()
				return nil
			}
		case syscall.SIGQUIT:
			{

				m.close()
				return nil
			}
		case syscall.SIGUSR1:
			{

				log.Info("try fork new")
				err := m.Fork() //传子进程需要的内容
				if err != nil {
					log.Error("fork new:", err)
				}
			}
		default:
			continue
		}
	}
}

func (m *Master) close() {
	log.Info("master close")
	m.cancelFunc()
	log.Debug("master shutdown http:", m.httpserver.Shutdown(context.Background()))

	m.masterTraffic.Close()
	m.workerTraffic.Close()
	m.stopService()
	log.Debug("try remove pid")

	if err := m.PID.Remove(); err != nil {
		log.Warn("remove pid:", err)
	}

}

func NewMasterHandle(pid *pidfile.PidFile) *Master {

	cancel, cancelFunc := context.WithCancel(context.Background())
	m := &Master{
		PID:                           pid,
		cancelFunc:                    cancelFunc,
		ctx:                           cancel,
		UnimplementedMasterServer:     service.UnimplementedMasterServer{},
		UnimplementedCtiServiceServer: service.UnimplementedCtiServiceServer{},
	}
	if _, has := eosc_args.GetEnv("MASTER_CONTINUE"); has {
		log.Info("init traffic from stdin")
		m.masterTraffic = traffic.NewController(os.Stdin)
		m.workerTraffic = traffic.NewController(os.Stdin)
	} else {
		log.Info("new traffic")

		m.masterTraffic = traffic.NewController(nil)
		m.workerTraffic = traffic.NewController(nil)
	}
	return m
}
