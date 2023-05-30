package process_master

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os/exec"
	"time"

	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/service"
)

var (
	ErrorProcessNotInit = errors.New("process not init")
)

type UnixClient struct {
	addr    string
	name    string
	client  *http.Client
	timeout time.Duration
}

func (uc *UnixClient) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	if uc.addr == "" {
		return nil, fmt.Errorf("%s %w", uc.name, ErrorProcessNotInit)
	}
	return net.DialTimeout("unix", uc.addr, uc.timeout)
}
func (uc *UnixClient) Update(process *exec.Cmd) {
	log.DebugF("unix client update: %s %s", uc.name, process)
	if process == nil {
		uc.addr = ""
		return
	}
	uc.addr = service.ServerUnixAddr(process.Process.Pid, uc.name)
}

func NewUnixClient(name string) *UnixClient {
	ul := &UnixClient{
		name: name,
	}
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

		fmt.Fprintf(writer, "%s %s", uc.name, ErrorProcessNotInit.Error())
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
	if resp.Header.Get("Content-Length") != "" {
		defer resp.Body.Close()
		io.Copy(writer, resp.Body)
		return
	}
	flusher, ok := writer.(http.Flusher)
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	for k, vs := range resp.Header {
		for _, v := range vs {
			writer.Header().Add(k, v)
		}
	}
	defer resp.Body.Close()

	buf := make([]byte, 4096)
	for {
		select {
		case <-request.Context().Done():
			return
		default:

			n, err2 := resp.Body.Read(buf)
			if n > 0 {
				_, err := writer.Write(buf[:n])
				if err != nil {
					return
				}
				flusher.Flush()
			}
			if err2 != nil {
				return
			}
		}

	}
}
