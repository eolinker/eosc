package traffic

import (
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	cmuxMatch "github.com/eolinker/eosc/traffic/cmux-match"
	"github.com/eolinker/eosc/traffic/mixl"
	"net"
	"strconv"
)

type MatcherData struct {
	data eosc.IUntyped
}

func NewMatcherData(tfConf ...*PbTraffic) *MatcherData {

	m := &MatcherData{
		data: eosc.NewUntyped(),
	}
	listeners, err := toListeners(tfConf)
	log.Debug("read listeners: ", len(listeners))
	if err != nil {
		log.Warn("read listeners:", err)
	}

	lm := make(map[int][]*net.TCPListener)
	for _, l := range listeners {
		p := readPort(l.Addr())
		lm[p] = append(lm[p], l)
	}
	for p, ls := range lm {
		mixlistener := mixl.NewMixListener(p, ls...)
		m.Set(p, &Matcher{
			mixListener: mixlistener,
			CMuxMatch:   cmuxMatch.NewMatch(mixlistener),
		})
	}
	return m
}

func (m *MatcherData) Set(port int, mux *Matcher) {
	m.data.Set(strconv.Itoa(port), mux)
}
func (m *MatcherData) Get(port int) *Matcher {
	o, has := m.data.Get(strconv.Itoa(port))
	if !has {
		return nil
	}
	return o.(*Matcher)
}
func (m *MatcherData) Del(port int) (*Matcher, bool) {
	o, ok := m.data.Del(strconv.Itoa(port))
	if ok {
		return o.(*Matcher), true
	}
	return nil, false
}
func (m *MatcherData) Clone() *MatcherData {
	return &MatcherData{
		data: m.data.Clone(),
	}
}
func (m *MatcherData) All() map[int]*Matcher {
	all := m.data.All()
	rs := make(map[int]*Matcher)
	for p, o := range all {
		port, err := strconv.Atoi(p)
		if err != nil {
			continue
		}
		rs[port] = o.(*Matcher)
	}
	return rs
}

type Matcher struct {
	mixListener *mixl.MixListener
	cmuxMatch.CMuxMatch
}

func NewMatcher(port int, ls ...*net.TCPListener) *Matcher {
	mixListener := mixl.NewMixListener(port, ls...)
	return &Matcher{mixListener: mixListener, CMuxMatch: cmuxMatch.NewMatch(mixListener)}
}

func (m *Matcher) Listeners() []*net.TCPListener {
	return m.mixListener.Listeners()
}
func (m *Matcher) Close() error {
	m.mixListener.Close()
	m.CMuxMatch.Close()
	return nil
}
