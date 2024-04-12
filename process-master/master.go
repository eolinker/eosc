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
	"fmt"
	"github.com/eolinker/eosc/process-master/proxy"
	"github.com/soheilhy/cmux"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/eolinker/eosc/etcd"
	"github.com/eolinker/eosc/process"
	"github.com/eolinker/eosc/process-master/extender"
	raft_service "github.com/eolinker/eosc/process-master/raft-service"
	"github.com/eolinker/eosc/router"
	"github.com/eolinker/eosc/traffic/mixl"
	"github.com/eolinker/eosc/utils"

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
	logWriter, logHandler := utils.InitMasterLog()
	handler.logHandler = logHandler
	log.Debug("master start:", os.Getpid(), ":", os.Getppid())

	pFile, err := pidfile.New()
	if err != nil {
		log.Errorf("the process-master is running:%v by:%d", err, os.Getpid())
		return
	}
	cfg := config.Load()
	log.Debug(cfg)
	master, err := NewMasterHandle(logWriter, cfg)
	if err != nil {
		log.Errorf("process-master[%d] start faild:%v", os.Getpid(), err)
		return
	}
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
	config           config.NConfig
	etcdServer       etcd.Etcd
	adminTraffic     *traffic.TrafficData
	workerTraffic    *traffic.TrafficData
	masterSrv        *grpc.Server
	ctx              context.Context
	cancelFunc       context.CancelFunc
	httpserver       []*http.Server
	logWriter        io.Writer
	dataController   *DataController
	workerController *WorkerController
	adminController  *AdminController
	dispatcherServe  *DispatcherServer
	adminUnixProxy   *proxy.UnixProxy
	workerUnixProxy  *proxy.UnixProxy
}

