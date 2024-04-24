/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package process_admin

import (
	"context"
	"encoding/json"
	admin "github.com/eolinker/eosc/process-admin/admin"
	api_apinto "github.com/eolinker/eosc/process-admin/api-apinto"
	"github.com/eolinker/eosc/process-admin/api-http"
	"github.com/eolinker/eosc/process-admin/data"
	"github.com/soheilhy/cmux"
	"time"

	"github.com/eolinker/eosc/config"
	grpc_unixsocket "github.com/eolinker/eosc/grpc-unixsocket"
	open_api "github.com/eolinker/eosc/open-api"
	"github.com/eolinker/eosc/professions"
	"github.com/eolinker/eosc/require"
	"github.com/eolinker/eosc/service"
	"github.com/eolinker/eosc/traffic"
	"github.com/eolinker/eosc/variable"

	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/julienschmidt/httprouter"

	"github.com/eolinker/eosc/process"

	"github.com/eolinker/eosc/extends"

	"github.com/eolinker/eosc/common/bean"

	"google.golang.org/protobuf/proto"

	"github.com/eolinker/eosc/utils"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
)

func Process() {
	//time.Sleep(time.Second)
	//utils.InitStdTransport(eosc.ProcessAdmin)
	log.Info("admin process start...")

	arg := readConfig()
	if arg == nil {
		arg = map[string]map[string][]byte{}
	}

	log.Debug("create admin process...")

	w, err := NewProcessAdmin(context.Background(), arg)
	if err != nil {

		w.writeOutput(process.StatusExit, err.Error())
		log.Error("new process admin error: ", err)
		return
	}
	w.writeOutput(process.StatusRunning, "")
	w.wait()
	log.Info("admin process end")
}

type ProcessAdmin struct {
	ctx    context.Context
	once   sync.Once
	reg    eosc.IExtenderDriverRegister
	router *httprouter.Router

	apiLocker    sync.Mutex
	server       *http.Server
	cx           cmux.CMux
	apintoServer *api_apinto.Server
}

func (pa *ProcessAdmin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pa.apiLocker.Lock()
	defer pa.apiLocker.Unlock()
	pa.router.ServeHTTP(w, r)
}

func (pa *ProcessAdmin) writeOutput(status int, msg string) {
	data := new(eosc.ProcessStatus)
	data.Status = int32(status)
	data.Msg = msg
	d, _ := proto.Marshal(data)
	err := utils.WriteFrame(os.Stdout, d)
	if err != nil {
		log.Error("write output error: ", err)
	}
}

func (pa *ProcessAdmin) wait() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	for {
		sig := <-sigc
		log.Infof("Caught signal pid:%d ppid:%d signal %s: .\n", os.Getpid(), os.Getppid(), sig.String())
		switch sig {
		case os.Interrupt, os.Kill:
			{
				pa.close()
				return
			}
		case syscall.SIGQUIT:
			{
				pa.close()
				return
			}
		case syscall.SIGUSR1:
			{
				pa.close()
				return
			}
		default:
			continue
		}
	}

}

