package traffic

import (
	cmuxMatch "github.com/eolinker/eosc/traffic/cmux-match"
	"github.com/eolinker/eosc/traffic/mixl"
	"net"
)

type MatcherData struct {
	data  map[string]*net.TCPListener
	ports map[int]*Matcher
}

func NewMatcherData(listeners map[string]*net.TCPListener) *MatcherData {

	m := &MatcherData{
		data:  make(map[string]*net.TCPListener),
		ports: make(map[int]*Matcher),
	}
	if listeners != nil {
		ports := make(map[int][]*net.TCPListener)
		for a, l := range listeners {
			_, port := readAddr(a)
			ports[port] = append(ports[port], l)
			m.data[a] = l
		}
		for port, ls := range ports {
			m.ports[port] = NewMatcher(port, ls...)
		}
	}

	return m
}

func (m *MatcherData) GetByPort(port int) *Matcher {
	v, has := m.ports[port]

	if !has {
		return nil
	}

	return v
}

func (m *MatcherData) clone() map[string]*net.TCPListener {
	ce := make(map[string]*net.TCPListener)
	for k, v := range m.data {
		ce[k] = v
	}
	return ce
}
func (m *MatcherData) All() map[string]*net.TCPListener {
	return m.data
}

type Matcher struct {
	mixListener *mixl.MixListener
	cmuxMatch.CMuxMatch
}

func NewMatcher(port int, ls ...*net.TCPListener) *Matcher {
	mixListener := mixl.NewMixListener(port, ls...)
	return &Matcher{mixListener: mixListener, CMuxMatch: cmuxMatch.NewMatch(mixListener)}
}

//	func (m *Matcher) Listeners() []*net.TCPListener {
//		return m.mixListener.Listeners()
//	}
func (m *Matcher) Close() error {
	m.mixListener.Close()
	m.CMuxMatch.Close()
	return nil
}
