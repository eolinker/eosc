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
	"encoding/json"
	"errors"
	"github.com/eolinker/eosc/etcd"
	"github.com/eolinker/eosc/process"
	"github.com/eolinker/eosc/process-master/extender"
	open_api "github.com/eolinker/eosc/process-master/open-api"
	"github.com/eolinker/eosc/utils"
	"io"
	"net"
	"net/http"
	"strings"

	raft_service "github.com/eolinker/eosc/process-master/raft-service"

	"github.com/eolinker/eosc/config"

	"github.com/eolinker/eosc/env"

	"github.com/eolinker/eosc"

	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/pidfile"
	"github.com/eolinker/eosc/traffic"

	"google.golang.org/grpc"

	"os"
	"os/signal"
	"syscall"
)

func Process() {
	ProcessDo(nil)
}
func ProcessDo(handler *MasterHandler) {
	logWriter := utils.InitMasterLog()
	log.Debug("master start:", os.Getpid(), ":", os.Getppid())

	pFile, err := pidfile.New()
	if err != nil {
		log.Errorf("the process-master is running:%v by:%d", err, os.Getpid())
		return
	}
	cfg := config.Load()

	master, err := NewMasterHandle(logWriter, cfg)
	if err != nil {
		log.Errorf("process-master[%d] start faild:%v", os.Getpid(), err)
		return
	}
	if err := master.Start(handler, cfg); err != nil {
		master.close()
		log.Errorf("process-master[%d] start faild:%v", os.Getpid(), err)
		return
	}

	master.Wait(pFile)
	pFile.Remove()
}

type Master struct {
	//service.UnimplementedMasterServer
	etcdServer       etcd.Etcd
	adminTraffic     traffic.IController
	workerTraffic    traffic.IController
	masterSrv        *grpc.Server
	ctx              context.Context
	cancelFunc       context.CancelFunc
	httpserver       *http.Server
	logWriter        io.Writer
	dataController   *DataController
	workerController *WorkerController
	adminController  *AdminController
	dispatcherServe  *DispatcherServer
	adminClient      *UnixClient
}

type MasterHandler struct {
	InitProfession func() []*eosc.ProfessionConfig
	VersionHandler func(etcd2 etcd.Etcd) http.Handler
}

func (mh *MasterHandler) initHandler() {
	if mh.VersionHandler == nil {
		mh.VersionHandler = func(server etcd.Etcd) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				version := server.Version()
				json.NewEncoder(w).Encode(version)
			})
		}
	}
}

func (m *Master) start(handler *MasterHandler, listensMsg *config.ListensMsg, etcdServer etcd.Etcd) error {

	if handler == nil {
		handler = new(MasterHandler)
	}

	handler.initHandler()

	raftService := raft_service.NewService(func(config map[string]map[string][]byte) map[string]map[string][]byte {
		if config == nil {
			config = make(map[string]map[string][]byte)
		}
		if handler.InitProfession != nil {

			if ps, has := config[eosc.NamespaceProfession]; !has || len(ps) == 0 {
				ps = make(map[string][]byte)
				pl := handler.InitProfession()
				for _, p := range pl {
					data, _ := json.Marshal(p)
					ps[p.Name] = data
				}
				config[eosc.NamespaceProfession] = ps
			}
		}
		return config
	})

	m.adminController = NewAdminConfig(raftService, process.NewProcessController(m.ctx, eosc.ProcessAdmin, m.logWriter, m.adminClient))
	m.workerController = NewWorkerController(m.workerTraffic, listensMsg, process.NewProcessController(m.ctx, eosc.ProcessWorker, m.logWriter))

	m.dispatcherServe = NewDispatcherServer()
	extenderManager := extender.NewManager(m.ctx, extender.GenCallbackList(m.dispatcherServe, m.workerController))
	m.dataController = NewDataController(raftService, extenderManager, m.dispatcherServe)

	etcdServer.Watch("/", raftService)
	etcdServer.HandlerLeader(m.adminController)

	return nil
}

