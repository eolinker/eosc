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
	"github.com/eolinker/eosc/config"
	grpc_unixsocket "github.com/eolinker/eosc/grpc-unixsocket"
	open_api "github.com/eolinker/eosc/open-api"
	"github.com/eolinker/eosc/professions"
	"github.com/eolinker/eosc/require"
	"github.com/eolinker/eosc/service"
	"github.com/eolinker/eosc/setting"
	"github.com/eolinker/eosc/traffic"
	"github.com/eolinker/eosc/variable"
	
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	
	"github.com/eolinker/eosc/process"
	
	"github.com/eolinker/eosc/extends"
	
	"github.com/eolinker/eosc/common/bean"
	
	"google.golang.org/protobuf/proto"
	
	"github.com/eolinker/eosc/utils"
	
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
)

func Process() {
	utils.InitStdTransport(eosc.ProcessAdmin)
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
	ctx        context.Context
	cancelFunc context.CancelFunc
	once       sync.Once
	reg        eosc.IExtenderDriverRegister
	router     *httprouter.Router
	
	apiLocker sync.Mutex
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
			
			}
		default:
			continue
		}
	}
	
}

//NewProcessAdmin 创建新的admin进程
//启动时通过stdin传输配置信息
func NewProcessAdmin(parent context.Context, arg map[string]map[string][]byte) (*ProcessAdmin, error) {
	cfg := &config.ListensMsg{}
	var tf traffic.ITraffic = traffic.NewEmptyTraffic()
	bean.Injection(&tf)
	bean.Injection(&cfg)
	register := initExtender(arg[eosc.NamespaceExtender])
	var extenderDrivers eosc.IExtenderDrivers = register
	bean.Injection(&extenderDrivers)
	//for namespace, a := range arg {
	//	log.Debug("namespace is ", namespace)
	//	for k, v := range a {
	//		log.Debug("key is ", k, " v is ", string(v))
	//	}
	//
	//}
	
	ctx, cancelFunc := context.WithCancel(parent)
	p := &ProcessAdmin{
		ctx:        ctx,
		cancelFunc: cancelFunc,
		router:     httprouter.New(),
	}
	extenderRequire := require.NewRequireManager()
	extenderData := NewExtenderData(arg[eosc.NamespaceExtender], extenderRequire)
	NewExtenderOpenApi(extenderData).Register(p.router)
	
	ps := professions.NewProfessions(register)
	
	ps = NewProfessionsRequire(ps, extenderRequire)
	ps.Reset(professionConfig(arg[eosc.NamespaceProfession]))
	
	vd := variable.NewVariables(arg[eosc.NamespaceVariable])
	
	settingApi := NewSettingApi(filerSetting(arg[eosc.NamespaceWorker], Setting, true), vd)
	wd := NewWorkerDatas(filerSetting(arg[eosc.NamespaceWorker], Setting, false))
	var iWorkers eosc.IWorkers = wd
	bean.Injection(&iWorkers)
	
	bean.Check()
	
	ws := NewWorkers(ps, wd, vd)
	
	// openAPI handler register
	NewProfessionApi(ps, wd).Register(p.router)
	NewWorkerApi(ws, settingApi.request).Register(p.router)
	settingApi.RegisterSetting(p.router)
	NewExportApi(extenderData, ps, ws).Register(p.router)
	NewVariableApi(extenderData, ws, vd, setting.GetSettings()).Register(p.router)
	
	p.router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := &open_api.Response{
			StatusCode: 404,
			Header:     nil,
			Data:       nil,
			Event:      nil,
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
		pa.cancelFunc()
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

func filerSetting(confs map[string][]byte, name string, yes bool) map[string][]byte {
	name = strings.ToLower(name)
	sets := make(map[string][]byte)
	for id, data := range confs {
		profession, _, _ := eosc.SplitWorkerId(id)
		if (strings.ToLower(profession) == name) == yes {
			sets[id] = data
		}
	}
	return sets
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
	
	addr := service.ServerUnixAddr(os.Getpid(), eosc.ProcessAdmin)
	syscall.Unlink(addr)
	log.Info("start admin unix server: ", addr)
	l, err := grpc_unixsocket.Listener(addr)
	if err != nil {
		return err
	}
	server := http.Server{Handler: pa}
	go func() {
		err := server.Serve(l)
		if err != nil {
			log.Info("http server error: ", err)
		}
		return
	}()
	go func() {
		<-pa.ctx.Done()
		server.Shutdown(context.Background())
		syscall.Unlink(addr)
	}()
	
	return nil
}
