package process_master

import (
	"context"
	"errors"
	"fmt"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/service"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"time"
)

var (
	ErrorAdminProcessNotInit = errors.New("admin process not init")
)

type UnixClient struct {
	addr    string
	client  *http.Client
	timeout time.Duration
}

func (uc *UnixClient) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	if uc.addr == "" {
		return nil, ErrorAdminProcessNotInit
	}
	return net.DialTimeout("unix", uc.addr, uc.timeout)
}
func (uc *UnixClient) Update(process *exec.Cmd) {
	if process == nil {
		uc.addr = ""
		return
	}
	uc.addr = service.ServerUnixAddr(process.Process.Pid, "admin")
}

func NewUnixClient() *UnixClient {
	ul := &UnixClient{}
	transport := &http.Transport{
		DialContext: ul.DialContext,
	}
	ul.client = &http.Client{Transport: transport}
	return ul
}
func (uc *UnixClient) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Debug("proxy:", request.RequestURI)

	if uc.addr == "" {
		writer.WriteHeader(http.StatusBadGateway)

		fmt.Fprintf(writer, "%v", ErrorAdminProcessNotInit)
		return
	}

	req, err := http.NewRequest(request.Method, request.RequestURI, request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v", err)
		log.Warnf("clone request to unix error:%v", err)
		return
	}
	req.URL.Scheme = "http"
	req.URL.Host = uc.addr
	req.Header = request.Header
	resp, err := uc.client.Do(req)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v", err)
		log.Errorf("proxy to unix err:%v", err)
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v", err)
		log.Errorf("fetch error: reading %v", err)
		return
	}

	for k, _ := range resp.Header {
		writer.Header().Set(k, resp.Header.Get(k))
	}
	writer.WriteHeader(resp.StatusCode)
	writer.Write(data)
}
