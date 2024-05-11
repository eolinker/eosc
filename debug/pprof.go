package debug

import (
	"fmt"
	"github.com/eolinker/eosc/env"
	"github.com/eolinker/eosc/log"
	"net"
	"net/http"
	"net/http/pprof"
)

var ()

func RunDebug(name string) {

	addr, has := env.GetEnv(fmt.Sprintf("PPROF_%s", name))
	log.Debug("pprof addr:", addr, " ", name, ": ", fmt.Sprintf("PPROF_%s", name))
	if !has || addr == "" {
		return
	}
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Warn("fail to listen pprof:", addr)

		return
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	lAddr := listen.Addr().(*net.TCPAddr)

	log.Infof("start pprof:\thttp%s:%d", lAddr.IP.String(), lAddr.Port)

	go func() {
		server := http.Server{
			Handler: mux,
		}
		err := server.Serve(listen)
		if err != nil {
			return
		}
	}()

}
