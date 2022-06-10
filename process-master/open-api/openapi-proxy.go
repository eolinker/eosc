package open_api

import (
	"encoding/json"
	"fmt"
	"github.com/eolinker/eosc/log"
	open_api "github.com/eolinker/eosc/open-api"
	"github.com/eolinker/eosc/raft"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
	"sync"
)

var (
	client *http.Client = &http.Client{Transport: http.DefaultTransport}
)

type IRaftSender interface {
	Send(event string, namespace string, key string, data []byte) error
	IsLeader() (bool, *raft.NodeInfo, error)
}

type OpenApiProxy struct {
	excludeRouter *httprouter.Router
	leaderHandler http.Handler
	raftSender    IRaftSender
	pool          sync.Pool
}

func NewOpenApiProxy(sender IRaftSender, leaderHandler http.Handler) *OpenApiProxy {
	p := &OpenApiProxy{
		excludeRouter: httprouter.New(),
		leaderHandler: leaderHandler,
		raftSender:    sender,
		pool: sync.Pool{New: func() interface{} {
			return NewTemplateWriter()
		}},
	}
	return p
}
func (p *OpenApiProxy) ExcludeHandle(method, path string, handler httprouter.Handle) {
	for _, ph := range formatPath(path) {
		p.excludeRouter.Handle(method, ph, handler)
	}
}
func (p *OpenApiProxy) ExcludeHandleFunc(method, path string, handler http.HandlerFunc) {
	for _, ph := range formatPath(path) {
		p.excludeRouter.HandlerFunc(method, ph, handler)
	}
}
func (p *OpenApiProxy) ExcludeHandlesFunc(path string, handler http.HandlerFunc) {
	p.ExcludeHandleFunc(http.MethodGet, path, handler)
	p.ExcludeHandleFunc(http.MethodHead, path, handler)
	p.ExcludeHandleFunc(http.MethodPost, path, handler)
	p.ExcludeHandleFunc(http.MethodPut, path, handler)
	p.ExcludeHandleFunc(http.MethodPatch, path, handler)
	p.ExcludeHandleFunc(http.MethodDelete, path, handler)
	p.ExcludeHandleFunc(http.MethodConnect, path, handler)
	p.ExcludeHandleFunc(http.MethodOptions, path, handler)
	p.ExcludeHandleFunc(http.MethodTrace, path, handler)
}
func (p *OpenApiProxy) ExcludeHandles(path string, handler httprouter.Handle) {
	p.ExcludeHandle(http.MethodGet, path, handler)
	p.ExcludeHandle(http.MethodHead, path, handler)
	p.ExcludeHandle(http.MethodPost, path, handler)
	p.ExcludeHandle(http.MethodPut, path, handler)
	p.ExcludeHandle(http.MethodPatch, path, handler)
	p.ExcludeHandle(http.MethodDelete, path, handler)
	p.ExcludeHandle(http.MethodConnect, path, handler)
	p.ExcludeHandle(http.MethodOptions, path, handler)
	p.ExcludeHandle(http.MethodTrace, path, handler)
}
func (p *OpenApiProxy) ExcludeHandler(method, path string, handler http.Handler) {
	for _, ph := range formatPath(path) {
		p.excludeRouter.Handler(method, ph, handler)
	}
}
func (p *OpenApiProxy) ExcludeHandlers(path string, handler http.Handler) {
	p.ExcludeHandler(http.MethodGet, path, handler)
	p.ExcludeHandler(http.MethodHead, path, handler)
	p.ExcludeHandler(http.MethodPost, path, handler)
	p.ExcludeHandler(http.MethodPut, path, handler)
	p.ExcludeHandler(http.MethodPatch, path, handler)
	p.ExcludeHandler(http.MethodDelete, path, handler)
	p.ExcludeHandler(http.MethodConnect, path, handler)
	p.ExcludeHandler(http.MethodOptions, path, handler)
	p.ExcludeHandler(http.MethodTrace, path, handler)
}

func (p *OpenApiProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, params, _ := p.excludeRouter.Lookup(r.Method, r.URL.Path)
	if handler != nil {
		handler(w, r, params)
		return
	}

	isLeader, leadNode, err := p.raftSender.IsLeader()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		log.Warnf("apinto no leader")
		return
	}

	if isLeader {
		p.doProxy(w, r)
	} else {
		p.doProxyToLeader(w, r, leadNode.Addr)
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
	if res.Event != nil {
		err := p.raftSender.Send(res.Event.Event, res.Event.Namespace, res.Event.Key, res.Event.Data)
		log.Debug("open api send:", res.Event)
		if err != nil {
			log.Errorf("open api raft:%v", err)
		}
	}
	if res.Header != nil {
		for k := range res.Header {
			w.Header().Set(k, res.Header.Get(k))
		}
	}

	w.WriteHeader(res.StatusCode)
	w.Write(res.Data)

}

func (p *OpenApiProxy) doProxyToLeader(w http.ResponseWriter, org *http.Request, leader string) {

	r, _ := http.NewRequest(org.Method, fmt.Sprintf("%s%s", leader, org.RequestURI), org.Body)
	r.Header = org.Header
	response, err := client.Do(r)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Header().Set("content-type", "application/js")
		fmt.Fprintf(w, `{"code":%d,"error":"%s"}`, http.StatusBadGateway, err.Error())
		return
	}

	response.Write(w)

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