type MasterHandler struct {
	InitProfession func() []*eosc.ProfessionConfig
	VersionHandler func(etcd2 etcd.Etcd) http.Handler
	logHandler     func(prefix string) http.Handler
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

func (m *Master) start(handler *MasterHandler, etcdServer etcd.Etcd) error {

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
	m.adminUnixProxy = proxy.NewUnixProxy(eosc.ProcessAdmin)
	m.workerUnixProxy = proxy.NewUnixProxy(eosc.ProcessWorker)
	m.adminController = NewAdminConfig(raftService, process.NewProcessController(m.ctx, eosc.ProcessAdmin, m.logWriter, m.adminUnixProxy))
	m.workerController = NewWorkerController(m.workerTraffic, m.config.Gateway, process.NewProcessController(m.ctx, eosc.ProcessWorker, m.logWriter, m.workerUnixProxy))

	m.dispatcherServe = NewDispatcherServer()
	extenderManager := extender.NewManager(m.ctx, extender.GenCallbackList(m.dispatcherServe, m.workerController))
	m.dataController = NewDataController(raftService, extenderManager, m.dispatcherServe)

	etcdServer.Watch("/", raftService)
	etcdServer.HandlerLeader(m.adminController)

	return nil
}

func (m *Master) Start(handler *MasterHandler) error {
	for _, uri := range m.config.Gateway.AdvertiseUrls {
		// 设置网关IP
		u, err := url.Parse(uri)
		if err != nil {
			log.Errorf("parse advertise url(%s) error:", uri, err)
			continue
		}
		index := strings.Index(u.Host, ":")
		host := u.Host
		if index >= 0 {
			host = u.Host[:index]
		}
		os.Setenv("gateway_ip", host)
		break
	}
	peerListener, err := m.listen(m.config.Peer)

	if err != nil {
		return err
	}

	clientListener, err := m.listen(m.config.Client)
	if err != nil {
		return err
	}

	/*
	   peer listener 监听后, etcd 的http请求由etcd server 处理, 其他的请求视为open api,
	   1. 如果当前节点是leader , 转发给 admin
	   2. 如果当前节点不是leader ,转发给 leader 的peer
	*/
	/*
		client listener 监听后:
		1. /system, /apinto/log/node/, 由master处理
		2. /apinto/ 转发给 worker处理
		3. 其他请求视为 open api,
			如当前节点是leader,转给 admin
			如果当前节点不是leader,转给 leader的peer
	*/

	peerListeners := utils.MatchMux(peerListener, etcdPaths)
	clienListeners := utils.MatchMux(clientListener, masterApiPaths, workerApiPaths)
	// 初始化etcd
	etcdApiListener := peerListeners[0]
	peerOpenApiListener := peerListeners[1]
	masterApiListener := clienListeners[0]
	workerApiListener := clienListeners[1]
	clientOpenApiListener := clienListeners[2]
	etcdMux := m.startHttpServer(etcdApiListener)
	etcdServer, err := etcd.NewServer(m.ctx, etcdMux, etcd.Config{
		PeerAdvertiseUrls:    m.config.Peer.AdvertiseUrls,
		ClientAdvertiseUrls:  m.config.Client.AdvertiseUrls,
		GatewayAdvertiseUrls: m.config.Gateway.AdvertiseUrls,
		DataDir:              env.DataDir(),
	})
	if err != nil {
		log.Error("start etcd error:", err)
		return err
	}

	m.etcdServer = etcdServer
	// 初始化etcd数据
	err = m.start(handler, etcdServer)
	if err != nil {
		return err
	}

	//openApiProxy := open_api.NewOpenApiProxy(m.etcdServer, m.adminUnixProxy)
	//etcdMux.Handle("/", http.NotFoundHandler()) // etcdMux 只会接受到指定前缀的请求

	masterMux := m.startHttpServer(masterApiListener)
	masterMux.Handle("/system/version", handler.VersionHandler(etcdServer))
	masterMux.HandleFunc("/system/info", m.EtcdInfoHandler)
	masterMux.HandleFunc("/system/nodes", m.EtcdNodesHandler)
	//node log
	logPrefix := fmt.Sprintf("%slog/node/", router.RouterPrefix)
	masterMux.Handle(logPrefix, handler.logHandler(logPrefix)) //master处理本地日志

	//masterMux.Handle(router.RouterPrefix, m.workerUnixProxy) //master转发至worker的路由
	//master转发至worker
	go doServer(workerApiListener, m.workerUnixProxy.ProxyToUnix)
	// 转发到leader -> admin
	go doServer(clientOpenApiListener, proxy.ProxyToLeader(m.etcdServer, m.adminUnixProxy))
	// 转发到 leader -> admin
	// todo 后续需要处理 peer 只有 leader 才能处理 open api, 否则根据协议报error , 以避免消息循环
	go doServer(peerOpenApiListener, proxy.ProxyToLeader(m.etcdServer, m.adminUnixProxy))

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

func (m *Master) listen(conf config.UrlConfig) (net.Listener, error) {
	tf := traffic.NewTraffic(m.adminTraffic)
	tcp, ssl := tf.Listen(conf.ListenUrls...)

	listener := make([]net.Listener, 0, len(tcp)+len(ssl))
	listener = append(listener, tcp...)
	if len(ssl) > 0 {
		cert, err := config.LoadCert(conf.Certificate, m.config.CertificateDir.Dir)
		if err != nil {
			return nil, err
		}
		tlsConf := &tls.Config{GetCertificate: cert.GetCertificate}
		for _, l := range ssl {
			listener = append(listener, tls.NewListener(l, tlsConf))
		}
	}
	if len(listener) == 0 {
		return nil, errors.New("no client listener")
	}
	return mixl.NewMixListener(0, listener...), nil
}

func cmuxListener(ln net.Listener) (httpListener net.Listener, apintoListener net.Listener) {
	cm := cmux.New(ln)

	httpListener = cm.Match(cmux.HTTP1Fast(), cmux.HTTP2())
	apintoListener = cm.Match(cmux.Any())
	go func() {
		cm.Serve()
	}()
	return httpListener, apintoListener
}

func (m *Master) startHttpServer(listen net.Listener) *http.ServeMux {
	mux := http.NewServeMux()

	server := &http.Server{
		Handler: mux,
	}
	// todo 暂时关闭keepalive
	server.SetKeepAlivesEnabled(false)
	go func() {
		err := server.Serve(listen)
		if err != nil {
			log.Warn(err)
		}
	}()
	m.httpserver = append(m.httpserver, server)
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

	httpservers := m.httpserver
	m.httpserver = nil
	for _, server := range httpservers {
		log.Debug("process-master shutdown http:", server.Shutdown(context.Background()))

	}

	m.adminTraffic.Close()
	m.workerTraffic.Close()
	m.dispatcherServe.Close()
	m.dataController.Close()
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
		config:     cfg,
	}
	var input io.Reader
	if _, has := env.GetEnv("MASTER_CONTINUE"); has {
		log.Info("Reset traffic from stdin")

		input = os.Stdin
	} else {
		input = nil
	}
	masterTraffic, err := traffic.ReadTraffic(input, config.GetListens(cfg.Client.ListenUrl, cfg.Peer.ListenUrl)...)
	if err != nil {
		return nil, err
	}
	m.adminTraffic = masterTraffic

	workerTraffic, err := traffic.ReadTraffic(input, config.GetListens(cfg.Gateway)...)
	if err != nil {
		return nil, err
	}
	m.workerTraffic = workerTraffic
	return m, nil
}
