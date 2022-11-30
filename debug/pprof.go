package debug

import (
	"github.com/eolinker/eosc/env"
	"github.com/eolinker/eosc/log"
	"net"
	"net/http"
	"net/http/pprof"
	"time"
)

var (
	mux = http.NewServeMux()
)

func Register(path string, handleFunc func(w http.ResponseWriter, r *http.Request)) {

	mux.HandleFunc(path, handleFunc)
}
func RunDebug() {

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	go func() {
		// 延迟3s启动

		time.Sleep(time.Second * 3)

		addr, has := env.GetEnv("pprof")
		if has {
			listen, err := net.Listen("tcp", addr)
			if err != nil {
				log.Warn("fail to listen pprof:", addr)
				return
			}
			lAddr := listen.Addr().(*net.TCPAddr)

			log.Infof("start pprof:\thttp%s:%d", lAddr.IP.String(), lAddr.Port)

			server := http.Server{
				Handler: mux,
			}
			err = server.Serve(listen)
			if err != nil {
				return
			}
		}

	}()

}
