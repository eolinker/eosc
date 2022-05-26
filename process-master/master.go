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
	"github.com/eolinker/eosc/process"
	"github.com/eolinker/eosc/utils"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/eolinker/eosc/process-master/extender"

	open_api "github.com/eolinker/eosc/process-master/open-api"

	raft_service "github.com/eolinker/eosc/process-master/raft-service"

	"github.com/eolinker/eosc/config"

	"github.com/eolinker/eosc/env"

	"github.com/eolinker/eosc"

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
	node          *raft.Node
	masterTraffic traffic.IController
	workerTraffic traffic.IController
	raftService   raft.IRaftService
	masterSrv     *grpc.Server
	ctx           context.Context
	cancelFunc    context.CancelFunc

	//PID           *pidfile.PidFile
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

func (m *Master) start(handler *MasterHandler, cfg *config.Config) error {

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
	m.workerController = NewWorkerController(m.workerTraffic, cfg, process.NewProcessController(m.ctx, eosc.ProcessWorker, m.logWriter))

	m.dispatcherServe = NewDispatcherServer()
	extenderManager := extender.NewManager(m.ctx, extender.GenCallbackList(m.dispatcherServe, m.workerController))
	m.dataController = NewDataController(raftService, extenderManager, m.dispatcherServe)

	node, err := raft.NewNode(raftService, m.adminController)
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

	m.startOpenApi(l)

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
func (m *Master) startOpenApi(l net.Listener) {
	sm := open_api.NewOpenApiProxy(m.node, m.openApiProxy)
	sm.ExcludeHandlers("/raft", m.node)

	m.httpserver = &http.Server{
		Handler: sm,
	}

	go func() {
		err := m.httpserver.Serve(l)
		if err != nil {
			log.Warn(err)
		}
	}()
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
