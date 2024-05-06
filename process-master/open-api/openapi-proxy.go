package open_api

import (
	"encoding/json"
	"fmt"
	open_api "github.com/eolinker/eosc/open-api"
	"io"
	"net/http"
	"strings"
	"sync"
)

var (
	client *http.Client = &http.Client{Transport: http.DefaultTransport}
)

type IRaftLeader interface {
	IsLeader() (bool, []string)
}

type OpenApiProxy struct {
	leaderHandler http.Handler
	raftSender    IRaftLeader
	pool          sync.Pool
}

func NewOpenApiProxy(sender IRaftLeader, leaderHandler http.Handler) *OpenApiProxy {
	p := &OpenApiProxy{
		leaderHandler: leaderHandler,
		raftSender:    sender,
		pool: sync.Pool{New: func() interface{} {
			return NewTemplateWriter()
		}},
	}
	return p
}
func (p *OpenApiProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	isLeader, leaderPeers := p.raftSender.IsLeader()

	if isLeader {
		p.doProxy(w, r)
	} else {
		p.doProxyToLeader(w, r, leaderPeers)
	}
}

func (p *OpenApiProxy) doProxy(w http.ResponseWriter, r *http.Request) {

	buf := p.pool.Get().(*_ProxyWriterBuffer)
	buf.Reset()
	defer p.pool.Put(buf)
	p.leaderHandler.ServeHTTP(buf, r)
	if buf.statusCode != http.StatusOK {
		buf.WriteTo(w)
		return
	}
	//buf.WriteHeaderTo(w)

	res := new(open_api.Response)
	err := json.Unmarshal(buf.buf.Bytes(), res)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"code":%d,"error":"%s","re","message":"%s"}`, http.StatusInternalServerError, err.Error(), buf.buf.String())
		return
	}
	//if res.Event != nil {
	//	for _, event := range res.Event {
	//		err := p.raftSender.Send(event.Event, event.Namespace, event.Key, event.Data)
	//		log.Debug("open api send:", res.Event)
	//		if err != nil {
	//			log.Errorf("open api raft:%v", err)
	//		}
	//	}
	//
	//}
	if res.Header != nil {
		for k := range res.Header {
			w.Header().Set(k, res.Header.Get(k))
		}
	}

	w.WriteHeader(res.StatusCode)
	w.Write(res.Data)

}

func (p *OpenApiProxy) doProxyToLeader(w http.ResponseWriter, org *http.Request, leaders []string) {
	var err error
	var response *http.Response
	for _, leader := range leaders {
		r, _ := http.NewRequest(org.Method, fmt.Sprintf("%s%s", leader, org.RequestURI), org.Body)
		r.Header = org.Header

		response, err = client.Do(r)
		if err != nil {
			continue
		}
		defer response.Body.Close()
		for key, value := range response.Header {
			w.Header().Set(key, strings.Join(value, ","))
		}
		w.WriteHeader(response.StatusCode)
		io.Copy(w, response.Body)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Header().Set("content-type", "application/js")
		fmt.Fprintf(w, `{"code":%d,"error":"%s"}`, http.StatusBadGateway, err.Error())
		return
	}

}

func formatPath(path string) []string {
	if len(path) == 0 {
		return []string{"/*more"}
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	i := strings.LastIndexAny(path, "/")
	last := path[i+1:]

	if len(last) == 0 { // 结尾为/，例如 /api/
		return []string{path, fmt.Sprintf("%s*more", path)}
	}
	if last == "*" { // 结尾为/*
		return []string{fmt.Sprintf("%smore", path)}
	}
	if strings.HasPrefix(last, "*") { // 结尾为 /*{any}
		return []string{path}
	}
	return []string{path, fmt.Sprintf("%s/*more", path)}
}
