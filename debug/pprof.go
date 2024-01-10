package debug

import (
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"

	"github.com/eolinker/eosc/env"
	"github.com/eolinker/eosc/log"
)

var (
	mux = http.NewServeMux()
)

func Register(path string, handleFunc func(w http.ResponseWriter, r *http.Request)) {
	if mux != nil {
		mux.HandleFunc(path, handleFunc)
	}
}
func RunDebug(name string) {

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	addr, has := env.GetEnv(fmt.Sprintf("PPROF_%s", name))
	log.Debug("pprof addr:", addr, " ", name, fmt.Sprintf("PPROF_%s", name))
	if has {

		listen, err := net.Listen("tcp", addr)
		if err != nil {
			log.Warn("fail to listen pprof:", addr)
			panic(err)
			return
		}
		lAddr := listen.Addr().(*net.TCPAddr)

		log.Infof("start pprof:\thttp%s:%d", lAddr.IP.String(), lAddr.Port)

		server := http.Server{
			Handler: mux,
		}
		go func() {
			err := server.Serve(listen)
			if err != nil {
				return
			}
		}()

	} else {
		mux = nil
	}

}
