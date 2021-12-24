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
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/eolinker/eosc/process-master/extenders"

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
	"github.com/eolinker/eosc/traffic"

	"google.golang.org/grpc"

	"os"
	"os/signal"
	"syscall"
)

func Process() {
	logWriter := utils.InitLogTransport(eosc.ProcessMaster)
	pFile, err := pidfile.New()
	if err != nil {
		log.Errorf("the process-master is running:%v by:%d", err, os.Getpid())
		return
	}
	master := NewMasterHandle(logWriter)
	if err := master.Start(nil); err != nil {
		master.close()
		log.Errorf("process-master[%d] start faild:%v", os.Getpid(), err)
		return
	}

	master.Wait(pFile)
	pFile.Remove()
}

type Master struct {
	//service.UnimplementedMasterServer
	node          *raft.Node
	masterTraffic traffic.IController
	workerTraffic traffic.IController
	raftService   raft.IService
	masterSrv     *grpc.Server
	ctx           context.Context
	cancelFunc    context.CancelFunc

	//PID           *pidfile.PidFile
	httpserver          *http.Server
	logWriter           io.Writer
	admin               *admin.Admin
	extenderSettingRaft *ExtenderSettingRaft
	workerController    *WorkerController
}

type MasterHandler struct {
	Professions eosc.IProfessions
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
	workerServiceProxy := NewWorkerServiceProxy()
	raftService := raft_service.NewService()

	workersConfig := NewWorkerConfigs()

	professionRaft := NewProfessionRaft(handler.Professions)
	extenderRaft := NewExtenderRaft(raftService)
	m.extenderSettingRaft = extenderRaft
	workerRaft := NewWorkersRaft(workersConfig, handler.Professions, workerServiceProxy, raftService)

	m.workerController = NewWorkerController(m.workerTraffic, cfg, extenderRaft.data, handler.Professions, workersConfig, workerServiceProxy, m.logWriter)

	m.admin = admin.NewAdmin(handler.Professions, workerRaft)
	raftService.AddEventHandler(m.workerController.raftEvent)
	raftService.AddCommitEventHandler(m.workerController.raftCommitEvent)

	raftService.SetHandlers(
		raft_service.NewCreateHandler(workers.SpaceWorker, workerRaft),
		raft_service.NewCreateHandler(professions.SpaceProfession, professionRaft),
		raft_service.NewCreateHandler(extenders.NamespaceExtenders, extenderRaft),
	)
	node, err := raft.NewNode(raftService)
	if err != nil {
		log.Error(err)
		return err
	}

	m.node = node

	return nil
}

func (m *Master) Start(handler *MasterHandler) error {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Error("get config error: ", err)
		return err
	}
	_, err = m.masterTraffic.Reset([]int{cfg.Admin.Listen})
	if err != nil {
		return err
	}
	_, err = m.workerTraffic.Reset(cfg.Ports())
	if err != nil {
		return err
	}
	err = m.start(handler, cfg)
	if err != nil {
		return err
	}

	// 监听master监听地址，用于接口处理
	l, err := m.masterTraffic.ListenTcp(cfg.Admin.IP, cfg.Admin.Listen)
	if err != nil {
		log.Error("master listen tcp error: ", err)
		return err
	}

	if strings.ToLower(cfg.Admin.Scheme) == "https" {
		// start https listener
		log.Debug("start https listener...")
		cert, err := config.NewCert([]*config.Certificate{cfg.Admin.Certificate}, cfg.CertificateDir.Dir)
		if err != nil {
			return err
		}
		l = tls.NewListener(l, &tls.Config{GetCertificate: cert.GetCertificate})
	}

	m.startHttp(l)

	log.Info("process-master start grpc service")
	err = m.startService()
	if err != nil {
		log.Error("process-master start  grpc server error: ", err.Error())
		return err
	}
	m.workerController.WaitStart()
	if _, has := env.GetEnv("MASTER_CONTINUE"); has {
		syscall.Kill(syscall.Getppid(), syscall.SIGQUIT)
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

	sm.Handle("/extender/", m.extenderSettingRaft)
	sm.Handle("/extenders", m.extenderSettingRaft)

	return sm
}
func (m *Master) Wait(pFile *pidfile.PidFile) error {

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
				err := m.Fork(pFile) //传子进程需要的内容
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
	if m.node == nil {
		return
	}
	log.Info("process-master close")
	log.Info("raft node close")
	m.node.Stop()

	m.cancelFunc()
	log.Debug("process-master shutdown http:", m.httpserver.Shutdown(context.Background()))
	m.masterTraffic.Close()

	m.workerTraffic.Close()

	m.stopService()
	log.Debug("try remove pid")

	m.workerController.Stop()

}

func NewMasterHandle(logWriter io.Writer) *Master {
	cancel, cancelFunc := context.WithCancel(context.Background())
	m := &Master{
		cancelFunc: cancelFunc,
		ctx:        cancel,
		logWriter:  logWriter,
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