func (m *Master) Start(handler *MasterHandler, cfg *config.NConfig) error {

	// 监听master监听地址，用于接口处理
	l := m.adminTraffic.ListenTcp(cfg.Admin.Listen, traffic.Http1)
	if l == nil {
		log.Error("master listen tcp error: ")
		return errors.New("not allow")
	}

	if strings.ToLower(cfg.Admin.Scheme) == "https" {
		// start https listener
		log.Debug("start https listener...")
		cert, err := config.LoadCert([]*config.Certificate{cfg.Admin.Certificate}, cfg.CertificateDir.Dir)
		if err != nil {
			return err
		}
		l = tls.NewListener(l, &tls.Config{GetCertificate: cert.GetCertificate})
	}
	mux := m.startHttpServer(l)
	etcdServer, err := etcd.NewServer(m.ctx, mux)
	if err != nil {
		log.Error("start etcd error:", err)
		return err
	}
	m.adminClient = NewUnixClient()
	m.etcdServer = etcdServer
	err = m.start(handler, cfg.Export(), etcdServer)
	if err != nil {
		return err
	}
	openApiProxy := open_api.NewOpenApiProxy(NewEtcdSender(m.etcdServer), m.adminClient)

	mux.Handle("/system/version", handler.VersionHandler(etcdServer))
	mux.HandleFunc("/system/info", m.EtcdInfoHandler)
	mux.HandleFunc("/system/nodes", m.EtcdNodesHandler)
	mux.Handle("/", openApiProxy)
	log.Info("process-master start grpc service")
	err = m.startService()
	if err != nil {
		log.Error("process-master start  grpc server error: ", err.Error())
		return err
	}

	if _, has := env.GetEnv("MASTER_CONTINUE"); has {
		log.Debug("master continue: call parent:", os.Getppid())
		syscall.Kill(os.Getppid(), syscall.SIGQUIT)
	}
	return nil

}

func (m *Master) startHttpServer(l net.Listener) *http.ServeMux {
	mux := http.NewServeMux()
	m.httpserver = &http.Server{
		Handler: mux,
	}
	go func() {
		err := m.httpserver.Serve(l)
		if err != nil {
			log.Warn(err)
		}
	}()
	return mux
}

func (m *Master) Wait(pFile *pidfile.PidFile) error {

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)

	for {
		sig := <-sigc
		log.Infof("Caught signal pid:%d ppid:%d signal %s: .\n", os.Getpid(), os.Getppid(), sig.String())
		//log.Debug(os.Interrupt.String(), sig.String(), sig == os.Interrupt)
		switch sig {
		case os.Interrupt, os.Kill:
			{
				if m.etcdServer != nil {
					m.etcdServer.Close()
					m.etcdServer = nil
				}
				m.close()
				return nil
			}
		case syscall.SIGQUIT:
			{
				if m.etcdServer != nil {
					m.etcdServer.Close()
					m.etcdServer = nil
				}
				m.close()
				return nil
			}
		case syscall.SIGUSR1:
			{
				if m.etcdServer != nil {
					m.etcdServer.Close()
					m.etcdServer = nil
				}
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
	if m.cancelFunc == nil {
		return
	}

	m.cancelFunc()
	m.cancelFunc = nil
	log.Debug("process-master shutdown http:", m.httpserver.Shutdown(context.Background()))

	m.adminTraffic.Close()
	m.workerTraffic.Close()
	m.dispatcherServe.Close()

	m.stopService()
	log.Debug("try remove pid")

	m.workerController.Stop()
	m.adminController.Stop()
}

func NewMasterHandle(logWriter io.Writer, cfg config.NConfig) (*Master, error) {

	cancel, cancelFunc := context.WithCancel(context.Background())
	m := &Master{
		cancelFunc: cancelFunc,
		ctx:        cancel,
		logWriter:  logWriter,
	}
	var input io.Reader
	if _, has := env.GetEnv("MASTER_CONTINUE"); has {
		log.Info("Reset traffic from stdin")

		input = os.Stdin
	} else {
		input = nil
	}
	masterTraffic, err := traffic.ReadController(input, config.GetListens(cfg.Client, cfg.Peer)...)
	if err != nil {
		return nil, err
	}
	m.adminTraffic = masterTraffic

	workerTraffic, err := traffic.ReadController(input, config.GetListens(cfg.Gateway)...)
	if err != nil {
		return nil, err
	}
	m.workerTraffic = workerTraffic
	return m, nil
}
