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
	"github.com/eolinker/eosc/etcd"
	"github.com/eolinker/eosc/process"
	open_api "github.com/eolinker/eosc/process-master/open-api"
	"github.com/eolinker/eosc/utils"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/eolinker/eosc/process-master/extender"

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
	pFile, err := pidfile.New()
	if err != nil {
		log.Errorf("the process-master is running:%v by:%d", err, os.Getpid())
		return
	}

	master := NewMasterHandle(logWriter)
	if err := master.Start(handler); err != nil {
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
	masterTraffic    traffic.IController
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
	openApiProxy     *UnixAdminProcess
}

type MasterHandler struct {
	InitProfession func() []*eosc.ProfessionConfig
}

func (mh *MasterHandler) initHandler() {

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
	}, func(config map[string]map[string][]byte) map[string]map[string][]byte {
		if config == nil {
			config = make(map[string]map[string][]byte)
		}
		if ps, has := config[eosc.NamespaceProfession]; has {
			pl := make([]*eosc.ProfessionConfig, 0, len(ps))
			for _, d := range ps {
				p := new(eosc.ProfessionConfig)
				if err := json.Unmarshal(d, p); err != nil {
					continue
				}
				pl = append(pl, p)
			}
			initWs := eosc.GenInitWorkerConfig(pl)
			ws, wsHas := config[eosc.NamespaceWorker]
			if !wsHas {
				ws = make(map[string][]byte)
			}
			for _, w := range initWs {
				if _, has := ws[w.Id]; !has {
					ws[w.Id], _ = json.Marshal(w)
				}
			}
			config[eosc.NamespaceWorker] = ws
		}
		return config
	})

	m.openApiProxy = NewUnixClient()
	m.adminController = NewAdminConfig(raftService, process.NewProcessController(m.ctx, eosc.ProcessAdmin, m.logWriter, m.openApiProxy))
	m.workerController = NewWorkerController(m.workerTraffic, listensMsg, process.NewProcessController(m.ctx, eosc.ProcessWorker, m.logWriter))

	m.dispatcherServe = NewDispatcherServer()
	extenderManager := extender.NewManager(m.ctx, extender.GenCallbackList(m.dispatcherServe, m.workerController))
	m.dataController = NewDataController(raftService, extenderManager, m.dispatcherServe)

	etcdServer.Watch("/", raftService)
	etcdServer.HandlerLeader(m.adminController)

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
	mux := m.startHttpServer(l)
	etcdServer, err := etcd.NewServer(m.ctx, mux)
	if err != nil {
		log.Error("start etcd error:", err)
		return err
	}
	m.etcdServer = etcdServer
	err = m.start(handler, cfg.Export(), etcdServer)
	if err != nil {
		return err
	}
	sm := open_api.NewOpenApiProxy(NewEtcdSender(m.etcdServer), m.openApiProxy)
	mux.Handle("/", sm)
	log.Info("process-master start grpc service")
	err = m.startService()
	if err != nil {
		log.Error("process-master start  grpc server error: ", err.Error())
		return err
	}

	if _, has := env.GetEnv("MASTER_CONTINUE"); has {
		syscall.Kill(syscall.Getppid(), syscall.SIGQUIT)
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
				m.etcdServer.Close()
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
	if m.etcdServer == nil {
		return
	}
	log.Info("process-master close")
	log.Info("raft node close")
	m.etcdServer.Close()

	m.cancelFunc()
	log.Debug("process-master shutdown http:", m.httpserver.Shutdown(context.Background()))

	m.masterTraffic.Close()
	m.workerTraffic.Close()
	m.dispatcherServe.Close()

	m.stopService()
	log.Debug("try remove pid")

	m.workerController.Stop()
	m.adminController.Stop()
}

func NewMasterHandle(logWriter io.Writer) *Master {
	cancel, cancelFunc := context.WithCancel(context.Background())
	m := &Master{
		cancelFunc: cancelFunc,
		ctx:        cancel,
		logWriter:  logWriter,
	}
	if _, has := env.GetEnv("MASTER_CONTINUE"); has {
		log.Info("Reset traffic from stdin")
		m.masterTraffic = traffic.NewController(os.Stdin)
		m.workerTraffic = traffic.NewController(os.Stdin)
	} else {
		log.Info("new traffic")

		m.masterTraffic = traffic.NewController(nil)
		m.workerTraffic = traffic.NewController(nil)
	}
	return m
}
