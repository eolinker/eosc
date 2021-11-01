/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package process_master

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"strings"

	"github.com/eolinker/eosc/config"

	"github.com/eolinker/eosc/env"

	"github.com/eolinker/eosc/utils"

	"github.com/eolinker/eosc"

	"github.com/eolinker/eosc/process-master/admin"
	"github.com/eolinker/eosc/process-master/professions"

	"github.com/eolinker/eosc/process-master/workers"

	raft_service "github.com/eolinker/eosc/raft/raft-service"

	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/pidfile"
	"github.com/eolinker/eosc/raft"
	"github.com/eolinker/eosc/service"
	"github.com/eolinker/eosc/traffic"

	"google.golang.org/grpc"

	"os"
	"os/signal"
	"syscall"
)

func Process() {
	utils.InitLogTransport(eosc.ProcessMaster)

	cfg, err := config.GetConfig()
	if err != nil {
		log.Error("get config error: ", err)
		return
	}
	file, err := pidfile.New()
	if err != nil {
		log.Errorf("the process-master is running:%v by:%d", err, os.Getpid())
		return
	}

	master := NewMasterHandle(file)
	if err := master.Start(nil, cfg); err != nil {
		master.close()
		log.Errorf("process-master[%d] start faild:%v", os.Getpid(), err)
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
	masterSrv     *grpc.Server
	ctx           context.Context
	cancelFunc    context.CancelFunc
	PID           *pidfile.PidFile
	httpserver    *http.Server
	cfg           *config.Config
	admin         *admin.Admin

	workerController *WorkerController
}

type MasterHandler struct {
	Professions eosc.IProfessionsData
}

func (mh *MasterHandler) initHandler() {
	if mh.Professions == nil {
		mh.Professions = professions.NewProfessions()
	}
}

func (m *Master) start(handler *MasterHandler, cfg *config.Config) error {
	if handler == nil {
		handler = new(MasterHandler)
	}

	handler.initHandler()
	s := raft_service.NewService()
	workersData := NewWorkersData(workers.NewTypedWorkers())
	professionRaft := NewProfessionRaft(handler.Professions)
	m.workerController = NewWorkerController(m.workerTraffic, professionRaft, workersData, cfg)
	m.workerController.Start()
	worker := NewWorkersRaft(workersData, handler.Professions, m.workerController, s, m.workerController)
	m.admin = admin.NewAdmin(handler.Professions, worker)
	s.SetHandlers(raft_service.NewCreateHandler(workers.SpaceWorker, worker), raft_service.NewCreateHandler(professions.SpaceProfession, professionRaft))
	node, err := raft.NewNode(s)
	if err != nil {
		log.Error(err)
		return err
	}

	m.node = node

	return nil
}

func (m *Master) Start(handler *MasterHandler, cfg *config.Config) error {

	m.cfg = cfg
	m.masterTraffic.Reset([]int{m.cfg.Admin.Listen})
	m.workerTraffic.Reset(cfg.Ports())
	err := m.start(handler, cfg)
	if err != nil {
		return err
	}

	// 监听master监听地址，用于接口处理
	l, err := m.masterTraffic.ListenTcp(m.cfg.Admin.IP, m.cfg.Admin.Listen)
	if err != nil {
		log.Error("master listen tcp error: ", err)
		return err
	}

	if strings.ToLower(m.cfg.Admin.Scheme) == "https" {
		// start https listener
		log.Debug("start https listener...")
		cert, err := config.NewCert([]*config.Certificate{m.cfg.Admin.Certificate}, m.cfg.CertificateDir.Dir)
		if err != nil {
			return err
		}
		l = tls.NewListener(l, &tls.Config{GetCertificate: cert.GetCertificate})
	}

	m.startHttp(l)
	if _, has := env.GetEnv("MASTER_CONTINUE"); has {
		syscall.Kill(syscall.Getppid(), syscall.SIGQUIT)
	}
	log.Info("process-master start grpc service")
	err = m.startService()
	if err != nil {
		log.Error("process-master start  grpc server error: ", err.Error())
		return err
	}
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
	sm.Handle("/raft", m.node)
	sm.Handle("/raft/", m.node)

	sm.Handle("/api/", m.admin)
	sm.Handle("/api", m.admin)

	return sm
}
func (m *Master) Wait() error {

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	for {
		sig := <-sigc
		log.Infof("Caught signal pid:%d ppid:%d signal %s: .\n", os.Getpid(), os.Getppid(), sig.String())
		//fmt.Println(os.Interrupt.String(), sig.String(), sig == os.Interrupt)
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
				m.node.Stop()
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

func (m *Master) Close() {
	m.close()
}

func (m *Master) close() {
	log.Info("process-master close")
	log.Info("raft node close")
	m.node.Stop()

	m.cancelFunc()
	log.Debug("process-master shutdown http:", m.httpserver.Shutdown(context.Background()))
	m.masterTraffic.Close()

	m.workerTraffic.Close()

	m.stopService()
	log.Debug("try remove pid")

	if err := m.PID.Remove(); err != nil {
		log.Warn("remove pid:", err)
	}
	m.workerController.Stop()

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
	if _, has := env.GetEnv("MASTER_CONTINUE"); has {
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
