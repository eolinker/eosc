package cmuxMatch

import (
	"github.com/soheilhy/cmux"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type MatchType int
type CMuxMatch interface {
	Match(match MatchType) net.Listener
	SetReadTimeout(time.Duration)
	Close() error
}

var (
	_ CMuxMatch = (*cMuxMatch)(nil)
)

const (
	Any MatchType = iota
	Http1
	Https
	Http2
	Websocket
	GRPC
	matchTypeMax
)

var (
	matchers    [][]cmux.Matcher
	matcherName []string
)

func init() {
	matchers = make([][]cmux.Matcher, matchTypeMax)
	matcherName = make([]string, matchTypeMax)

	matchers[Any] = []cmux.Matcher{func(reader io.Reader) bool {
		return true
	}}
	matchers[Http1] = []cmux.Matcher{cmux.HTTP1Fast(), cmux.HTTP2()}
	matchers[Https] = []cmux.Matcher{cmux.TLS()}
	matchers[Http2] = []cmux.Matcher{cmux.HTTP2()}
	matchers[Websocket] = []cmux.Matcher{cmux.HTTP1HeaderField("Upgrade", "websocket")}
	matchers[GRPC] = []cmux.Matcher{cmux.HTTP2HeaderFieldPrefix("content-type", "application/grpc")}

	matcherName[Any] = "Any"
	matcherName[Http1] = "Http1"
	matcherName[Https] = "Https"
	matcherName[Http2] = "Http2"
	matcherName[Websocket] = "Websocket"
	matcherName[GRPC] = "GRPC"

}
func (t MatchType) String() string {
	if t > matchTypeMax || t < 0 {
		return "unknown"
	}
	return matcherName[t]
}
func (t MatchType) matcher() []cmux.Matcher {
	return matchers[t]
}

type cMuxMatch struct {
	cMux      cmux.CMux
	listeners []*shutListener

	root        *ListenerProxy
	lock        sync.Mutex
	readTimeOut time.Duration
}

func (m *cMuxMatch) SetReadTimeout(readTimeOut time.Duration) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.readTimeOut = readTimeOut
	if m.cMux != nil {
		m.cMux.SetReadTimeout(readTimeOut)
	}
}

func (m *cMuxMatch) Match(match MatchType) net.Listener {
	if match >= matchTypeMax || match < 0 {
		panic("invalid match type")
	}
	m.lock.Lock()
	defer m.lock.Unlock()

	if l := m.listeners[match]; l == nil {
		m.listeners[match] = newListener()
		m.rebuild()
	}
	return m.listeners[match]

}
func (m *cMuxMatch) rebuild() {
	m.root = m.root.Replace()

	if m.cMux != nil {
		m.cMux.Close()
		m.cMux = nil
	}

	nc := cmux.New(m.root)
	if m.readTimeOut != 0 {
		nc.SetReadTimeout(m.readTimeOut)
	}

	for i := GRPC; i >= Any; i-- {
		l := m.listeners[i]
		if l != nil {
			ms := i.matcher()
			l.reset(nc.Match(ms...))
		}
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	m.cMux = nc
	go func(nc cmux.CMux) {
		wg.Done()
		err := nc.Serve()
		if err != nil {
			log.Println("m")
			return
		}
	}(nc)
	wg.Wait()
}

func (m *cMuxMatch) Close() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.cMux != nil {
		m.cMux.Close()
		m.cMux = nil
	}

	m.root.ShutDown()
	for i, l := range m.listeners {
		if l != nil {
			l.Shutdown()
			m.listeners[i] = nil
		}
	}
	return nil
}

func NewMatch(l net.Listener) CMuxMatch {
	if l == nil {
		panic("mast init listener")
	}

	shutdown := make(chan struct{})
	m := &cMuxMatch{
		root:      NewListenerProxy(l, shutdown),
		listeners: make([]*shutListener, matchTypeMax),
	}
	go func() {
		<-shutdown
		m.Close()
	}()
	return m
}
