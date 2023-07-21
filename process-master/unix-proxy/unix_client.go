/*
 * Copyright (c) 2023. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package unix_proxy

import (
	"bufio"
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

const BuffSize = 4096

var (
	ErrorProcessNotInit = errors.New("process not init")
)

type UnixClient struct {
	addr    string
	name    string
	client  http.RoundTripper
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
	ul.client = transport
	return ul
}
func (uc *UnixClient) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	log.Debug("proxy:", request.RequestURI)

	if uc.addr == "" {
		w.WriteHeader(http.StatusBadGateway)

		fmt.Fprintf(w, "%s %s", uc.name, ErrorProcessNotInit.Error())
		return
	}
	request.URL.Scheme = "http"
	request.URL.Host = uc.name
	if !tokenListContainsValue(request.Header, "Connection", "Upgrade") {
		response, err := uc.client.RoundTrip(request)
		if err != nil {
			return
		}

		defer func() {
			response.Body.Close()
		}()
		w.WriteHeader(response.StatusCode)
		for k, vs := range response.Header {
			for _, v := range vs {
				w.Header().Add(k, v)
			}
		}
		io.Copy(w, response.Body)
	} else {

		h, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		var brw *bufio.ReadWriter
		netConn, brw, err := h.Hijack()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer func() {

			netConn.Close()
		}()
		if brw.Reader.Buffered() > 0 {

			return
		}

		upstream, resp, err := uc.DialContextUpgrade(request)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer upstream.Close()
		err = resp.Write(netConn)
		if err != nil {
			return
		}
		go func() {
			io.Copy(netConn, upstream)
		}()
		io.Copy(upstream, netConn)
	}

}
