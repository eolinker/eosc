package main

import (
	"context"
	"fmt"
	"github.com/eolinker/eosc/config"
	"github.com/eolinker/eosc/etcd"
	"github.com/eolinker/eosc/log"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {

	transport := log.NewTransport(os.Stderr, log.DebugLevel)
	transport.SetFormatter(&log.LineFormatter{
		TimestampFormat:  "2006-01-02 15:04:05",
		CallerPrettyfier: nil,
	})
	log.Reset(transport)
	log.SetPrefix(fmt.Sprintf("[demo]"))
	conf, err := config.GetConfig()
	if err != nil {
		log.Debug(conf)
		return
	}
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", conf.Admin.IP, conf.Admin.Listen))
	if err != nil {
		log.Debug(err)
		return
	}

	mux := http.NewServeMux()
	ser := http.Server{
		Handler: mux,
	}

	go ser.Serve(listen)
	server, err := etcd.NewServer(context.Background(), mux)
	if err != nil {
		log.Debug(err)
		return
	}
	mux.HandleFunc("/do/leave", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, server.Leave())
	})
	mux.HandleFunc("/do/join", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		target := r.FormValue("target")
		if target == "" {
			return
		}
		if !strings.HasPrefix(target, "http") {
			target = "http://" + target
		}
		fmt.Fprint(w, server.Join(target))

	})

	server.Watch("/", new(DemoHandler))

	select {}

}

type DemoHandler struct {
}

func (d *DemoHandler) Put(key string, value []byte) error {
	log.Debugf("put:%s=%s\n", key, string(value))
	return nil
}

func (d *DemoHandler) Delete(key string) error {
	log.Debugf("delete:%s\n", key)
	return nil
}

func (d *DemoHandler) Reset(values []*etcd.KValue) {
	log.Debug("reset start===============")
	for _, v := range values {
		log.Debugf("\t%s=%s\n", string(v.Key), string(v.Value))
	}
	log.Debug("reset end=============")
}