// NewProcessAdmin 创建新的admin进程
// 启动时通过stdin传输配置信息
func NewProcessAdmin(parent context.Context, arg map[string]map[string][]byte) (*ProcessAdmin, error) {

	var tf traffic.ITraffic = traffic.NewEmptyTraffic()
	bean.Injection(&tf)
	var listenUrl = new(config.ListenUrl)
	bean.Injection(&listenUrl)
	register := initExtender(arg[eosc.NamespaceExtender])
	var extenderDrivers eosc.IExtenderDrivers = register
	bean.Injection(&extenderDrivers)

	p := &ProcessAdmin{
		ctx:    parent,
		router: httprouter.New(),
		server: &http.Server{},
	}
	p.server.SetKeepAlivesEnabled(false)
	p.server.Handler = p
	extenderRequire := require.NewRequireManager()
	extenderData := data.NewExtenderData(arg[eosc.NamespaceExtender], extenderRequire)

	ps := professions.NewProfessions(register)

	ps = NewProfessionsRequire(ps, extenderRequire)

	ps.Reset(utils.MapValue(utils.MapType(arg[eosc.NamespaceProfession], func(k string, v []byte) (*eosc.ProfessionConfig, bool) {
		c := new(eosc.ProfessionConfig)
		err := json.Unmarshal(v, c)
		if err != nil {
			log.Error("read profession config:", err)
			return nil, false
		}
		return c, true
	})))

	vd := variable.NewVariables(arg[eosc.NamespaceVariable])
	bean.Injection(&vd)
	wd := admin.NewImlAdminData(arg[eosc.NamespaceWorker], ps, vd, arg[eosc.NamespaceCustomer])
	p.apintoServer = api_apinto.NewServer(wd)

	var iWorkers eosc.IWorkers = wd
	bean.Injection(&iWorkers)

	_ = bean.Check()
	api_http.NewExtenderOpenApi(extenderData).Register(p.router)

	settingApi := api_http.NewSettingApi(wd)

	// openAPI handler register
	api_http.NewProfessionApi(wd).Register(p.router)
	api_http.NewWorkerApi(wd, settingApi).Register(p.router)
	settingApi.RegisterSetting(p.router)
	api_http.NewExportApi(extenderData, wd).Register(p.router)
	api_http.NewHashApi(wd).Register(p.router)
	api_http.NewVariableApi(wd).Register(p.router)

	p.router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := &open_api.Response{
			StatusCode: 404,
			Header:     nil,
			Data:       nil,
		}

		data, _ := json.Marshal(response)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(200)
		w.Write(data)
	})
	p.OpenApiServer()

	return p, nil
}

func (pa *ProcessAdmin) close() {
	pa.once.Do(func() {

		timeout, cancel := context.WithTimeout(pa.ctx, time.Second*3)
		err := pa.server.Shutdown(timeout)
		defer cancel()
		if err != nil {
			log.Warn("shutdown server error: ", err)
			return
		}
	})
}
func initExtender(config map[string][]byte) extends.IExtenderRegister {
	register := extends.InitRegister()
	extenderConfig := make(map[string]string)
	for k, v := range config {
		extenderConfig[k] = string(v)
	}
	extends.LoadPlugins(extenderConfig, register)
	return register
}

func readConfig() map[string]map[string][]byte {
	conf := make(map[string]map[string][]byte)

	data, err := utils.ReadFrame(os.Stdin)
	if err != nil {
		log.Warn("read arg fail:", err)
		return conf
	}
	err = json.Unmarshal(data, &conf)
	if err != nil {
		log.Warn("unmarshal arg fail:", err)
	}
	log.Debug("read arg:")
	for namespace, vs := range conf {
		for k, v := range vs {
			log.DebugF("read:[%s:%s]=%s\n", namespace, k, string(v))
		}
	}
	return conf
}

func (pa *ProcessAdmin) OpenApiServer() error {

	addr := service.ServerAddr(os.Getpid(), eosc.ProcessAdmin)
	syscall.Unlink(addr)
	log.Info("start admin unix server: ", addr)
	l, err := grpc_unixsocket.Listener(addr)
	if err != nil {
		return err
	}

	pa.cx = cmux.New(l)
	httpListener := pa.cx.Match(cmux.HTTP1Fast(), cmux.HTTP2())
	apintoListener := pa.cx.Match(api_apinto.Matcher())
	unknownListener := pa.cx.Match(cmux.Any())

	go func() {
		err := pa.server.Serve(httpListener)
		if err != nil {
			log.Info("http server error: ", err)
		}
		return
	}()
	go func() {
		err := pa.apintoServer.Server(apintoListener)
		if err != nil {
			return
		}
	}()

	go func() {
		for {
			conn, err := unknownListener.Accept()
			if err != nil {
				return
			}
			log.Warn("unknown conn: ", conn.RemoteAddr())
			conn.Write([]byte("-ERR unknown proto"))
			conn.Close()
		}

	}()
	go func() {
		err := pa.cx.Serve()
		if err != nil {

			return
		}
	}()
	return nil
}
